package keeper

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/line/ostracon/libs/log"
	ostos "github.com/line/ostracon/libs/os"

	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/store/prefix"
	store "github.com/line/lbm-sdk/store/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/types/module"
	clienttypes "github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	ibcexported "github.com/line/lbm-sdk/x/ibc/core/exported"
	"github.com/line/lbm-sdk/x/upgrade/types"
)

// UpgradeInfoFileName file to store upgrade information
const UpgradeInfoFileName string = "upgrade-info.json"

type Keeper struct {
	homePath           string
	skipUpgradeHeights map[int64]bool
	storeKey           sdk.StoreKey
	cdc                codec.BinaryMarshaler
	upgradeHandlers    map[string]types.UpgradeHandler
}

// NewKeeper constructs an upgrade Keeper
func NewKeeper(skipUpgradeHeights map[int64]bool, storeKey sdk.StoreKey, cdc codec.BinaryMarshaler, homePath string) Keeper {
	return Keeper{
		homePath:           homePath,
		skipUpgradeHeights: skipUpgradeHeights,
		storeKey:           storeKey,
		cdc:                cdc,
		upgradeHandlers:    map[string]types.UpgradeHandler{},
	}
}

// SetUpgradeHandler sets an UpgradeHandler for the upgrade specified by name. This handler will be called when the upgrade
// with this name is applied. In order for an upgrade with the given name to proceed, a handler for this upgrade
// must be set even if it is a no-op function.
func (k Keeper) SetUpgradeHandler(name string, upgradeHandler types.UpgradeHandler) {
	k.upgradeHandlers[name] = upgradeHandler
}

// SetModuleVersionMap saves a given version map to state
func (k Keeper) SetModuleVersionMap(ctx sdk.Context, vm module.VersionMap) {
	if len(vm) > 0 {
		store := ctx.KVStore(k.storeKey)
		versionStore := prefix.NewStore(store, []byte{types.VersionMapByte})
		// Even though the underlying store (cachekv) store is sorted, we still
		// prefer a deterministic iteration order of the map, to avoid undesired
		// surprises if we ever change stores.
		sortedModNames := make([]string, 0, len(vm))

		for key := range vm {
			sortedModNames = append(sortedModNames, key)
		}
		sort.Strings(sortedModNames)

		for _, modName := range sortedModNames {
			ver := vm[modName]
			nameBytes := []byte(modName)
			verBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(verBytes, ver)
			versionStore.Set(nameBytes, verBytes)
		}
	}
}

// GetModuleVersionMap returns a map of key module name and value module consensus version
// as defined in ADR-041.
func (k Keeper) GetModuleVersionMap(ctx sdk.Context) module.VersionMap {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, []byte{types.VersionMapByte})

	vm := make(module.VersionMap)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		moduleBytes := it.Key()
		// first byte is prefix key, so we remove it here
		name := string(moduleBytes[1:])
		moduleVersion := binary.BigEndian.Uint64(it.Value())
		vm[name] = moduleVersion
	}

	return vm
}

// GetModuleVersions gets a slice of module consensus versions
func (k Keeper) GetModuleVersions(ctx sdk.Context) []*types.ModuleVersion {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, []byte{types.VersionMapByte})
	defer it.Close()

	mv := make([]*types.ModuleVersion, 0)
	for ; it.Valid(); it.Next() {
		moduleBytes := it.Key()
		name := string(moduleBytes[1:])
		moduleVersion := binary.BigEndian.Uint64(it.Value())
		mv = append(mv, &types.ModuleVersion{
			Name:    name,
			Version: moduleVersion,
		})
	}
	return mv
}

// gets the version for a given module, and returns true if it exists, false otherwise
func (k Keeper) getModuleVersion(ctx sdk.Context, name string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, []byte{types.VersionMapByte})
	defer it.Close()

	for ; it.Valid(); it.Next() {
		moduleName := string(it.Key()[1:])
		if moduleName == name {
			version := binary.BigEndian.Uint64(it.Value())
			return version, true
		}
	}
	return 0, false
}

// ScheduleUpgrade schedules an upgrade based on the specified plan.
// If there is another Plan already scheduled, it will overwrite it
// (implicitly cancelling the current plan)
// ScheduleUpgrade will also write the upgraded client to the upgraded client path
// if an upgraded client is specified in the plan
func (k Keeper) ScheduleUpgrade(ctx sdk.Context, plan types.Plan) error {
	if err := plan.ValidateBasic(); err != nil {
		return err
	}

	if plan.Time.Unix() > 0 {
		if !plan.Time.After(ctx.BlockHeader().Time) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "upgrade cannot be scheduled in the past")
		}
	} else if plan.Height <= ctx.BlockHeight() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "upgrade cannot be scheduled in the past")
	}

	if k.GetDoneHeight(ctx, plan.Name) != 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "upgrade with name %s has already been completed", plan.Name)
	}

	store := ctx.KVStore(k.storeKey)

	// clear any old IBC state stored by previous plan
	oldPlan, exists := k.GetUpgradePlan(ctx)
	if exists && oldPlan.IsIBCPlan() {
		k.ClearIBCState(ctx, oldPlan.Height-1)
	}

	bz := k.cdc.MustMarshalBinaryBare(&plan)
	store.Set(types.PlanKey(), bz)

	if plan.IsIBCPlan() {
		// Set UpgradedClientState in store
		clientState, err := clienttypes.UnpackClientState(plan.UpgradedClientState)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "could not unpack clientstate: %v", err)
		}
		// sets the new upgraded client in last height committed on this chain is at plan.Height,
		// since the chain will panic at plan.Height and new chain will resume at plan.Height
		return k.SetUpgradedClient(ctx, plan.Height, clientState)
	}
	return nil
}

// SetUpgradedClient sets the expected upgraded client for the next version of this chain at the last height the current chain will commit.
func (k Keeper) SetUpgradedClient(ctx sdk.Context, planHeight int64, cs ibcexported.ClientState) error {
	store := ctx.KVStore(k.storeKey)

	// zero out any custom fields before setting
	cs = cs.ZeroCustomFields()
	bz, err := clienttypes.MarshalClientState(k.cdc, cs)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "could not marshal clientstate: %v", err)
	}

	store.Set(types.UpgradedClientKey(planHeight), bz)
	return nil
}

// GetUpgradedClient gets the expected upgraded client for the next version of this chain
func (k Keeper) GetUpgradedClient(ctx sdk.Context, height int64) (ibcexported.ClientState, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.UpgradedClientKey(height))
	if len(bz) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "upgraded client not found in store for height %d", height)
	}

	clientState, err := clienttypes.UnmarshalClientState(k.cdc, bz)
	if err != nil {
		return nil, err
	}
	return clientState, nil
}

// SetUpgradedConsensusState set the expected upgraded consensus state for the next version of this chain
// using the last height committed on this chain.
func (k Keeper) SetUpgradedConsensusState(ctx sdk.Context, planHeight int64, cs ibcexported.ConsensusState) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := clienttypes.MarshalConsensusState(k.cdc, cs)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "could not marshal consensus state: %v", err)
	}

	store.Set(types.UpgradedConsStateKey(planHeight), bz)
	return nil
}

// GetUpgradedConsensusState set the expected upgraded consensus state for the next version of this chain
func (k Keeper) GetUpgradedConsensusState(ctx sdk.Context, lastHeight int64) (ibcexported.ConsensusState, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.UpgradedConsStateKey(lastHeight))
	if len(bz) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "upgraded consensus state not found in store for height: %d", lastHeight)
	}
	consState, err := clienttypes.UnmarshalConsensusState(k.cdc, bz)
	if err != nil {
		return nil, err
	}
	return consState, nil
}

// GetDoneHeight returns the height at which the given upgrade was executed
func (k Keeper) GetDoneHeight(ctx sdk.Context, name string) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.DoneByte})
	bz := store.Get([]byte(name))
	if len(bz) == 0 {
		return 0
	}

	return int64(binary.BigEndian.Uint64(bz))
}

// ClearIBCState clears any planned IBC state
func (k Keeper) ClearIBCState(ctx sdk.Context, lastHeight int64) {
	// delete IBC client and consensus state from store if this is IBC plan
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UpgradedClientKey(lastHeight))
	store.Delete(types.UpgradedConsStateKey(lastHeight))
}

// ClearUpgradePlan clears any schedule upgrade
func (k Keeper) ClearUpgradePlan(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PlanKey())
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetUpgradePlan returns the currently scheduled Plan if any, setting havePlan to true if there is a scheduled
// upgrade or false if there is none
func (k Keeper) GetUpgradePlan(ctx sdk.Context) (plan types.Plan, havePlan bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PlanKey())
	if bz == nil {
		return plan, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &plan)
	return plan, true
}

// setDone marks this upgrade name as being done so the name can't be reused accidentally
func (k Keeper) setDone(ctx sdk.Context, name string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{types.DoneByte})
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(ctx.BlockHeight()))
	store.Set([]byte(name), bz)
}

// HasHandler returns true iff there is a handler registered for this name
func (k Keeper) HasHandler(name string) bool {
	_, ok := k.upgradeHandlers[name]
	return ok
}

// ApplyUpgrade will execute the handler associated with the Plan and mark the plan as done.
func (k Keeper) ApplyUpgrade(ctx sdk.Context, plan types.Plan) {
	handler := k.upgradeHandlers[plan.Name]
	if handler == nil {
		panic("ApplyUpgrade should never be called without first checking HasHandler")
	}

	updatedVM, err := handler(ctx, plan, k.GetModuleVersionMap(ctx))
	if err != nil {
		panic(err)
	}

	k.SetModuleVersionMap(ctx, updatedVM)

	// Must clear IBC state after upgrade is applied as it is stored separately from the upgrade plan.
	// This will prevent resubmission of upgrade msg after upgrade is already completed.
	if plan.IsIBCPlan() {
		k.ClearIBCState(ctx, plan.Height-1)
	}
	k.ClearUpgradePlan(ctx)
	k.setDone(ctx, plan.Name)
}

// IsSkipHeight checks if the given height is part of skipUpgradeHeights
func (k Keeper) IsSkipHeight(height int64) bool {
	return k.skipUpgradeHeights[height]
}

// DumpUpgradeInfoToDisk writes upgrade information to UpgradeInfoFileName.
func (k Keeper) DumpUpgradeInfoToDisk(height int64, name string) error {
	upgradeInfoFilePath, err := k.GetUpgradeInfoPath()
	if err != nil {
		return err
	}

	upgradeInfo := store.UpgradeInfo{
		Name:   name,
		Height: height,
	}
	info, err := json.Marshal(upgradeInfo)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(upgradeInfoFilePath, info, 0600)
}

// GetUpgradeInfoPath returns the upgrade info file path
func (k Keeper) GetUpgradeInfoPath() (string, error) {
	upgradeInfoFileDir := path.Join(k.getHomeDir(), "data")
	err := ostos.EnsureDir(upgradeInfoFileDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return filepath.Join(upgradeInfoFileDir, UpgradeInfoFileName), nil
}

// getHomeDir returns the height at which the given upgrade was executed
func (k Keeper) getHomeDir() string {
	return k.homePath
}

// ReadUpgradeInfoFromDisk returns the name and height of the upgrade which is
// written to disk by the old binary when panicking. An error is returned if
// the upgrade path directory cannot be created or if the file exists and
// cannot be read or if the upgrade info fails to unmarshal.
func (k Keeper) ReadUpgradeInfoFromDisk() (store.UpgradeInfo, error) {
	var upgradeInfo store.UpgradeInfo

	upgradeInfoPath, err := k.GetUpgradeInfoPath()
	if err != nil {
		return upgradeInfo, err
	}

	data, err := ioutil.ReadFile(upgradeInfoPath)
	if err != nil {
		// if file does not exist, assume there are no upgrades
		if os.IsNotExist(err) {
			return upgradeInfo, nil
		}

		return upgradeInfo, err
	}

	if err := json.Unmarshal(data, &upgradeInfo); err != nil {
		return upgradeInfo, err
	}

	return upgradeInfo, nil
}
