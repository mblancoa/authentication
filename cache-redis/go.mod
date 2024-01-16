module github.com/x/authentication/cache-redis

go 1.21.3

require (
	github.com/x/authentication/domain latest
)

replace (
	github.com/x/authentication/domain latest  => ../domain
)