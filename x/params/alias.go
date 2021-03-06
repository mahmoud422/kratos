package params

// nolint

import (
	"github.com/KuChain-io/kuchain/x/params/keeper"
	"github.com/KuChain-io/kuchain/x/params/types"
)

const (
	StoreKey  = types.StoreKey
	TStoreKey = types.TStoreKey
)

var (
	// functions aliases
	NewKeeper       = keeper.NewKeeper
	NewKeyTable     = types.NewKeyTable
	NewParamSetPair = types.NewParamSetPair
)

type (
	Keeper           = keeper.Keeper
	ParamSetPair     = types.ParamSetPair
	ParamSetPairs    = types.ParamSetPairs
	ParamSet         = types.ParamSet
	Subspace         = types.Subspace
	ReadOnlySubspace = types.ReadOnlySubspace
	KeyTable         = types.KeyTable
)
