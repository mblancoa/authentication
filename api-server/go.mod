module github.com/x/authentication/apiserver

go 1.21.3-

require (
	github.com/x/authentication/domain latest
	github.com/x/authentication/repository-mongodb latest
	github.com/x/authentication/cache-redis latest
	github.com/x/authentication/api-domain latest
)

replace (
	github.com/x/authentication/domain latest  => ../domain
	github.com/x/authentication/repository-mongodb latest => ../repository-mongodb
	github.com/x/authentication/cache-redis latest => ../cache-redis
	github.com/x/authentication/api-domain latest => ./../api-domain
)