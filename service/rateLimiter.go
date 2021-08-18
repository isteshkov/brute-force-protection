package service

type RateLimiter interface {
	AttemptByLogin(login string) bool
	AttemptByPassword(password string) bool
	AttemptByIP(ip string) bool

	CleanBucketByLogin(login string) error
	CleanBucketByIP(ip string) error
}
