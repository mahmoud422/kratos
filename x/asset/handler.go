package asset

import (
	"strconv"

	"github.com/KuChain-io/kuchain/chain/msg"
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/asset/keeper"
	"github.com/KuChain-io/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.AssetCoinsKeeper) msg.Handler {
	return func(ctx chainTypes.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.KuMsg:
			return handleMsgTransfer(ctx)
		case *types.MsgCreateCoin:
			return handleMsgCreate(ctx, k, msg)
		case *types.MsgIssueCoin:
			return handleMsgIssue(ctx, k, msg)
		case *types.MsgBurnCoin:
			return handleMsgBurn(ctx, k, msg)
		case *types.MsgLockCoin:
			return handleMsgLockCoin(ctx, k, msg)
		case *types.MsgUnlockCoin:
			return handleMsgUnlockCoin(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized asset message type: %T", msg)
		}
	}
}

// handleMsgTransfer Handle KuMsg for transfer.
func handleMsgTransfer(ctx chainTypes.Context) (*sdk.Result, error) {
	// no need process transfer
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgCreate Handle Msg create coin type
func handleMsgCreate(ctx chainTypes.Context, k keeper.AssetCoinsKeeper, msg *types.MsgCreateCoin) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgCreateCoinData{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg create coin data unmarshal error")
	}

	logger.Debug("handle coin create",
		"creator", msgData.Creator,
		"symbol", msgData.Symbol,
		"max_supply", msgData.MaxSupply,
		"isCanIssue", msgData.CanIssue,
		"isCanLock", msgData.CanLock,
		"issueHeight", msgData.IssueToHeight,
		"initSupply", msgData.InitSupply,
		"desc", string(msgData.Desc))

	ctx.RequireAccount(msgData.Creator)

	if err := k.Create(ctx.Context(),
		msgData.Creator, msgData.Symbol, msgData.MaxSupply,
		msgData.CanIssue, msgData.CanLock, msgData.IssueToHeight, msgData.InitSupply, msgData.Desc); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg create coin %s", msgData.Symbol)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, msgData.Creator.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msgData.Symbol.String()),
			sdk.NewAttribute(types.AttributeKeyMaxSupply, msgData.MaxSupply.String()),
			sdk.NewAttribute(types.AttributeKeyCanIssue, strconv.FormatBool(msgData.CanIssue)),
			sdk.NewAttribute(types.AttributeKeyCanLock, strconv.FormatBool(msgData.CanLock)),
			sdk.NewAttribute(types.AttributeKeyIssueToHeight, strconv.FormatInt(msgData.IssueToHeight, 10)),
			sdk.NewAttribute(types.AttributeKeyInit, msgData.InitSupply.String()),
			sdk.NewAttribute(types.AttributeKeyDescription, string(msgData.Desc)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgIssue Handle Msg Issue coin
func handleMsgIssue(ctx chainTypes.Context, k keeper.AssetCoinsKeeper, msg *types.MsgIssueCoin) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgIssueCoinData{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg issue coin data unmarshal error")
	}

	logger.Debug("handle coin issue",
		"creator", msgData.Creator,
		"symbol", msgData.Symbol,
		"amount", msgData.Amount)

	ctx.RequireAccount(msgData.Creator)

	stat, err := k.GetCoinStat(ctx.Context(), msgData.Creator, msgData.Symbol)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "get coin stat from coin %s", msg.Amount.String())
	}

	// if coins cannot be issue, if there is 1000 blocks after coin created, no one can issue
	if !stat.CanIssue && (ctx.BlockHeight() > (stat.CreateHeight + 5)) { // FIXME: for test
		return nil, sdkerrors.Wrapf(types.ErrAssetCoinCannotBeLock, "coin %s cannot be issue after 1000 block from coin create", msg.Amount.String())
	}

	if err := k.Issue(ctx.Context(), msgData.Creator, msgData.Symbol, msgData.Amount); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg issue coin %s", msgData.Symbol)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIssue,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyCreator, msgData.Creator.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msgData.Symbol.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msgData.Amount.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}

// handleMsgBurn Handle Msg Burn coin
func handleMsgBurn(ctx chainTypes.Context, k keeper.AssetCoinsKeeper, msg *types.MsgBurnCoin) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgBurnCoinData{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg burn coin data unmarshal error")
	}

	logger.Debug("handle coin burn",
		"id", msgData.Id,
		"amount", msgData.Amount)

	ctx.RequireAuth(msgData.Id)

	if err := k.Burn(ctx.Context(), msgData.Id, msgData.Amount); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg burn coin %s", msgData.Id)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIssue,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyFrom, msgData.Id.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msgData.Amount.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}

// handleMsgLockCoin Handle Msg lock coin
func handleMsgLockCoin(ctx chainTypes.Context, k keeper.AssetCoinsKeeper, msg *types.MsgLockCoin) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgLockCoinData{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg lock coin data unmarshal error")
	}

	logger.Debug("handle coin lock",
		"id", msgData.Id,
		"amount", msgData.Amount,
		"height", msgData.UnlockBlockHeight)

	ctx.RequireAuth(msgData.Id)

	for _, c := range msgData.Amount {
		creator, symbol, err := chainTypes.CoinAccountsFromDenom(c.Denom)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "get creator and symbol from coin %s", msg.Amount.String())
		}

		stat, err := k.GetCoinStat(ctx.Context(), creator, symbol)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "get coin stat from coin %s", msg.Amount.String())
		}

		if !stat.CanLock {
			return nil, sdkerrors.Wrapf(types.ErrAssetCoinCannotBeLock, "coin %s cannot be locked", msg.Amount.String())
		}
	}

	if err := k.LockCoins(ctx.Context(), msgData.Id, msgData.UnlockBlockHeight, msgData.Amount); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg lock coin %s", msgData.Id)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyFrom, msgData.Id.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msgData.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyUnlockHeight, strconv.Itoa(int(msgData.UnlockBlockHeight))),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgUnlockCoin Handle Msg lock coin
func handleMsgUnlockCoin(ctx chainTypes.Context, k keeper.AssetCoinsKeeper, msg *types.MsgUnlockCoin) (*sdk.Result, error) {
	logger := ctx.Logger()

	msgData := types.MsgUnlockCoinData{}
	if err := msg.UnmarshalData(Cdc(), &msgData); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg unlock coin data unmarshal error")
	}

	logger.Debug("handle coin lock",
		"id", msgData.Id,
		"amount", msgData.Amount)

	ctx.RequireAuth(msgData.Id)

	if err := k.UnLockCoins(ctx.Context(), msgData.Id, msgData.Amount); err != nil {
		return nil, sdkerrors.Wrapf(err, "msg unlock coin %s", msgData.Id)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyFrom, msgData.Id.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msgData.Amount.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
