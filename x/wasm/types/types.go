package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	wasmvmtypes "github.com/line/wasmvm/types"
)

const (
	defaultMemoryCacheSize   uint32 = 100 // in MiB
	defaultQueryGasLimit     uint64 = 3000000
	defaultContractDebugMode        = false
)

var AllContractStatus = []ContractStatus{
	ContractStatusInactive,
	ContractStatusActive,
}

func (c ContractStatus) String() string {
	switch c {
	case ContractStatusActive:
		return "Active"
	case ContractStatusInactive:
		return "Inactive"
	}
	return "Unspecified"
}

func (c *ContractStatus) UnmarshalText(text []byte) error {
	for _, v := range AllContractStatus {
		if v.String() == string(text) {
			*c = v
			return nil
		}
	}
	*c = ContractStatusUnspecified
	return nil
}

func (c ContractStatus) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *ContractStatus) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(c)
}

func (c *ContractStatus) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	return json.Unmarshal(data, c)
}

func (m Model) ValidateBasic() error {
	if len(m.Key) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "key")
	}
	return nil
}

func (c CodeInfo) ValidateBasic() error {
	if len(c.CodeHash) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "code hash")
	}
	if err := sdk.ValidateAccAddress(c.Creator); err != nil {
		return sdkerrors.Wrap(err, "creator")
	}
	if err := validateSourceURL(c.Source); err != nil {
		return sdkerrors.Wrap(err, "source")
	}
	if err := validateBuilder(c.Builder); err != nil {
		return sdkerrors.Wrap(err, "builder")
	}
	if err := c.InstantiateConfig.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "instantiate config")
	}
	return nil
}

// NewCodeInfo fills a new Contract struct
func NewCodeInfo(codeHash []byte, creator sdk.AccAddress, source string, builder string, instantiatePermission AccessConfig) CodeInfo {
	return CodeInfo{
		CodeHash:          codeHash,
		Creator:           creator.String(),
		Source:            source,
		Builder:           builder,
		InstantiateConfig: instantiatePermission,
	}
}

var AllCodeHistoryTypes = []ContractCodeHistoryOperationType{ContractCodeHistoryOperationTypeGenesis, ContractCodeHistoryOperationTypeInit, ContractCodeHistoryOperationTypeMigrate}

// NewContractInfo creates a new instance of a given WASM contract info
func NewContractInfo(codeID uint64, creator, admin sdk.AccAddress, label string, createdAt *AbsoluteTxPosition, status ContractStatus) ContractInfo {
	var adminAddr string
	if !admin.Empty() {
		adminAddr = admin.String()
	}
	return ContractInfo{
		CodeID:  codeID,
		Creator: creator.String(),
		Admin:   adminAddr,
		Label:   label,
		Created: createdAt,
		Status:  status,
	}
}

// validatable is an optional interface that can be implemented by an ContractInfoExtension to enable validation
type validatable interface {
	ValidateBasic() error
}

// ValidateBasic does syntax checks on the data. If an extension is set and has the `ValidateBasic() error` method, then
// the method is called as well. It is recommend to implement `ValidateBasic` so that the data is verified in the setter
// but also in the genesis import process.
func (c *ContractInfo) ValidateBasic() error {
	if c.CodeID == 0 {
		return sdkerrors.Wrap(ErrEmpty, "code id")
	}
	if err := sdk.ValidateAccAddress(c.Creator); err != nil {
		return sdkerrors.Wrap(err, "creator")
	}
	if len(c.Admin) != 0 {
		if err := sdk.ValidateAccAddress(c.Admin); err != nil {
			return sdkerrors.Wrap(err, "admin")
		}
	}
	if err := validateLabel(c.Label); err != nil {
		return sdkerrors.Wrap(err, "label")
	}
	found := false
	for _, v := range AllContractStatus {
		if c.Status == v {
			found = true
			break
		}
	}
	if !found || c.Status == ContractStatusUnspecified {
		return sdkerrors.Wrap(ErrInvalidMsg, "invalid status")
	}
	if c.Extension == nil {
		return nil
	}

	e, ok := c.Extension.GetCachedValue().(validatable)
	if !ok {
		return nil
	}
	if err := e.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "extension")
	}
	return nil
}

// SetExtension set new extension data. Calls `ValidateBasic() error` on non nil values when method is implemented by
// the extension.
func (c *ContractInfo) SetExtension(ext ContractInfoExtension) error {
	if ext == nil {
		c.Extension = nil
		return nil
	}
	if e, ok := ext.(validatable); ok {
		if err := e.ValidateBasic(); err != nil {
			return err
		}
	}
	any, err := codectypes.NewAnyWithValue(ext)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	c.Extension = any
	return nil
}

// ReadExtension copies the extension value to the pointer passed as argument so that there is no need to cast
// For example with a custom extension of type `MyContractDetails` it will look as following:
// 		var d MyContractDetails
//		if err := info.ReadExtension(&d); err != nil {
//			return nil, sdkerrors.Wrap(err, "extension")
//		}
func (c *ContractInfo) ReadExtension(e ContractInfoExtension) error {
	rv := reflect.ValueOf(e)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidType, "not a pointer")
	}
	if c.Extension == nil {
		return nil
	}

	cached := c.Extension.GetCachedValue()
	elem := reflect.ValueOf(cached).Elem()
	if !elem.Type().AssignableTo(rv.Elem().Type()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "extension is of type %s but argument of %s", elem.Type(), rv.Elem().Type())
	}
	rv.Elem().Set(elem)
	return nil
}

func (c ContractInfo) InitialHistory(initMsg []byte) ContractCodeHistoryEntry {
	return ContractCodeHistoryEntry{
		Operation: ContractCodeHistoryOperationTypeInit,
		CodeID:    c.CodeID,
		Updated:   c.Created,
		Msg:       initMsg,
	}
}

func (c *ContractInfo) AddMigration(ctx sdk.Context, codeID uint64, msg []byte) ContractCodeHistoryEntry {
	h := ContractCodeHistoryEntry{
		Operation: ContractCodeHistoryOperationTypeMigrate,
		CodeID:    codeID,
		Updated:   NewAbsoluteTxPosition(ctx),
		Msg:       msg,
	}
	c.CodeID = codeID
	return h
}

// ResetFromGenesis resets contracts timestamp and history.
func (c *ContractInfo) ResetFromGenesis(ctx sdk.Context) ContractCodeHistoryEntry {
	c.Created = NewAbsoluteTxPosition(ctx)
	return ContractCodeHistoryEntry{
		Operation: ContractCodeHistoryOperationTypeGenesis,
		CodeID:    c.CodeID,
		Updated:   c.Created,
	}
}

// AdminAddr convert into sdk.AccAddress or nil when not set
func (c *ContractInfo) AdminAddr() sdk.AccAddress {
	if c.Admin == "" {
		return ""
	}
	admin := sdk.AccAddress(c.Admin)
	return admin
}

// ContractInfoExtension defines the extension point for custom data to be stored with a contract info
type ContractInfoExtension interface {
	proto.Message
	String() string
}

var _ codectypes.UnpackInterfacesMessage = &ContractInfo{}

// UnpackInterfaces implements codectypes.UnpackInterfaces
func (c *ContractInfo) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var details ContractInfoExtension
	if err := unpacker.UnpackAny(c.Extension, &details); err != nil {
		return err
	}
	return codectypes.UnpackInterfaces(details, unpacker)
}

// NewAbsoluteTxPosition gets a block position from the context
func NewAbsoluteTxPosition(ctx sdk.Context) *AbsoluteTxPosition {
	// we must safely handle nil gas meters
	var index uint64
	meter := ctx.BlockGasMeter()
	if meter != nil {
		index = meter.GasConsumed()
	}
	height := ctx.BlockHeight()
	if height < 0 {
		panic(fmt.Sprintf("unsupported height: %d", height))
	}
	return &AbsoluteTxPosition{
		BlockHeight: uint64(height),
		TxIndex:     index,
	}
}

// LessThan can be used to sort
func (a *AbsoluteTxPosition) LessThan(b *AbsoluteTxPosition) bool {
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}
	return a.BlockHeight < b.BlockHeight || (a.BlockHeight == b.BlockHeight && a.TxIndex < b.TxIndex)
}

// AbsoluteTxPositionLen number of elements in byte representation
const AbsoluteTxPositionLen = 16

// Bytes encodes the object into a 16 byte representation with big endian block height and tx index.
func (a *AbsoluteTxPosition) Bytes() []byte {
	if a == nil {
		panic("object must not be nil")
	}
	r := make([]byte, AbsoluteTxPositionLen)
	copy(r[0:], sdk.Uint64ToBigEndian(a.BlockHeight))
	copy(r[8:], sdk.Uint64ToBigEndian(a.TxIndex))
	return r
}

// NewEnv initializes the environment for a contract instance
func NewEnv(ctx sdk.Context, contractAddr sdk.AccAddress) wasmvmtypes.Env {
	// safety checks before casting below
	if ctx.BlockHeight() < 0 {
		panic("Block height must never be negative")
	}
	nano := ctx.BlockTime().UnixNano()
	if nano < 1 {
		panic("Block (unix) time must never be empty or negative ")
	}
	env := wasmvmtypes.Env{
		Block: wasmvmtypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(nano),
			ChainID: ctx.ChainID(),
		},
		Contract: wasmvmtypes.ContractInfo{
			Address: contractAddr.String(),
		},
	}
	return env
}

// NewInfo initializes the MessageInfo for a contract instance
func NewInfo(creator sdk.AccAddress, deposit sdk.Coins) wasmvmtypes.MessageInfo {
	return wasmvmtypes.MessageInfo{
		Sender: creator.String(),
		Funds:  NewWasmCoins(deposit),
	}
}

// NewWasmCoins translates between Cosmos SDK coins and Wasm coins
func NewWasmCoins(cosmosCoins sdk.Coins) (wasmCoins []wasmvmtypes.Coin) {
	for _, coin := range cosmosCoins {
		wasmCoin := wasmvmtypes.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		}
		wasmCoins = append(wasmCoins, wasmCoin)
	}
	return wasmCoins
}

const CustomEventType = "wasm"
const AttributeKeyContractAddr = "contract_address"

// ParseEvents converts wasm LogAttributes into an sdk.Events
func ParseEvents(wasmOutputAttrs []wasmvmtypes.EventAttribute, contractAddr sdk.AccAddress) sdk.Events {
	// we always tag with the contract address issuing this event
	attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddr, contractAddr.String())}

	// append attributes from wasm to the sdk.Event
	for _, l := range wasmOutputAttrs {
		// and reserve the contract_address key for our use (not contract)
		if l.Key != AttributeKeyContractAddr {
			attr := sdk.NewAttribute(l.Key, l.Value)
			attrs = append(attrs, attr)
		}
	}

	// each wasm invokation always returns one sdk.Event
	return sdk.Events{sdk.NewEvent(CustomEventType, attrs...)}
}

// WasmConfig is the extra config required for wasm
type WasmConfig struct {
	SmartQueryGasLimit uint64
	// MemoryCacheSize in MiB not bytes
	MemoryCacheSize uint32
	// ContractDebugMode log what contract print
	ContractDebugMode bool
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWasmConfig() WasmConfig {
	return WasmConfig{
		SmartQueryGasLimit: defaultQueryGasLimit,
		MemoryCacheSize:    defaultMemoryCacheSize,
		ContractDebugMode:  defaultContractDebugMode,
	}
}
