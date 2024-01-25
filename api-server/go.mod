module github.com/x/authentication/apiserver

go 1.21.3-

require (
	github.com/x/authentication/core latest
	github.com/x/authentication/repository-mongodb latest
	github.com/x/authentication/cache-redis latest
)

replace (
	github.com/x/authentication/core latest  => ./../core
	github.com/x/authentication/repository-mongodb latest => ../repository-mongodb
	github.com/x/authentication/cache-redis latest => ../cache-redis
)