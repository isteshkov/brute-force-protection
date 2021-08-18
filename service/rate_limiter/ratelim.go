package ratelimiter

import (
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
)

type RateLim struct {
	LoginBunch    Bunch
	PasswordBunch Bunch
	IPBunch       Bunch

	MaxLoginAttempts    int
	MaxPasswordAttempts int
	MaxIPAttempts       int

	logger logging.Logger
}

func NewRateLim(maxLogin, maxPass, maxIP int, l logging.Logger) *RateLim {
	return &RateLim{
		LoginBunch:          NewBunch(),
		PasswordBunch:       NewBunch(),
		IPBunch:             NewBunch(),
		MaxLoginAttempts:    maxLogin,
		MaxPasswordAttempts: maxPass,
		MaxIPAttempts:       maxIP,
		logger:              l,
	}
}

func (rl *RateLim) AttemptByLogin(login string) bool {
	bucket, err := rl.LoginBunch.Get(login)
	if err != nil {
		rl.LoginBunch.Set(login, NewBucket(WithAttempts(1)))
		return true
	}

	if bucket.IsMinuteSpent() {
		bucket.Reset()
		return true
	}

	if bucket.Attempts <= rl.MaxLoginAttempts {
		bucket.Incr()
		return true
	}

	return false
}

func (rl *RateLim) AttemptByPassword(password string) bool {
	bucket, err := rl.PasswordBunch.Get(password)
	if err != nil {
		rl.PasswordBunch.Set(password, NewBucket(WithAttempts(1)))
		return true
	}

	if bucket.IsMinuteSpent() {
		bucket.Reset()
		return true
	}

	if bucket.Attempts <= rl.MaxPasswordAttempts {
		bucket.Incr()
		return true
	}

	return false
}

func (rl *RateLim) AttemptByIP(ip string) bool {
	bucket, err := rl.IPBunch.Get(ip)
	if err != nil {
		rl.IPBunch.Set(ip, NewBucket(WithAttempts(1)))
		return true
	}

	if bucket.IsMinuteSpent() {
		bucket.Reset()
		return true
	}

	if bucket.Attempts <= rl.MaxIPAttempts {
		bucket.Incr()
		return true
	}

	return false
}

func (rl *RateLim) CleanBucketByLogin(login string) error {
	return rl.LoginBunch.Delete(login)
}

func (rl *RateLim) CleanBucketByIP(ip string) error {
	return rl.IPBunch.Delete(ip)
}
