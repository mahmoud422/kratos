package keeper

import (
	chainTypes "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Implements StakingHooks interface
var _ types.StakingHooks = Keeper{}

// AfterValidatorCreated - call hook if registered
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.AfterValidatorCreated(ctx, valAddr)
	}
}

// BeforeValidatorModified - call hook if registered
func (k Keeper) BeforeValidatorModified(ctx sdk.Context, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorModified(ctx, valAddr)
	}
}

// AfterValidatorRemoved - call hook if registered
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
	}
}

// AfterValidatorBonded - call hook if registered
func (k Keeper) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.AfterValidatorBonded(ctx, consAddr, valAddr)
	}
}

// AfterValidatorBeginUnbonding - call hook if registered
func (k Keeper) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	}
}

// BeforeDelegationCreated - call hook if registered
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
	}
}

// BeforeDelegationSharesModified - call hook if registered
func (k Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	}
}

// BeforeDelegationRemoved - call hook if registered
func (k Keeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	}
}

// AfterDelegationModified - call hook if registered
func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr chainTypes.AccountID, valAddr chainTypes.AccountID) {
	if k.hooks != nil {
		k.hooks.AfterDelegationModified(ctx, delAddr, valAddr)
	}
}

// BeforeValidatorSlashed - call hook if registered
func (k Keeper) BeforeValidatorSlashed(ctx sdk.Context, valAddr chainTypes.AccountID, fraction sdk.Dec) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
}
