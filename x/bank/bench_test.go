package bank_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

var moduleAccAddr = supply.NewModuleAddress(staking.BondedPoolName)

func BenchmarkOneBankSendTxPerBlock(b *testing.B) {
	// Add an account at genesis
	acc := auth.BaseAccount{
		Address: addr1,
	}

	// construct genesis state
	genAccs := []authexported.GenesisAccount{&acc}
	benchmarkApp := simapp.SetupWithGenesisAccounts(genAccs)
	ctx := benchmarkApp.BaseApp.NewContext(false, abci.Header{})

	// some value conceivably higher than the benchmarks would ever go
	err := benchmarkApp.BankKeeper.SetBalances(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 100000000000)))
	require.NoError(b, err)

	benchmarkApp.Commit()

	// Precompute all txs
	txs := simapp.GenSequenceOfTxs([]sdk.Msg{sendMsg1}, []uint64{0}, []uint64{uint64(0)}, b.N, priv1)
	b.ResetTimer()

	height := int64(3)

	// Run this with a profiler, so its easy to distinguish what time comes from
	// Committing, and what time comes from Check/Deliver Tx.
	for i := 0; i < b.N; i++ {
		benchmarkApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: height}})
		_, _, err := benchmarkApp.Check(txs[i])
		if err != nil {
			panic("something is broken in checking transaction")
		}

		benchmarkApp.Deliver(txs[i])
		benchmarkApp.EndBlock(abci.RequestEndBlock{Height: height})
		benchmarkApp.Commit()
		height++
	}
}

func BenchmarkOneBankMultiSendTxPerBlock(b *testing.B) {
	// Add an account at genesis
	acc := auth.BaseAccount{
		Address: addr1,
	}

	// Construct genesis state
	genAccs := []authexported.GenesisAccount{&acc}
	benchmarkApp := simapp.SetupWithGenesisAccounts(genAccs)
	ctx := benchmarkApp.BaseApp.NewContext(false, abci.Header{})

	// some value conceivably higher than the benchmarks would ever go
	err := benchmarkApp.BankKeeper.SetBalances(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 100000000000)))
	require.NoError(b, err)

	benchmarkApp.Commit()

	// Precompute all txs
	txs := simapp.GenSequenceOfTxs([]sdk.Msg{multiSendMsg1}, []uint64{0}, []uint64{uint64(0)}, b.N, priv1)
	b.ResetTimer()

	height := int64(3)

	// Run this with a profiler, so its easy to distinguish what time comes from
	// Committing, and what time comes from Check/Deliver Tx.
	for i := 0; i < b.N; i++ {
		benchmarkApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: height}})
		_, _, err := benchmarkApp.Check(txs[i])
		if err != nil {
			panic("something is broken in checking transaction")
		}

		benchmarkApp.Deliver(txs[i])
		benchmarkApp.EndBlock(abci.RequestEndBlock{Height: height})
		benchmarkApp.Commit()
		height++
	}
}
