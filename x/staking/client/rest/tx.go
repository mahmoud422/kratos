package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/KuChain-io/kuchain/chain/client/txutil"
	chaintype "github.com/KuChain-io/kuchain/chain/types"
	rest "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	ctx := txutil.NewKuCLICtx(cliCtx)

	r.HandleFunc(
		"/staking/delegations",
		postDelegationsHandlerFn(ctx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/unbonding_delegations",
		postUnbondingDelegationsHandlerFn(ctx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/redelegations",
		postRedelegationsHandlerFn(ctx),
	).Methods("POST")
}

type (
	// DelegateRequest defines the properties of a delegation request's body.
	DelegateRequest struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		DelegatorAcc string       `json:"delegator_acc" yaml:"delegator_acc"`
		ValidatorAcc string       `json:"validator_acc" yaml:"validator_acc"`
		Amount       string       `json:"amount" yaml:"amount"`
	}

	// RedelegateRequest defines the properties of a redelegate request's body.
	RedelegateRequest struct {
		BaseReq         rest.BaseReq `json:"base_req" yaml:"base_req"`
		DelegatorAcc    string       `json:"delegator_acc" yaml:"delegator_acc"`
		ValidatorSrcAcc string       `json:"validator_src_acc" yaml:"validator_src_acc"`
		ValidatorDstAcc string       `json:"validator_dst_acc" yaml:"validator_dst_acc"`
		Amount          string       `json:"amount" yaml:"amount"`
	}

	// UndelegateRequest defines the properties of a undelegate request's body.
	UndelegateRequest struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		DelegatorAcc string       `json:"delegator_acc" yaml:"delegator_acc"`
		ValidatorAcc string       `json:"validator_acc" yaml:"validator_acc"`
		Amount       string       `json:"amount" yaml:"amount"`
	}
)

func postDelegationsHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DelegateRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		amount, err := sdk.ParseCoin(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("amount parse error, %v", err))
			return
		}

		delAccountID, err := chaintype.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("delegate accountID error, %v", err))
			return
		}
		valAccountID, err := chaintype.NewAccountIDFromStr(req.ValidatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validator accountID error, %v", err))
			return
		}

		delAccAddress, err := txutil.QueryAccountAuth(cliCtx, delAccountID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", delAccountID, err))
			return
		}

		msg := types.NewKuMsgDelegate(delAccAddress, delAccountID, valAccountID, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postRedelegationsHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RedelegateRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delAccountID, err := chaintype.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("delegate acccount error, %v", err))
			return
		}

		valSrcAccID, err := chaintype.NewAccountIDFromStr(req.ValidatorSrcAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("src-validator error, %v", err))
			return
		}

		valDstAccID, err := chaintype.NewAccountIDFromStr(req.ValidatorDstAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("dst-validator error, %v", err))
			return
		}

		amount, err := sdk.ParseCoin(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		delAccAddress, err := txutil.QueryAccountAuth(cliCtx, delAccountID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", delAccountID, err))
			return
		}

		msg := types.NewKuMsgRedelegate(delAccAddress, delAccountID, valSrcAccID, valDstAccID, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postUnbondingDelegationsHandlerFn(cliCtx txutil.KuCLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UndelegateRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delAccountID, err := chaintype.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("delegate account id error, %v", err))
			return
		}

		amount, err := sdk.ParseCoin(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("parse amount error, %v", err))
			return
		}

		valAddr, err := chaintype.NewAccountIDFromStr(req.ValidatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("val account id error, %v", err))
			return
		}
		delAccAddress, err := txutil.QueryAccountAuth(cliCtx, delAccountID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("query account %s auth error, %v", delAccountID, err))
			return
		}

		msg := types.NewKuMsgUnbond(delAccAddress, delAccountID, valAddr, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
