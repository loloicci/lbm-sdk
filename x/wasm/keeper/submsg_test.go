package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/line/lbm-sdk/types"
	wasmvmtypes "github.com/line/wasmvm/types"
	"github.com/stretchr/testify/require"
)

// test handing of submessages, very closely related to the reflect_test

// Try a simple send, no gas limit to for a sanity check before trying table tests
func TestDispatchSubMsgSuccessCase(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, ReflectFeatures, nil, nil)
	accKeeper, keeper, bankKeeper := keepers.AccountKeeper, keepers.WasmKeeper, keepers.BankKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))

	creator := createFakeFundedAccount(t, ctx, accKeeper, bankKeeper, deposit)
	creatorBalance := deposit.Sub(contractStart)
	_, _, fred := keyPubAddr()

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keepers.ContractKeeper.Create(ctx, creator, reflectCode, "", "", nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, "", []byte("{}"), "reflect contract 1", contractStart)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// check some account values
	checkAccount(t, ctx, accKeeper, bankKeeper, contractAddr, contractStart)
	checkAccount(t, ctx, accKeeper, bankKeeper, creator, creatorBalance)
	checkAccount(t, ctx, accKeeper, bankKeeper, fred, nil)

	// creator can send contract's tokens to fred (using SendMsg)
	msg := wasmvmtypes.CosmosMsg{
		Bank: &wasmvmtypes.BankMsg{
			Send: &wasmvmtypes.SendMsg{
				ToAddress: fred.String(),
				Amount: []wasmvmtypes.Coin{{
					Denom:  "denom",
					Amount: "15000",
				}},
			},
		},
	}
	reflectSend := ReflectHandleMsg{
		ReflectSubCall: &reflectSubPayload{
			Msgs: []wasmvmtypes.SubMsg{{
				ID:      7,
				Msg:     msg,
				ReplyOn: wasmvmtypes.ReplyAlways,
			}},
		},
	}
	reflectSendBz, err := json.Marshal(reflectSend)
	require.NoError(t, err)
	_, err = keepers.ContractKeeper.Execute(ctx, contractAddr, creator, reflectSendBz, nil)
	require.NoError(t, err)

	// fred got coins
	checkAccount(t, ctx, accKeeper, bankKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin("denom", 15000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, bankKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin("denom", 25000)))
	checkAccount(t, ctx, accKeeper, bankKeeper, creator, creatorBalance)

	// query the reflect state to ensure the result was stored
	query := ReflectQueryMsg{
		SubCallResult: &SubCall{ID: 7},
	}
	queryBz, err := json.Marshal(query)
	require.NoError(t, err)
	queryRes, err := keeper.QuerySmart(ctx, contractAddr, queryBz)
	require.NoError(t, err)

	var res wasmvmtypes.Reply
	err = json.Unmarshal(queryRes, &res)
	require.NoError(t, err)
	assert.Equal(t, uint64(7), res.ID)
	assert.Empty(t, res.Result.Err)
	require.NotNil(t, res.Result.Ok)
	sub := res.Result.Ok
	assert.Empty(t, sub.Data)
	require.Len(t, sub.Events, 3)

	transfer := sub.Events[0]
	assert.Equal(t, "transfer", transfer.Type)
	assert.Equal(t, wasmvmtypes.EventAttribute{
		Key:   "recipient",
		Value: fred.String(),
	}, transfer.Attributes[0])

	sender := sub.Events[1]
	assert.Equal(t, "message", sender.Type)
	assert.Equal(t, wasmvmtypes.EventAttribute{
		Key:   "sender",
		Value: contractAddr.String(),
	}, sender.Attributes[0])

	// where does this come from?
	module := sub.Events[2]
	assert.Equal(t, "message", module.Type)
	assert.Equal(t, wasmvmtypes.EventAttribute{
		Key:   "module",
		Value: "bank",
	}, module.Attributes[0])

}

func TestDispatchSubMsgErrorHandling(t *testing.T) {
	fundedDenom := "funds"
	fundedAmount := 1_000_000
	ctxGasLimit := uint64(1_000_000)
	subGasLimit := uint64(300_000)

	// prep - create one chain and upload the code
	ctx, keepers := CreateTestInput(t, false, ReflectFeatures, nil, nil)
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	ctx = ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())
	accKeeper, keeper, bankKeeper := keepers.AccountKeeper, keepers.WasmKeeper, keepers.BankKeeper
	contractStart := sdk.NewCoins(sdk.NewInt64Coin(fundedDenom, int64(fundedAmount)))
	uploader := createFakeFundedAccount(t, ctx, accKeeper, bankKeeper, contractStart.Add(contractStart...))

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	reflectID, err := keepers.ContractKeeper.Create(ctx, uploader, reflectCode, "", "", nil)
	require.NoError(t, err)

	// create hackatom contract for testing (for infinite loop)
	hackatomCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)
	hackatomID, err := keepers.ContractKeeper.Create(ctx, uploader, hackatomCode, "", "", nil)
	require.NoError(t, err)
	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	hackatomAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, hackatomID, uploader, "", initMsgBz, "hackatom demo", contractStart)
	require.NoError(t, err)

	validBankSend := func(contract, emptyAccount string) wasmvmtypes.CosmosMsg {
		return wasmvmtypes.CosmosMsg{
			Bank: &wasmvmtypes.BankMsg{
				Send: &wasmvmtypes.SendMsg{
					ToAddress: emptyAccount,
					Amount: []wasmvmtypes.Coin{{
						Denom:  fundedDenom,
						Amount: strconv.Itoa(fundedAmount / 2),
					}},
				},
			},
		}
	}

	invalidBankSend := func(contract, emptyAccount string) wasmvmtypes.CosmosMsg {
		return wasmvmtypes.CosmosMsg{
			Bank: &wasmvmtypes.BankMsg{
				Send: &wasmvmtypes.SendMsg{
					ToAddress: emptyAccount,
					Amount: []wasmvmtypes.Coin{{
						Denom:  fundedDenom,
						Amount: strconv.Itoa(fundedAmount * 2),
					}},
				},
			},
		}
	}

	infiniteLoop := func(contract, emptyAccount string) wasmvmtypes.CosmosMsg {
		return wasmvmtypes.CosmosMsg{
			Wasm: &wasmvmtypes.WasmMsg{
				Execute: &wasmvmtypes.ExecuteMsg{
					ContractAddr: hackatomAddr.String(),
					Msg:          []byte(`{"cpu_loop":{}}`),
				},
			},
		}
	}

	instantiateContract := func(contract, emptyAccount string) wasmvmtypes.CosmosMsg {
		return wasmvmtypes.CosmosMsg{
			Wasm: &wasmvmtypes.WasmMsg{
				Instantiate: &wasmvmtypes.InstantiateMsg{
					CodeID: reflectID,
					Msg:    []byte("{}"),
					Label:  "subcall reflect",
				},
			},
		}
	}

	type assertion func(t *testing.T, ctx sdk.Context, contract, emptyAccount string, response wasmvmtypes.SubcallResult)

	assertReturnedEvents := func(expectedEvents int) assertion {
		return func(t *testing.T, ctx sdk.Context, contract, emptyAccount string, response wasmvmtypes.SubcallResult) {
			assert.Len(t, response.Ok.Events, expectedEvents)
		}
	}

	assertGasUsed := func(minGas, maxGas uint64) assertion {
		return func(t *testing.T, ctx sdk.Context, contract, emptyAccount string, response wasmvmtypes.SubcallResult) {
			gasUsed := ctx.GasMeter().GasConsumed()
			assert.True(t, gasUsed >= minGas, "Used %d gas (less than expected %d)", gasUsed, minGas)
			assert.True(t, gasUsed <= maxGas, "Used %d gas (more than expected %d)", gasUsed, maxGas)
		}
	}

	assertErrorString := func(shouldContain string) assertion {
		return func(t *testing.T, ctx sdk.Context, contract, emptyAccount string, response wasmvmtypes.SubcallResult) {
			assert.Contains(t, response.Err, shouldContain)
		}
	}

	assertGotContractAddr := func(t *testing.T, ctx sdk.Context, contract, emptyAccount string, response wasmvmtypes.SubcallResult) {
		// should get the events emitted on new contract
		event := response.Ok.Events[0]
		assert.Equal(t, event.Type, "wasm")
		assert.Equal(t, event.Attributes[0].Key, "contract_address")
		eventAddr := event.Attributes[0].Value
		assert.NotEqual(t, contract, eventAddr)

		// data field is the raw canonical address
		// QUESTION: why not types.MsgInstantiateContractResponse? difference between calling Router and Service?
		resAddr := sdk.AccAddress(response.Ok.Data)
		assert.Equal(t, eventAddr, resAddr.String())
	}

	cases := map[string]struct {
		submsgID uint64
		// we will generate message from the
		msg      func(contract, emptyAccount string) wasmvmtypes.CosmosMsg
		gasLimit *uint64

		// true if we expect this to throw out of gas panic
		isOutOfGasPanic bool
		// true if we expect this execute to return an error (can be false when submessage errors)
		executeError bool
		// true if we expect submessage to return an error (but execute to return success)
		subMsgError bool
		// make assertions after dispatch
		resultAssertions []assertion
	}{
		"send tokens": {
			submsgID: 5,
			msg:      validBankSend,
			// note we charge another 40k for the reply call
			resultAssertions: []assertion{assertReturnedEvents(3), assertGasUsed(117000, 130000)},
		},
		"not enough tokens": {
			submsgID:    6,
			msg:         invalidBankSend,
			subMsgError: true,
			// uses less gas than the send tokens (cost of bank transfer)
			resultAssertions: []assertion{assertGasUsed(96000, 110000), assertErrorString("insufficient funds")},
		},
		"out of gas panic with no gas limit": {
			submsgID:        7,
			msg:             infiniteLoop,
			isOutOfGasPanic: true,
		},

		"send tokens with limit": {
			submsgID: 15,
			msg:      validBankSend,
			gasLimit: &subGasLimit,
			// uses same gas as call without limit
			resultAssertions: []assertion{assertReturnedEvents(3), assertGasUsed(117000, 130000)},
		},
		"not enough tokens with limit": {
			submsgID:    16,
			msg:         invalidBankSend,
			subMsgError: true,
			gasLimit:    &subGasLimit,
			// uses same gas as call without limit
			resultAssertions: []assertion{assertGasUsed(96000, 110000), assertErrorString("insufficient funds")},
		},
		"out of gas caught with gas limit": {
			submsgID:    17,
			msg:         infiniteLoop,
			subMsgError: true,
			gasLimit:    &subGasLimit,
			// uses all the subGasLimit, plus the 92k or so for the main contract
			resultAssertions: []assertion{assertGasUsed(subGasLimit+92000, subGasLimit+106000), assertErrorString("out of gas")},
		},

		"instantiate contract gets address in data and events": {
			submsgID:         21,
			msg:              instantiateContract,
			resultAssertions: []assertion{assertReturnedEvents(1), assertGotContractAddr},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			creator := createFakeFundedAccount(t, ctx, accKeeper, bankKeeper, contractStart)
			_, _, empty := keyPubAddr()

			contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, reflectID, creator, "", []byte("{}"), fmt.Sprintf("contract %s", name), contractStart)
			require.NoError(t, err)

			msg := tc.msg(contractAddr.String(), empty.String())
			reflectSend := ReflectHandleMsg{
				ReflectSubCall: &reflectSubPayload{
					Msgs: []wasmvmtypes.SubMsg{{
						ID:       tc.submsgID,
						Msg:      msg,
						GasLimit: tc.gasLimit,
						ReplyOn:  wasmvmtypes.ReplyAlways,
					}},
				},
			}
			reflectSendBz, err := json.Marshal(reflectSend)
			require.NoError(t, err)

			execCtx := ctx.WithGasMeter(sdk.NewGasMeter(ctxGasLimit))
			defer func() {
				if tc.isOutOfGasPanic {
					r := recover()
					require.NotNil(t, r, "expected panic")
					if _, ok := r.(sdk.ErrorOutOfGas); !ok {
						t.Fatalf("Expected OutOfGas panic, got: %#v\n", r)
					}
				}
			}()
			_, err = keepers.ContractKeeper.Execute(execCtx, contractAddr, creator, reflectSendBz, nil)

			if tc.executeError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// query the reply
				query := ReflectQueryMsg{
					SubCallResult: &SubCall{ID: tc.submsgID},
				}
				queryBz, err := json.Marshal(query)
				require.NoError(t, err)
				queryRes, err := keeper.QuerySmart(ctx, contractAddr, queryBz)
				require.NoError(t, err)
				var res wasmvmtypes.Reply
				err = json.Unmarshal(queryRes, &res)
				require.NoError(t, err)
				assert.Equal(t, tc.submsgID, res.ID)

				if tc.subMsgError {
					require.NotEmpty(t, res.Result.Err)
					require.Nil(t, res.Result.Ok)
				} else {
					require.Empty(t, res.Result.Err)
					require.NotNil(t, res.Result.Ok)
				}

				for _, assertion := range tc.resultAssertions {
					assertion(t, execCtx, contractAddr.String(), empty.String(), res.Result)
				}

			}
		})
	}
}

// Test an error case, where the Encoded doesn't return any sdk.Msg and we trigger(ed) a null pointer exception.
// This occurs with the IBC encoder. Test this.
func TestDispatchSubMsgEncodeToNoSdkMsg(t *testing.T) {
	// fake out the bank handle to return success with no data
	nilEncoder := func(sender sdk.AccAddress, msg *wasmvmtypes.BankMsg) ([]sdk.Msg, error) {
		return nil, nil
	}
	customEncoders := &MessageEncoders{
		Bank: nilEncoder,
	}

	ctx, keepers := CreateTestInput(t, false, ReflectFeatures, nil, nil, WithMessageHandler(NewSDKMessageHandler(nil, nil, customEncoders)))
	accKeeper, keeper, bankKeeper := keepers.AccountKeeper, keepers.WasmKeeper, keepers.BankKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))

	creator := createFakeFundedAccount(t, ctx, accKeeper, bankKeeper, deposit)
	_, _, fred := keyPubAddr()

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keepers.ContractKeeper.Create(ctx, creator, reflectCode, "", "", nil)
	require.NoError(t, err)

	// creator instantiates a contract and gives it tokens
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, "", []byte("{}"), "reflect contract 1", contractStart)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// creator can send contract's tokens to fred (using SendMsg)
	msg := wasmvmtypes.CosmosMsg{
		Bank: &wasmvmtypes.BankMsg{
			Send: &wasmvmtypes.SendMsg{
				ToAddress: fred.String(),
				Amount: []wasmvmtypes.Coin{{
					Denom:  "denom",
					Amount: "15000",
				}},
			},
		},
	}
	reflectSend := ReflectHandleMsg{
		ReflectSubCall: &reflectSubPayload{
			Msgs: []wasmvmtypes.SubMsg{{
				ID:      7,
				Msg:     msg,
				ReplyOn: wasmvmtypes.ReplyAlways,
			}},
		},
	}
	reflectSendBz, err := json.Marshal(reflectSend)
	require.NoError(t, err)
	_, err = keepers.ContractKeeper.Execute(ctx, contractAddr, creator, reflectSendBz, nil)
	require.NoError(t, err)

	// query the reflect state to ensure the result was stored
	query := ReflectQueryMsg{
		SubCallResult: &SubCall{ID: 7},
	}
	queryBz, err := json.Marshal(query)
	require.NoError(t, err)
	queryRes, err := keeper.QuerySmart(ctx, contractAddr, queryBz)
	require.NoError(t, err)

	var res wasmvmtypes.Reply
	err = json.Unmarshal(queryRes, &res)
	require.NoError(t, err)
	assert.Equal(t, uint64(7), res.ID)
	assert.Empty(t, res.Result.Err)
	require.NotNil(t, res.Result.Ok)
	sub := res.Result.Ok
	assert.Empty(t, sub.Data)
	require.Len(t, sub.Events, 0)
}

// Try a simple send, no gas limit to for a sanity check before trying table tests
func TestDispatchSubMsgConditionalReplyOn(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, ReflectFeatures, nil, nil)
	accKeeper, keeper, bankKeeper := keepers.AccountKeeper, keepers.WasmKeeper, keepers.BankKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))

	creator := createFakeFundedAccount(t, ctx, accKeeper, bankKeeper, deposit)
	_, _, fred := keyPubAddr()

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keepers.ContractKeeper.Create(ctx, creator, reflectCode, "", "", nil)
	require.NoError(t, err)

	// creator instantiates a contract and gives it tokens
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, "", []byte("{}"), "reflect contract 1", contractStart)
	require.NoError(t, err)

	goodSend := wasmvmtypes.CosmosMsg{
		Bank: &wasmvmtypes.BankMsg{
			Send: &wasmvmtypes.SendMsg{
				ToAddress: fred.String(),
				Amount: []wasmvmtypes.Coin{{
					Denom:  "denom",
					Amount: "1000",
				}},
			},
		},
	}
	failSend := wasmvmtypes.CosmosMsg{
		Bank: &wasmvmtypes.BankMsg{
			Send: &wasmvmtypes.SendMsg{
				ToAddress: fred.String(),
				Amount: []wasmvmtypes.Coin{{
					Denom:  "no-such-token",
					Amount: "777777",
				}},
			},
		},
	}

	cases := map[string]struct {
		// true for wasmvmtypes.ReplySuccess, false for wasmvmtypes.ReplyError
		replyOnSuccess bool
		msg            wasmvmtypes.CosmosMsg
		// true if the call should return an error (it wasn't handled)
		expectError bool
		// true if the reflect contract wrote the response (success or error) - it was captured
		writeResult bool
	}{
		"all good, reply success": {
			replyOnSuccess: true,
			msg:            goodSend,
			expectError:    false,
			writeResult:    true,
		},
		"all good, reply error": {
			replyOnSuccess: false,
			msg:            goodSend,
			expectError:    false,
			writeResult:    false,
		},
		"bad msg, reply success": {
			replyOnSuccess: true,
			msg:            failSend,
			expectError:    true,
			writeResult:    false,
		},
		"bad msg, reply error": {
			replyOnSuccess: false,
			msg:            failSend,
			expectError:    false,
			writeResult:    true,
		},
	}

	var id uint64 = 0
	for name, tc := range cases {
		id++
		t.Run(name, func(t *testing.T) {
			subMsg := wasmvmtypes.SubMsg{
				ID:      id,
				Msg:     tc.msg,
				ReplyOn: wasmvmtypes.ReplySuccess,
			}
			if !tc.replyOnSuccess {
				subMsg.ReplyOn = wasmvmtypes.ReplyError
			}

			reflectSend := ReflectHandleMsg{
				ReflectSubCall: &reflectSubPayload{
					Msgs: []wasmvmtypes.SubMsg{subMsg},
				},
			}
			reflectSendBz, err := json.Marshal(reflectSend)
			require.NoError(t, err)
			_, err = keepers.ContractKeeper.Execute(ctx, contractAddr, creator, reflectSendBz, nil)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// query the reflect state to check if the result was stored
			query := ReflectQueryMsg{
				SubCallResult: &SubCall{ID: id},
			}
			queryBz, err := json.Marshal(query)
			require.NoError(t, err)
			queryRes, err := keeper.QuerySmart(ctx, contractAddr, queryBz)
			if tc.writeResult {
				// we got some data for this call
				require.NoError(t, err)
				var res wasmvmtypes.Reply
				err = json.Unmarshal(queryRes, &res)
				require.NoError(t, err)
				require.Equal(t, id, res.ID)
			} else {
				// nothing should be there -> error
				require.Error(t, err)
			}
		})
	}
}
