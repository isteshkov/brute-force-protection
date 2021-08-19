package ratelimiter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/isteshkov/brute-force-protection/domain/common"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
)

func TestRateLim_AttemptByIP(t *testing.T) {
	logger, err := logging.NewLogger(&logging.Config{})
	assert.NoError(t, err)

	maxIPAttempts := 5
	rateLim := NewRateLim(1, 1, maxIPAttempts, logger)

	ip := RandomIP()

	allow := rateLim.AttemptByIP(ip)
	assert.True(t, allow)

	bucket, err := rateLim.IPBunch.Get(ip)
	assert.NoError(t, err)
	assert.Equal(t, 1, bucket.Attempts)

	for i := 1; i <= maxIPAttempts+1; i++ {
		allow = rateLim.AttemptByIP(ip)
		bucket, err = rateLim.IPBunch.Get(ip)
		assert.NoError(t, err)
		if i <= maxIPAttempts {
			assert.True(t, allow)
			assert.Equal(t, i+1, bucket.Attempts)
			continue
		}
		assert.False(t, allow)
		assert.Equal(t, i, bucket.Attempts)
	}

	bucket.Reset()

	allow = rateLim.AttemptByIP(ip)
	assert.True(t, allow)
	bucket, err = rateLim.IPBunch.Get(ip)
	assert.NoError(t, err)
	assert.Equal(t, 1, bucket.Attempts)
}

func TestRateLim_AttemptByLogin(t *testing.T) {
	logger, err := logging.NewLogger(&logging.Config{})
	assert.NoError(t, err)

	maxLoginAttempts := 5
	rateLim := NewRateLim(maxLoginAttempts, 1, 1, logger)

	login := common.RandString(15)

	allow := rateLim.AttemptByLogin(login)
	assert.True(t, allow)

	bucket, err := rateLim.LoginBunch.Get(login)
	assert.NoError(t, err)
	assert.Equal(t, 1, bucket.Attempts)

	for i := 1; i <= maxLoginAttempts+1; i++ {
		allow = rateLim.AttemptByLogin(login)
		bucket, err = rateLim.LoginBunch.Get(login)
		assert.NoError(t, err)
		if i <= maxLoginAttempts {
			assert.True(t, allow)
			assert.Equal(t, i+1, bucket.Attempts)
			continue
		}
		assert.False(t, allow)
		assert.Equal(t, i, bucket.Attempts)
	}

	bucket.Reset()

	allow = rateLim.AttemptByLogin(login)
	assert.True(t, allow)
	bucket, err = rateLim.LoginBunch.Get(login)
	assert.NoError(t, err)
	assert.Equal(t, 1, bucket.Attempts)
}

func TestRateLim_AttemptByPassword(t *testing.T) {
	logger, err := logging.NewLogger(&logging.Config{})
	assert.NoError(t, err)

	maxPasswordAttempts := 5
	rateLim := NewRateLim(1, maxPasswordAttempts, 1, logger)

	password := common.RandString(50)

	allow := rateLim.AttemptByPassword(password)
	assert.True(t, allow)

	bucket, err := rateLim.PasswordBunch.Get(password)
	assert.NoError(t, err)
	assert.Equal(t, 1, bucket.Attempts)

	for i := 1; i <= maxPasswordAttempts+1; i++ {
		allow = rateLim.AttemptByPassword(password)
		bucket, err = rateLim.PasswordBunch.Get(password)
		assert.NoError(t, err)
		if i <= maxPasswordAttempts {
			assert.True(t, allow)
			assert.Equal(t, i+1, bucket.Attempts)
			continue
		}
		assert.False(t, allow)
		assert.Equal(t, i, bucket.Attempts)
	}

	bucket.Reset()

	allow = rateLim.AttemptByPassword(password)
	assert.True(t, allow)
	bucket, err = rateLim.PasswordBunch.Get(password)
	assert.NoError(t, err)
	assert.Equal(t, 1, bucket.Attempts)
}

func TestRateLim_CleanBucketByIP(t *testing.T) {
	logger, err := logging.NewLogger(&logging.Config{})
	assert.NoError(t, err)

	rateLim := NewRateLim(1, 1, 1, logger)

	key := "192.168.0.1/24"
	srcBucket := NewBucket(WithAttempts(1))
	srcBucket.Incr()
	rateLim.IPBunch.Set(key, srcBucket)

	bucket, err := rateLim.IPBunch.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, srcBucket, bucket)

	err = rateLim.IPBunch.Delete(key)
	assert.NoError(t, err)

	bucketNil, err := rateLim.IPBunch.Get(key)
	assert.Error(t, err)
	assert.Nil(t, bucketNil)
	assert.Contains(t, err.Error(), "no such key")
}

func TestRateLim_CleanBucketByLogin(t *testing.T) {
	logger, err := logging.NewLogger(&logging.Config{})
	assert.NoError(t, err)

	rateLim := NewRateLim(1, 1, 1, logger)

	key := "Login_test"
	srcBucket := NewBucket(WithAttempts(1))
	srcBucket.Incr()
	rateLim.LoginBunch.Set(key, srcBucket)

	bucket, err := rateLim.LoginBunch.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, srcBucket, bucket)

	err = rateLim.LoginBunch.Delete(key)
	assert.NoError(t, err)

	bucketNil, err := rateLim.LoginBunch.Get(key)
	assert.Error(t, err)
	assert.Nil(t, bucketNil)
	assert.Contains(t, err.Error(), "no such key")
}

func RandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d/%d",
		common.RandInt(192, 198),
		common.RandInt(128, 191),
		common.RandInt(10, 50),
		common.RandInt(0, 254),
		common.RandInt(18, 24))
}
