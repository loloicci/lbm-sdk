package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/wasm/types"
	abci "github.com/line/ostracon/abci/types"
)

// ValidatorSetSource is a subset of the staking keeper
type ValidatorSetSource interface {
	ApplyAndReturnValidatorSetUpdates(sdk.Context) (updates []abci.ValidatorUpdate, err error)
}

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper *Keeper, data types.GenesisState, stakingKeeper ValidatorSetSource, msgHandler sdk.Handler) ([]abci.ValidatorUpdate, error) {
	contractKeeper := NewGovPermissionKeeper(keeper)
	keeper.setParams(ctx, data.Params)
	var maxCodeID uint64
	for i, code := range data.Codes {
		err := keeper.importCode(ctx, code.CodeID, code.CodeInfo, code.CodeBytes)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "code %d with id: %d", i, code.CodeID)
		}
		if code.CodeID > maxCodeID {
			maxCodeID = code.CodeID
		}
		if code.Pinned {
			if err := contractKeeper.PinCode(ctx, code.CodeID); err != nil {
				return nil, sdkerrors.Wrapf(err, "contract number %d", i)
			}
		}
	}

	var maxContractID int
	for i, contract := range data.Contracts {
		err := sdk.ValidateAccAddress(contract.ContractAddress)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "address in contract number %d", i)
		}
		err = keeper.importContract(ctx, sdk.AccAddress(contract.ContractAddress), &contract.ContractInfo,
			contract.ContractState)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "contract number %d", i)
		}
		maxContractID = i + 1 // not ideal but max(contractID) is not persisted otherwise
	}

	for i, seq := range data.Sequences {
		err := keeper.importAutoIncrementID(ctx, seq.IDKey, seq.Value)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "sequence number %d", i)
		}
	}

	// sanity check seq values
	if keeper.peekAutoIncrementID(ctx, types.KeyLastCodeID) <= maxCodeID {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "seq %s must be greater %d ", string(types.KeyLastCodeID), maxCodeID)
	}
	if keeper.peekAutoIncrementID(ctx, types.KeyLastInstanceID) <= uint64(maxContractID) {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "seq %s must be greater %d ", string(types.KeyLastInstanceID), maxContractID)
	}

	if len(data.GenMsgs) == 0 {
		return nil, nil
	}
	for _, genTx := range data.GenMsgs {
		msg := genTx.AsMsg()
		if msg == nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "unknown message")
		}
		_, err := msgHandler(ctx, msg)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "genesis")
		}
	}
	return stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper *Keeper) *types.GenesisState {
	var genState types.GenesisState

	genState.Params = keeper.GetParams(ctx)

	keeper.IterateCodeInfos(ctx, func(codeID uint64, info types.CodeInfo) bool {
		bytecode, err := keeper.GetByteCode(ctx, codeID)
		if err != nil {
			panic(err)
		}
		genState.Codes = append(genState.Codes, types.Code{
			CodeID:    codeID,
			CodeInfo:  info,
			CodeBytes: bytecode,
			Pinned:    keeper.IsPinnedCode(ctx, codeID),
		})
		return false
	})

	keeper.IterateContractInfo(ctx, func(addr sdk.AccAddress, contract types.ContractInfo) bool {
		contractStateIterator := keeper.GetContractState(ctx, addr)
		var state []types.Model
		for ; contractStateIterator.Valid(); contractStateIterator.Next() {
			m := types.Model{
				Key:   contractStateIterator.Key(),
				Value: contractStateIterator.Value(),
			}
			state = append(state, m)
		}
		// redact contract info
		contract.Created = nil
		genState.Contracts = append(genState.Contracts, types.Contract{
			ContractAddress: addr.String(),
			ContractInfo:    contract,
			ContractState:   state,
		})

		return false
	})

	for _, k := range [][]byte{types.KeyLastCodeID, types.KeyLastInstanceID} {
		genState.Sequences = append(genState.Sequences, types.Sequence{
			IDKey: k,
			Value: keeper.peekAutoIncrementID(ctx, k),
		})
	}

	return &genState
}
