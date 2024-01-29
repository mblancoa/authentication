module github.com/mblancoa/authentication/apiserver

go 1.21.3-

require (
	github.com/mblancoa/authentication/core latest
	github.com/mblancoa/authentication/repository-mongodb latest
	github.com/mblancoa/authentication/cache-redis latest
)

replace (
	github.com/mblancoa/authentication/core latest  => ./../core
	github.com/mblancoa/authentication/repository-mongodb latest => ../repository-mongodb
	github.com/mblancoa/authentication/cache-redis latest => ../cache-redis
)