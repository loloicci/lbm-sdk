package types_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	sdk "github.com/line/lbm-sdk/types"
	bank "github.com/line/lbm-sdk/x/bank/types"
)

func TestBalanceValidate(t *testing.T) {
	testCases := []struct {
		name    string
		balance bank.Balance
		expErr  bool
	}{
		{
			"valid balance",
			bank.Balance{
				Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
				Coins:   sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
			},
			false,
		},
		{"empty balance", bank.Balance{}, true},
		{
			"nil balance coins",
			bank.Balance{
				Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
			},
			false,
		},
		{
			"dup coins",
			bank.Balance{
				Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
				Coins: sdk.Coins{
					sdk.NewInt64Coin("uatom", 1),
					sdk.NewInt64Coin("uatom", 1),
				},
			},
			true,
		},
		{
			"invalid coin denom",
			bank.Balance{
				Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
				Coins: sdk.Coins{
					sdk.Coin{Denom: "", Amount: sdk.OneInt()},
				},
			},
			true,
		},
		{
			"negative coin",
			bank.Balance{
				Address: "link1mejkku76a2ec35262rdqddggzwrgtrh52t3t0c",
				Coins: sdk.Coins{
					sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(-1)},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			err := tc.balance.Validate()

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBalance_GetAddress(t *testing.T) {
	tests := []struct {
		name      string
		Address   string
		wantPanic bool
	}{
		{"empty address", "", false},
		{"malformed address", "invalid", false},
		{"valid address", "link1qz9c0r2jvpkccx67d5svg8kms6eu3k832hdc6p", false},
	}
	// Balance.GetAddress() does not validate the address.
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b := bank.Balance{Address: tt.Address}
			if tt.wantPanic {
				require.Panics(t, func() { b.GetAddress() })
			} else {
				require.True(t, b.GetAddress().Equals(sdk.AccAddress(tt.Address)))
			}
		})
	}
}

func TestSanitizeBalances(t *testing.T) {
	// 1. Generate balances
	tokens := sdk.TokensFromConsensusPower(81)
	coin := sdk.NewCoin("benchcoin", tokens)
	coins := sdk.Coins{coin}
	addrs, _ := makeRandomAddressesAndPublicKeys(20)

	var balances []bank.Balance
	for _, addr := range addrs {
		balances = append(balances, bank.Balance{
			Address: addr.String(),
			Coins:   coins,
		})
	}
	// 2. Sort the values.
	sorted := bank.SanitizeGenesisBalances(balances)

	// 3. Compare and ensure that all the values are sorted in ascending order.
	// Invariant after sorting:
	//    a[i] <= a[i+1...n]
	for i := 0; i < len(sorted); i++ {
		ai := sorted[i]
		// Ensure that every single value that comes after i is less than it.
		for j := i + 1; j < len(sorted); j++ {
			aj := sorted[j]

			if got := strings.Compare(ai.GetAddress().String(), aj.GetAddress().String()); got > 0 {
				t.Errorf("Balance(%d) > Balance(%d)", i, j)
			}
		}
	}
}

func makeRandomAddressesAndPublicKeys(n int) (accL []sdk.AccAddress, pkL []*ed25519.PubKey) {
	for i := 0; i < n; i++ {
		pk := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
		pkL = append(pkL, pk)
		accL = append(accL, sdk.BytesToAccAddress(pk.Address()))
	}
	return accL, pkL
}

var sink, revert interface{}

func BenchmarkSanitizeBalances500(b *testing.B) {
	benchmarkSanitizeBalances(b, 500)
}

func BenchmarkSanitizeBalances1000(b *testing.B) {
	benchmarkSanitizeBalances(b, 1000)
}

func benchmarkSanitizeBalances(b *testing.B, nAddresses int) {
	b.ReportAllocs()
	tokens := sdk.TokensFromConsensusPower(81)
	coin := sdk.NewCoin("benchcoin", tokens)
	coins := sdk.Coins{coin}
	addrs, _ := makeRandomAddressesAndPublicKeys(nAddresses)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var balances []bank.Balance
		for _, addr := range addrs {
			balances = append(balances, bank.Balance{
				Address: addr.String(),
				Coins:   coins,
			})
		}
		sink = bank.SanitizeGenesisBalances(balances)
	}
	if sink == nil {
		b.Fatal("Benchmark did not run")
	}
	sink = revert
}
