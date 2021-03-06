package rest

import (
	"net/http"

	"github.com/KuChain-io/kuchain/chain/client/txutil"
	chainType "github.com/KuChain-io/kuchain/chain/types"
	rest "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/distribution/client/common"
	"github.com/KuChain-io/kuchain/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Withdraw all delegator rewards
	r.HandleFunc(
		"/distribution/delegators/rewards",
		withdrawDelegatorRewardsHandlerFn(cliCtx, queryRoute),
	).Methods("POST")

	// Withdraw delegation rewards
	r.HandleFunc(
		"/distribution/delegators_validator/rewards",
		withdrawDelegationRewardsHandlerFn(cliCtx),
	).Methods("POST")

	// Replace the rewards withdrawal address
	r.HandleFunc(
		"/distribution/delegators/withdraw_account",
		setDelegatorWithdrawalAddrHandlerFn(cliCtx),
	).Methods("POST")

	// Withdraw validator rewards and commission
	r.HandleFunc(
		"/distribution/validators/rewards",
		withdrawValidatorRewardsHandlerFn(cliCtx),
	).Methods("POST")

}

type (
	withdrawRewardsReq struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		DelegatorAcc string       `json:"delegator_acc" yaml:"delegator_acc"`
		ValidatorAcc string       `json:"validator_acc" yaml:"validator_acc"`
	}

	setWithdrawalAddrReq struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		DelegatorAcc string       `json:"delegator_acc" yaml:"delegator_acc"`
		WithdrawAcc  string       `json:"withdraw_acc" yaml:"withdraw_acc"`
	}

	fundCommunityPoolReq struct {
		BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
		Amount       string       `json:"amount" yaml:"amount"`
		DepositorAcc string       `json:"depositor_acc" yaml:"depositor_acc"`
	}
)

// Withdraw delegator rewards
func withdrawDelegatorRewardsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delegatorAcc, err := chainType.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(delegatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, err := common.WithdrawAllDelegatorRewards(cliCtx, auth, queryRoute, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, msg)
	}
}

// Withdraw delegation rewards
func withdrawDelegationRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delegatorAcc, err := chainType.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		validatorAcc, err := chainType.NewAccountIDFromStr(req.ValidatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(delegatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgWithdrawDelegatorReward(auth, delegatorAcc, validatorAcc)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, []sdk.Msg{msg})
	}
}

// Replace the rewards withdrawal address
func setDelegatorWithdrawalAddrHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setWithdrawalAddrReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) { //bugs, x
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		delegatorAcc, err := chainType.NewAccountIDFromStr(req.DelegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		withdrawAcc, err := chainType.NewAccountIDFromStr(req.WithdrawAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(delegatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, delegatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSetWithdrawAccountId(auth, delegatorAcc, withdrawAcc)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, []sdk.Msg{msg})
	}
}

// Withdraw validator rewards and commission
func withdrawValidatorRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		validatorAcc, err := chainType.NewAccountIDFromStr(req.ValidatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(validatorAcc)
		auth, err := txutil.QueryAccountAuth(ctx, validatorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgWithdrawDelegatorReward(auth, validatorAcc, validatorAcc)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, []sdk.Msg{msg})
	}
}

func fundCommunityPoolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req fundCommunityPoolReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()

		amount, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		depositor, err := chainType.NewAccountIDFromStr(req.DepositorAcc)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := txutil.NewKuCLICtx(cliCtx).WithFromAccount(depositor)
		auth, err := txutil.QueryAccountAuth(ctx, depositor)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgFundCommunityPool(auth, amount, depositor)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txutil.WriteGenerateStdTxResponse(w, txutil.NewKuCLICtx(cliCtx), req.BaseReq, []sdk.Msg{msg})
	}
}

func checkDelegatorAddressVar(w http.ResponseWriter, r *http.Request) (chainType.AccountID, bool) {
	accID, err := chainType.NewAccountIDFromStr(mux.Vars(r)["delegatorAddr"])
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return chainType.EmptyAccountID(), false
	}

	return accID, true
}

func checkValidatorAddressVar(w http.ResponseWriter, r *http.Request) (chainType.AccountID, bool) {
	// FIXME: support accountID
	addr, err := chainType.NewAccountIDFromStr(mux.Vars(r)["validatorAddr"])
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return chainType.EmptyAccountID(), false
	}

	return addr, true
}
