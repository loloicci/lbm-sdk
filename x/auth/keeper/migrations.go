package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper AccountKeeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper AccountKeeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return nil
}
