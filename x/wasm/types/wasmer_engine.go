package types

import (
	wasmvm "github.com/line/wasmvm"
	wasmvmtypes "github.com/line/wasmvm/types"
)

// WasmerEngine defines the WASM contract runtime engine.
type WasmerEngine interface {

	// Create will compile the wasm code, and store the resulting pre-compile
	// as well as the original code. Both can be referenced later via CodeID
	// This must be done one time for given code, after which it can be
	// instatitated many times, and each instance called many times.
	//
	// For example, the code for all ERC-20 contracts should be the same.
	// This function stores the code for that contract only once, but it can
	// be instantiated with custom inputs in the future.
	Create(code wasmvm.WasmCode) (wasmvm.Checksum, error)

	// AnalyzeCode will statically analyze the code.
	// Currently just reports if it exposes all IBC entry points.
	AnalyzeCode(checksum wasmvm.Checksum) (*wasmvmtypes.AnalysisReport, error)

	// Instantiate will create a new contract based on the given codeID.
	// We can set the initMsg (contract "genesis") here, and it then receives
	// an account and address and can be invoked (Execute) many times.
	//
	// Storage should be set with a PrefixedKVStore that this code can safely access.
	//
	// Under the hood, we may recompile the wasm, use a cached native compile, or even use a cached instance
	// for performance.
	Instantiate(
		code wasmvm.Checksum,
		env wasmvmtypes.Env,
		info wasmvmtypes.MessageInfo,
		initMsg []byte,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.Response, uint64, error)

	// Execute calls a given contract. Since the only difference between contracts with the same CodeID is the
	// data in their local storage, and their address in the outside world, we need no ContractID here.
	// (That is a detail for the external, sdk-facing, side).
	//
	// The caller is responsible for passing the correct `store` (which must have been initialized exactly once),
	// and setting the env with relevant info on this instance (address, balance, etc)
	Execute(
		code wasmvm.Checksum,
		env wasmvmtypes.Env,
		info wasmvmtypes.MessageInfo,
		executeMsg []byte,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.Response, uint64, error)

	// Query allows a client to execute a contract-specific query. If the result is not empty, it should be
	// valid json-encoded data to return to the client.
	// The meaning of path and data can be determined by the code. Path is the suffix of the abci.QueryRequest.Path
	Query(
		code wasmvm.Checksum,
		env wasmvmtypes.Env,
		queryMsg []byte,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) ([]byte, uint64, error)

	// Migrate will migrate an existing contract to a new code binary.
	// This takes storage of the data from the original contract and the CodeID of the new contract that should
	// replace it. This allows it to run a migration step if needed, or return an error if unable to migrate
	// the given data.
	//
	// MigrateMsg has some data on how to perform the migration.
	Migrate(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		migrateMsg []byte,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.Response, uint64, error)

	// Sudo runs an existing contract in read/write mode (like Execute), but is never exposed to external callers
	// (either transactions or government proposals), but can only be called by other native Go modules directly.
	//
	// This allows a contract to expose custom "super user" functions or priviledged operations that can be
	// deeply integrated with native modules.
	Sudo(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		sudoMsg []byte,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.Response, uint64, error)

	// Reply is called on the original dispatching contract after running a submessage
	Reply(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		reply wasmvmtypes.Reply,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.Response, uint64, error)

	// GetCode will load the original wasm code for the given code id.
	// This will only succeed if that code id was previously returned from
	// a call to Create.
	//
	// This can be used so that the (short) code id (hash) is stored in the iavl tree
	// and the larger binary blobs (wasm and pre-compiles) are all managed by the
	// rust library
	GetCode(code wasmvm.Checksum) (wasmvm.WasmCode, error)

	// Cleanup should be called when no longer using this to free resources on the rust-side
	Cleanup()

	// IBCChannelOpen is available on IBC-enabled contracts and is a hook to call into
	// during the handshake pahse
	IBCChannelOpen(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		channel wasmvmtypes.IBCChannel,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (uint64, error)

	// IBCChannelConnect is available on IBC-enabled contracts and is a hook to call into
	// during the handshake pahse
	IBCChannelConnect(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		channel wasmvmtypes.IBCChannel,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)

	// IBCChannelClose is available on IBC-enabled contracts and is a hook to call into
	// at the end of the channel lifetime
	IBCChannelClose(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		channel wasmvmtypes.IBCChannel,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)

	// IBCPacketReceive is available on IBC-enabled contracts and is called when an incoming
	// packet is received on a channel belonging to this contract
	IBCPacketReceive(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		packet wasmvmtypes.IBCPacket,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCReceiveResponse, uint64, error)
	// IBCPacketAck is available on IBC-enabled contracts and is called when an
	// the response for an outgoing packet (previously sent by this contract)
	// is received
	IBCPacketAck(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		ack wasmvmtypes.IBCAcknowledgement,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)

	// IBCPacketTimeout is available on IBC-enabled contracts and is called when an
	// outgoing packet (previously sent by this contract) will provably never be executed.
	// Usually handled like ack returning an error
	IBCPacketTimeout(
		codeID wasmvm.Checksum,
		env wasmvmtypes.Env,
		packet wasmvmtypes.IBCPacket,
		store wasmvm.KVStore,
		goapi wasmvm.GoAPI,
		querier wasmvm.Querier,
		gasMeter wasmvm.GasMeter,
		gasLimit uint64,
	) (*wasmvmtypes.IBCBasicResponse, uint64, error)

	// Pin pins a code to an in-memory cache, such that is
	// always loaded quickly when executed.
	// Pin is idempotent.
	Pin(checksum wasmvm.Checksum) error

	// Unpin removes the guarantee of a contract to be pinned (see Pin).
	// After calling this, the code may or may not remain in memory depending on
	// the implementor's choice.
	// Unpin is idempotent.
	Unpin(checksum wasmvm.Checksum) error

	// GetMetrics some internal metrics for monitoring purposes.
	GetMetrics() (*wasmvmtypes.Metrics, error)
}
