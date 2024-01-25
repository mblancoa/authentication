module github.com/mblancoa/authentication/repository-mongodb

go 1.21.3

require (
	github.com/mblancoa/authentication/domain latest
)

replace (
	github.com/mblancoa/authentication/domain latest  => ./../core
)