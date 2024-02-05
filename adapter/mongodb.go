package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/tools"
)

type MongoDbCredentialsService struct {
	credentialsRepository MongoDbCredentialsRepository
}

func NewMongoDbCredentialsService(credentialsRepository MongoDbCredentialsRepository) core.CredentialsPersistenceService {
	return &MongoDbCredentialsService{credentialsRepository: credentialsRepository}
}

func (m *MongoDbCredentialsService) ExistsCredentialsByIdAndPassword(credentials core.Credentials) (core.Credentials, bool) {
	credentialsDB, err := m.credentialsRepository.FindById(context.Background(), credentials.ID)
	if err != nil {
		return core.Credentials{}, false
	}
	var result core.Credentials
	tools.Mapper(credentialsDB, &result)

	return result, true
}

func (m *MongoDbCredentialsService) InsertCredentials(credentials core.Credentials) (core.Credentials, error) {
	ctx := context.Background()
	var credentialsDB CredentialsDB
	tools.Mapper(credentials, &credentialsDB)

	_, err := m.credentialsRepository.InsertOne(ctx, &credentialsDB)
	if err != nil {
		return core.Credentials{}, err
	}

	insertionDB, err := m.credentialsRepository.FindById(ctx, credentials.ID)
	var result core.Credentials
	tools.Mapper(insertionDB, &result)

	return result, nil
}

func (m *MongoDbCredentialsService) FindCredentialsById(credentials core.FullCredentials) (core.FullCredentials, error) {
	credentialsDB, err := m.credentialsRepository.FindById(context.Background(), credentials.ID)
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

	_, err := m.credentialsRepository.UpdateById(context.Background(), &credentialsDB, credentials.ID)
	return err
}

type MongoDbUserService struct {
	userRepository MongoDbUserRepository
}

func NewMongoDbUserService(userRepository MongoDbUserRepository) core.UserPersistenceService {
	return &MongoDbUserService{userRepository: userRepository}
}

func (m *MongoDbUserService) FindUserById(id string) (core.User, error) {
	userDB, err := m.userRepository.FindById(context.Background(), id)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) FindUserByEmail(email string) (core.User, error) {
	userDB, err := m.userRepository.FindByEmail(context.Background(), email)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) FindUserByPhoneNumber(phoneNumber string) (core.User, error) {
	userDB, err := m.userRepository.FindByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		return core.User{}, err
	}

	var result core.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) InsertUser(user core.User) (core.User, error) {
	//TODO
	panic("implement me")
}

func (m *MongoDbUserService) UpdateUser(user core.User) (core.User, error) {
	//TODO
	panic("implement me")
}
