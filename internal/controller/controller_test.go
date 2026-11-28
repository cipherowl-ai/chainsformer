package controller

import (
	"testing"

	"go.uber.org/fx"

	"github.com/coinbase/chainsformer/internal/chainstorage"
	"github.com/coinbase/chainsformer/internal/config"
	"github.com/coinbase/chainsformer/internal/utils/testapp"
	"github.com/coinbase/chainsformer/internal/utils/testutil"
)

func TestNewController(t *testing.T) {
	testapp.TestAllConfigs(t, func(t *testing.T, cfg *config.Config) {
		require := testutil.Require(t)

		var controller Controller
		testapp.New(
			t,
			testapp.WithConfig(cfg),
			Module,
			chainstorage.Module,
			fx.Populate(&controller),
		)

		require.NotNil(controller)
	})
}

func TestNewController_ServesBlockTransactions(t *testing.T) {
	// TestConfigs covers only bitcoin, so build each chain's config here.
	for _, configName := range []string{
		"bitcoin-mainnet", "bitcoincash-mainnet", "dash-mainnet", "dogecoin-mainnet", "litecoin-mainnet", "zcash-mainnet",
		"tron-mainnet",
	} {
		t.Run(configName, func(t *testing.T) {
			blockchain, network, err := config.ParseConfigName(configName)
			testutil.Require(t).NoError(err)
			testapp.TestAllEnvs(t, func(t *testing.T, base *config.Config) {
				require := testutil.Require(t)
				cfg, err := config.New(
					config.WithEnvironment(base.Env()),
					config.WithBlockchain(blockchain),
					config.WithNetwork(network),
				)
				require.NoError(err)
				require.Equal(blockchain, cfg.Blockchain())

				var controller Controller
				testapp.New(
					t,
					testapp.WithConfig(cfg),
					Module,
					chainstorage.Module,
					fx.Populate(&controller),
				)

				var names []string
				for _, table := range controller.Tables() {
					names = append(names, table.GetTableName())
				}
				require.Contains(names, "table=block_transactions/format=native/encoding=none")
			})
		})
	}
}
