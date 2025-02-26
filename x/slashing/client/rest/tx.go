package rest

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/tx"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/rest"
	"github.com/line/lbm-sdk/x/slashing/types"
)

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/slashing/validators/{validatorAddr}/unjail", NewUnjailRequestHandlerFn(clientCtx)).Methods("POST")
}

// Unjail TX body
type UnjailReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

// NewUnjailRequestHandlerFn returns an HTTP REST handler for creating a MsgUnjail
// transaction.
func NewUnjailRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Validator := vars["validatorAddr"]

		var req UnjailReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressToBytes(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		valAddr, err := sdk.ValAddressToBytes(bech32Validator)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if !bytes.Equal(fromAddr, valAddr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own validator address")
			return
		}

		msg := types.NewMsgUnjail(sdk.ValAddress(bech32Validator))
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
