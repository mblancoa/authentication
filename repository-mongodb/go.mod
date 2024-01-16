module github.com/x/authentication/repository-mongodb

go 1.21.3

require (
	github.com/x/authentication/domain latest
)

replace (
	github.com/x/authentication/domain latest  => ../domain
)