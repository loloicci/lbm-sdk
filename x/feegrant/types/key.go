package types

import (
	sdk "github.com/line/lbm-sdk/types"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "feegrant"

	// StoreKey is the store key string for supply
	StoreKey = ModuleName

	// RouterKey is the message route for supply
	RouterKey = ModuleName

	// QuerierRoute is the querier route for supply
	QuerierRoute = ModuleName
)

var (
	// FeeAllowanceKeyPrefix is the set of the kvstore for fee allowance data
	FeeAllowanceKeyPrefix = []byte{0x00}
)

// FeeAllowanceKey is the canonical key to store a grant from granter to grantee
// We store by grantee first to allow searching by everyone who granted to you
func FeeAllowanceKey(granter sdk.AccAddress, grantee sdk.AccAddress) []byte {
	return append(FeeAllowancePrefixByGrantee(grantee), granter.Bytes()...)
}

// FeeAllowancePrefixByGrantee returns a prefix to scan for all grants to this given address.
func FeeAllowancePrefixByGrantee(grantee sdk.AccAddress) []byte {
	return append(FeeAllowanceKeyPrefix, grantee.Bytes()...)
}
