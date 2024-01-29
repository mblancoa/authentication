module github.com/mblancoa/authentication/cache-redis

go 1.21.3

require (
	github.com/mblancoa/authentication/core latest
)

replace (
	github.com/mblancoa/authentication/core latest  => ./../core
)