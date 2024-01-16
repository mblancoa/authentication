module github.com/x/authentication/api-client

go 1.21.3

require (
	github.com/x/authentication/api-domain latest
)

replace (
	github.com/x/authentication/api-domain latest => ./../api-domain
)