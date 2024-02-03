package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/tools"
)

type MongoDbCredentialsService struct {
	credentialsRepository *MongoDbCredentialsRepository
}

func NewMongoDbCredentialsService(credentialsRepository *MongoDbCredentialsRepository) core.CredentialsPersistenceService {
	return &MongoDbCredentialsService{credentialsRepository: credentialsRepository}
}

func (m *MongoDbCredentialsService) ExistsCredentialsByIdAndPassword(credentials core.Credentials) (core.Credentials, bool) {
	var credentialsDB CredentialsDB
	tools.Mapper(credentials, &credentialsDB)

	credentialsDB, ok := m.credentialsRepository.ExistsCredentialsByIdAndPassword(context.Background(), credentialsDB)
	var result core.Credentials
	tools.Mapper(credentialsDB, &result)

	return result, ok
}

func (m *MongoDbCredentialsService) InsertCredentials(credentials core.Credentials) (core.Credentials, error) {
	var credentialsDB CredentialsDB
	tools.Mapper(credentials, &credentialsDB)

	credentialsDB, err := m.credentialsRepository.InsertCredentials(context.Background(), credentialsDB)
	if err != nil {
		return core.Credentials{}, err
	}

	var result core.Credentials
	tools.Mapper(credentialsDB, &result)

	return result, nil
}

func (m *MongoDbCredentialsService) FindCredentialsById(credentials core.FullCredentials) (core.FullCredentials, error) {
	var credentialsDB CredentialsDB
	tools.Mapper(credentials, &credentialsDB)

	credentialsDB, err := m.credentialsRepository.FindCredentialsById(context.Background(), credentialsDB)
	if err != nil {
		return core.FullCredentials{}, err
	}

	var result core.FullCredentials
	tools.Mapper(credentialsDB, &result)

	return result, nil
}

func (m *MongoDbCredentialsService) UpdateCredentials(credentials core.FullCredentials) error {
	var credentialsDB CredentialsDB
	tools.Mapper(credentials, &credentialsDB)

	return m.credentialsRepository.UpdateCredentials(context.Background(), credentialsDB)
}

type MongoDbUserService struct {
	userRepository *MongoDbUserRepository
}

func NewMongoDbUserService(userRepository *MongoDbUserRepository) core.UserPersistenceService {
	return &MongoDbUserService{userRepository: userRepository}
}

func (m *MongoDbUserService) FindUserById(id string) (core.User, error) {
	userDB, err := m.userRepository.FindUserById(context.Background(), id)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) FindUserByEmail(email string) (core.User, error) {
	userDB, err := m.userRepository.FindUserByEmail(context.Background(), email)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) FindUserByPhoneNumber(phoneNumber string) (core.User, error) {
	userDB, err := m.userRepository.FindUserByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) InsertUser(user core.User) (core.User, error) {
	var userDB UserDB
	tools.Mapper(user, &userDB)

	userDB, err := m.userRepository.InsertUser(context.Background(), user)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) UpdateUser(user core.User) (core.User, error) {
	var userDB UserDB
	tools.Mapper(user, &userDB)

	userDB, err := m.userRepository.UpdateUser(context.Background(), user)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}
