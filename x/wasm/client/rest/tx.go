package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	wasmUtils "github.com/line/lbm-sdk/x/wasm/client/utils"
	"github.com/line/lbm-sdk/x/wasm/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/wasm/code", storeCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/code/{codeId}", instantiateContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/codeinit", storeCodeAndInstantiateContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/{contractAddr}", executeContractHandlerFn(cliCtx)).Methods("POST")
}

type storeCodeReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WasmBytes []byte       `json:"wasm_bytes"`
}

type instantiateContractReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Label   string       `json:"label" yaml:"label"`
	Deposit sdk.Coins    `json:"deposit" yaml:"deposit"`
	Admin   string       `json:"admin,omitempty" yaml:"admin"`
	InitMsg []byte       `json:"init_msg" yaml:"init_msg"`
}

type storeCodeAndInstantiateContractReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WasmBytes []byte       `json:"wasm_bytes"`
	Label     string       `json:"label" yaml:"label"`
	Deposit   sdk.Coins    `json:"deposit" yaml:"deposit"`
	Admin     string       `json:"admin,omitempty" yaml:"admin"`
	InitMsg   []byte       `json:"init_msg" yaml:"init_msg"`
}

type executeContractReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	ExecMsg []byte       `json:"exec_msg" yaml:"exec_msg"`
	Amount  sdk.Coins    `json:"coins" yaml:"coins"`
}

func storeCodeHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req storeCodeReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var err error
		wasm := req.WasmBytes

		// gzip the wasm file
		if wasmUtils.IsWasm(wasm) {
			wasm, err = wasmUtils.GzipIt(wasm)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if !wasmUtils.IsGzip(wasm) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid input file, use wasm binary or zip")
			return
		}

		// build and sign the transaction, then broadcast to Tendermint
		msg := types.MsgStoreCode{
			Sender:       req.BaseReq.From,
			WASMByteCode: wasm,
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func instantiateContractHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req instantiateContractReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}
		vars := mux.Vars(r)
		codeIDVar := vars["codeId"]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// get the id of the code to instantiate
		codeID, err := strconv.ParseUint(codeIDVar, 10, 64)
		if err != nil {
			return
		}

		msg := types.MsgInstantiateContract{
			Sender:  req.BaseReq.From,
			CodeID:  codeID,
			Label:   req.Label,
			Funds:   req.Deposit,
			InitMsg: req.InitMsg,
			Admin:   req.Admin,
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func storeCodeAndInstantiateContractHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req storeCodeAndInstantiateContractReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var err error
		wasm := req.WasmBytes

		// gzip the wasm file
		if wasmUtils.IsWasm(wasm) {
			wasm, err = wasmUtils.GzipIt(wasm)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if !wasmUtils.IsGzip(wasm) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid input file, use wasm binary or zip")
			return
		}

		// build and sign the transaction, then broadcast to Tendermint
		msg := types.MsgStoreCodeAndInstantiateContract{
			Sender:       req.BaseReq.From,
			WASMByteCode: wasm,
			Label:        req.Label,
			Funds:        req.Deposit,
			InitMsg:      req.InitMsg,
			Admin:        req.Admin,
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func executeContractHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req executeContractReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}
		vars := mux.Vars(r)
		contractAddr := vars["contractAddr"]

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgExecuteContract{
			Sender:   req.BaseReq.From,
			Contract: contractAddr,
			Msg:      req.ExecMsg,
			Funds:    req.Amount,
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}
