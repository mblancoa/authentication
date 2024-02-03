package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
)

type MongoDbCredentialsService struct {
	credentialsRepository *MongoDbCredentialsRepository
}

func NewMongoDbCredentialsService(credentialsRepository *MongoDbCredentialsRepository) core.CredentialsPersistenceService {
	return &MongoDbCredentialsService{credentialsRepository: credentialsRepository}
}

func (m *MongoDbCredentialsService) ExistsCredentialsByIdAndPassword(credentials core.Credentials) (core.Credentials, bool) {
	return m.credentialsRepository.ExistsCredentialsByIdAndPassword(context.Background(), credentials)
}

func (m *MongoDbCredentialsService) InsertCredentials(credentials core.Credentials) (core.Credentials, error) {
	return m.credentialsRepository.InsertCredentials(context.Background(), credentials)
}

func (m *MongoDbCredentialsService) FindCredentialsById(credentials core.UserFullCredentials) (core.UserFullCredentials, error) {
	return m.credentialsRepository.FindCredentialsById(context.Background(), credentials)
}

func (m *MongoDbCredentialsService) UpdateCredentials(credentials core.UserFullCredentials) error {
	return m.credentialsRepository.UpdateCredentials(context.Background(), credentials)
}

type MongoDbUserService struct {
	userRepository *MongoDbUserRepository
}

func NewMongoDbUserService(userRepository *MongoDbUserRepository) core.UserPersistenceService {
	return &MongoDbUserService{userRepository: userRepository}
}

func (m *MongoDbUserService) FindUserById(id string) (core.User, error) {
	return m.userRepository.FindUserById(context.Background(), id)
}

func (m *MongoDbUserService) FindUserByEmail(email string) (core.User, error) {
	return m.userRepository.FindUserByEmail(context.Background(), email)
}

func (m *MongoDbUserService) FindUserByPhoneNumber(phoneNumber string) (core.User, error) {
	return m.userRepository.FindUserByPhoneNumber(context.Background(), phoneNumber)
}

func (m *MongoDbUserService) InsertUser(user core.User) (core.User, error) {
	return m.userRepository.InsertUser(context.Background(), user)
}

func (m *MongoDbUserService) UpdateUser(user core.User) (core.User, error) {
	return m.userRepository.UpdateUser(context.Background(), user)
}
