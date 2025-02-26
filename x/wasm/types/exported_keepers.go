package types

import (
	sdk "github.com/line/lbm-sdk/types"
	capabilitytypes "github.com/line/lbm-sdk/x/capability/types"
	types2 "github.com/line/wasmvm/types"
)

// ViewKeeper provides read only operations
type ViewKeeper interface {
	GetContractHistory(ctx sdk.Context, contractAddr sdk.AccAddress) []ContractCodeHistoryEntry
	QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
	QueryRaw(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) []byte
	HasContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool
	GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) *ContractInfo
	IterateContractInfo(ctx sdk.Context, cb func(sdk.AccAddress, ContractInfo) bool)
	IterateContractsByCode(ctx sdk.Context, codeID uint64, cb func(address sdk.AccAddress) bool)
	GetContractState(ctx sdk.Context, contractAddress sdk.AccAddress) sdk.Iterator
	GetCodeInfo(ctx sdk.Context, codeID uint64) *CodeInfo
	IterateCodeInfos(ctx sdk.Context, cb func(uint64, CodeInfo) bool)
	GetByteCode(ctx sdk.Context, codeID uint64) ([]byte, error)
	IsPinnedCode(ctx sdk.Context, codeID uint64) bool
}

// ContractOpsKeeper contains mutable operations on a contract.
type ContractOpsKeeper interface {
	// Create uploads and compiles a WASM contract, returning a short identifier for the contract
	Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte, source string, builder string, instantiateAccess *AccessConfig) (codeID uint64, err error)

	// Instantiate creates an instance of a WASM contract
	Instantiate(ctx sdk.Context, codeID uint64, creator, admin sdk.AccAddress, initMsg []byte, label string, deposit sdk.Coins) (sdk.AccAddress, []byte, error)

	// Execute executes the contract instance
	Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) (*sdk.Result, error)

	// Migrate allows to upgrade a contract to a new code with data migration.
	Migrate(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newCodeID uint64, msg []byte) (*sdk.Result, error)

	// UpdateContractAdmin sets the admin value on the ContractInfo. It must be a valid address (use ClearContractAdmin to remove it)
	UpdateContractAdmin(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newAdmin sdk.AccAddress) error

	// ClearContractAdmin sets the admin value on the ContractInfo to nil, to disable further migrations/ updates.
	ClearContractAdmin(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress) error

	// PinCode pins the wasm contract in wasmvm cache
	PinCode(ctx sdk.Context, codeID uint64) error

	// UnpinCode removes the wasm contract from wasmvm cache
	UnpinCode(ctx sdk.Context, codeID uint64) error

	// SetContractInfoExtension updates the extension point data that is stored with the contract info
	SetContractInfoExtension(ctx sdk.Context, contract sdk.AccAddress, extra ContractInfoExtension) error

	// UpdateContractStatus sets a new status of the contract on the ContractInfo.
	UpdateContractStatus(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, status ContractStatus) error
}

// IBCContractKeeper IBC lifecycle event handler
type IBCContractKeeper interface {
	OnOpenChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		channel types2.IBCChannel,
	) error
	OnConnectChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		channel types2.IBCChannel,
	) error
	OnCloseChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		channel types2.IBCChannel,
	) error
	OnRecvPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		packet types2.IBCPacket,
	) ([]byte, error)
	OnAckPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		acknowledgement types2.IBCAcknowledgement,
	) error
	OnTimeoutPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		packet types2.IBCPacket,
	) error
	// ClaimCapability allows the transfer module to claim a capability
	//that IBC module passes to it
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
	// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
	AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool
}
