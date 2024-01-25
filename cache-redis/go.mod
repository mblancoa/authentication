module github.com/x/authentication/cache-redis

go 1.21.3

require (
	github.com/x/authentication/core latest
)

replace (
	github.com/x/authentication/core latest  => ./../core
)