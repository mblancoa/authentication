package adapter

import "github.com/mblancoa/authentication/core"

type MongoDbUserCredentialsRepository struct {
}

func NewMongoDbUserCredentialsRepository() *MongoDbUserCredentialsRepository {
	return &MongoDbUserCredentialsRepository{}
}

func (m *MongoDbUserCredentialsRepository) ExistsUserCredentialsByIdAndPassword(credentials core.UserCredentials) (core.UserCredentials, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserCredentialsRepository) InsertUserCredentials(credentials core.UserCredentials) (core.UserCredentials, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserCredentialsRepository) FindUserCredentialsById(credentials core.UserFullCredentials) (core.UserFullCredentials, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserCredentialsRepository) UpdateUserCredentials(credentials core.UserFullCredentials) error {
	//TODO implement me
	panic("implement me")
}

type MongoDbUserRepository struct {
}

func NewMongoDbUserRepository() *MongoDbUserRepository {
	return &MongoDbUserRepository{}
}

func (m *MongoDbUserRepository) FindUserById(id string) (core.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDbUserRepository) FindUserByEmail(email string) (core.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserRepository) FindUserByPhoneNumber(phoneNumber string) (core.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserService) InsertUser(user core.User) (core.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserRepository) UpdateUser(user core.User) (core.User, error) {
	//TODO implement me
	panic("implement me")
}
