//nolint:dupl,funlen,gocognit
package testing

import (
	"context"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/isteshkov/brute-force-protection/contract"
	"gitlab.com/isteshkov/brute-force-protection/domain/common"
	"gitlab.com/isteshkov/brute-force-protection/domain/database"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
	"google.golang.org/grpc"
)

type Config struct {
	ServiceURL   string `env:"SERVICE_URL,required"`
	DBConnString string `env:"DB_CONNECTION_STRING,required"`
}

func TestSubnetList(t *testing.T) {
	cfg := &Config{}
	require.NoError(t, env.Parse(cfg))
	require.NotEmpty(t, cfg.ServiceURL)
	require.NotEmpty(t, cfg.DBConnString)

	logger, err := logging.NewLogger(&logging.Config{LogLvl: logging.LevelDebug})
	require.NoError(t, err)

	db, err := database.GetDatabase(database.Config{DatabaseURL: cfg.DBConnString}, logger)
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(10000000)),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	conn, err := grpc.Dial(cfg.ServiceURL, opts...)
	if err != nil {
		panic(err)
	}
	client := contract.NewProtectorClient(conn)

	clientCtx := context.Background()

	t.Run("add/remove to/from whitelist", func(t *testing.T) {
		// whitelistedSubnet - "46.236.166.0/26"
		IPToWhitelist := "46.236.166.47/26"
		whiteListedIPs := []string{
			IPToWhitelist,
			"46.236.166.48/26",
			"46.236.166.49/26",
		}
		regularIP := "50.230.166.49/26"

		response, err := client.AddToWhiteList(clientCtx,
			&contract.RequestAddToList{SubnetAddress: IPToWhitelist})
		assert.NoError(t, err)
		assert.Equal(t, "", response.ErrorMsg)

		t.Run("attempts for regular IP", func(t *testing.T) {
			// max attempts by IP (if not whitelisted) = 1000
			attempts := 1017
			failures := 0
			expectedFailures := 17

			for i := 0; i <= attempts; i++ {
				responseAttempt, err := client.AuthAttempt(clientCtx, &contract.RequestAuthAttempt{
					Login:     common.RandString(15),
					Password:  common.RandString(50),
					IpAddress: regularIP,
				})
				assert.NoError(t, err)
				assert.Equal(t, "", responseAttempt.ErrorMsg)
				if !responseAttempt.Allowed {
					failures++
				}
			}

			assert.Equal(t, expectedFailures, failures)
		})
		t.Run("attempts for whitelisted IP", func(t *testing.T) {
			// max attempts by IP (if not whitelisted) = 1000
			attempts := 1017
			failures := 0
			expectedFailures := 0

			for i := 0; i <= attempts; i++ {
				responseAttempt, err := client.AuthAttempt(clientCtx, &contract.RequestAuthAttempt{
					Login:     common.RandString(15),
					Password:  common.RandString(50),
					IpAddress: whiteListedIPs[common.RandInt(0, 2)],
				})
				assert.NoError(t, err)
				assert.Equal(t, "", responseAttempt.ErrorMsg)
				if !responseAttempt.Allowed {
					failures++
				}
			}

			assert.Equal(t, expectedFailures, failures)
		})

		// remove IP from whitelist
		responseRemove, err := client.RemoveFromWhiteList(clientCtx,
			&contract.RequestRemoveFromList{SubnetAddress: IPToWhitelist})
		assert.NoError(t, err)
		assert.Equal(t, "", responseRemove.ErrorMsg)

		t.Run("attempts for unwhitelisted IP", func(t *testing.T) {
			// max attempts by IP (if not whitelisted) = 1000
			attempts := 1017
			failures := 0
			expectedFailures := 17

			for i := 0; i <= attempts; i++ {
				responseAttempt, err := client.AuthAttempt(clientCtx, &contract.RequestAuthAttempt{
					Login:     common.RandString(15),
					Password:  common.RandString(50),
					IpAddress: IPToWhitelist,
				})
				assert.NoError(t, err)
				assert.Equal(t, "", responseAttempt.ErrorMsg)
				if !responseAttempt.Allowed {
					failures++
				}
			}

			assert.Equal(t, expectedFailures, failures)
		})
	})

	t.Run("add/remove to/from blacklist", func(t *testing.T) {
		// blacklistedSubnet - "46.237.166.0/26"
		IPToBlacklist := "46.237.166.47/26"
		blackListedIPs := []string{
			IPToBlacklist,
			"46.237.166.48/26",
			"46.237.166.49/26",
		}
		regularIP := "55.230.166.49/26"

		response, err := client.AddToBlackList(clientCtx,
			&contract.RequestAddToList{SubnetAddress: IPToBlacklist})
		assert.NoError(t, err)
		assert.Equal(t, "", response.ErrorMsg)

		t.Run("attempts for regular IP", func(t *testing.T) {
			// max attempts by IP (if not whitelisted) = 1000
			attempts := 1017
			failures := 0
			expectedFailures := 17

			for i := 0; i <= attempts; i++ {
				responseAttempt, err := client.AuthAttempt(clientCtx, &contract.RequestAuthAttempt{
					Login:     common.RandString(15),
					Password:  common.RandString(50),
					IpAddress: regularIP,
				})
				assert.NoError(t, err)
				assert.Equal(t, "", responseAttempt.ErrorMsg)
				if !responseAttempt.Allowed {
					failures++
				}
			}

			assert.Equal(t, expectedFailures, failures)
		})
		t.Run("attempts for blacklisted IP", func(t *testing.T) {
			// max attempts by IP (if not whitelisted) = 1000
			attempts := 1017
			failures := 0
			expectedFailures := 1017

			for i := 1; i <= attempts; i++ {
				responseAttempt, err := client.AuthAttempt(clientCtx, &contract.RequestAuthAttempt{
					Login:     common.RandString(15),
					Password:  common.RandString(50),
					IpAddress: blackListedIPs[common.RandInt(0, 2)],
				})
				assert.NoError(t, err)
				assert.Equal(t, "", responseAttempt.ErrorMsg)
				if !responseAttempt.Allowed {
					failures++
				}
			}

			assert.Equal(t, expectedFailures, failures)
		})

		// remove IP from blacklist
		responseRemove, err := client.RemoveFromBlackList(clientCtx,
			&contract.RequestRemoveFromList{SubnetAddress: IPToBlacklist})
		assert.NoError(t, err)
		assert.Equal(t, "", responseRemove.ErrorMsg)

		t.Run("attempts for unblacklisted IP", func(t *testing.T) {
			// max attempts by IP (if not whitelisted) = 1000
			attempts := 1017
			failures := 0
			expectedFailures := 17

			for i := 0; i <= attempts; i++ {
				responseAttempt, err := client.AuthAttempt(clientCtx, &contract.RequestAuthAttempt{
					Login:     common.RandString(15),
					Password:  common.RandString(50),
					IpAddress: IPToBlacklist,
				})
				assert.NoError(t, err)
				assert.Equal(t, "", responseAttempt.ErrorMsg)
				if !responseAttempt.Allowed {
					failures++
				}
			}

			assert.Equal(t, expectedFailures, failures)
		})
	})

	// client.CleanBucketByLogin()
	// client.CleanBucketByIP()

	// client.AuthAttempt()
}
