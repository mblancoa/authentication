package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
)

type MongoDbCredentialsService struct {
	credentialsRepository MongoDbCredentialsRepository
}

func NewMongoDbCredentialsService(credentialsRepository MongoDbCredentialsRepository) core.CredentialsPersistenceService {
	return &MongoDbCredentialsService{credentialsRepository: credentialsRepository}
}

func (m *MongoDbCredentialsService) CheckCredentials(credentials core.Credentials, maxAttempts int) (core.Credentials, error) {
	credentialsDB, err := m.credentialsRepository.FindById(context.Background(), credentials.Id)
	if err != nil {
		return core.Credentials{}, errors.NewNotFoundError(err.Error())
	}
	if credentialsDB.Password != credentials.Password {
		credentialsDB.Attempts++
		if credentialsDB.Attempts == 3 {
			credentialsDB.State = core.Blocked
		}
		_, err = m.credentialsRepository.UpdateById(context.Background(), credentialsDB, credentialsDB.Id)
		if err != nil {
			return core.Credentials{}, errors.NewGenericErrorByCause("Error updating credentials attempts", err)
		}
	} else if credentialsDB.State == core.Active && credentialsDB.Attempts != 0 {
		credentialsDB.Attempts = 0
		_, err = m.credentialsRepository.UpdateById(context.Background(), credentialsDB, credentialsDB.Id)
		if err != nil {
			return core.Credentials{}, errors.NewGenericErrorByCause("Error updating credentials state", err)
		}
	}
	var result core.Credentials
	tools.Mapper(credentialsDB, &result)

	return result, nil
}

func (m *MongoDbCredentialsService) InsertCredentials(credentials core.Credentials) (core.Credentials, error) {
	ctx := context.Background()
	credentialsDB := CredentialsDB{}
	tools.Mapper(&credentials, &credentialsDB)

	_, err := m.credentialsRepository.InsertOne(ctx, &credentialsDB)
	if err != nil {
		return core.Credentials{}, errors.NewGenericErrorByCause("Error inserting credentials", err)
	}

	return credentials, nil
}

func (m *MongoDbCredentialsService) FindCredentialsByUserId(id string) (core.FullCredentials, error) {
	credentialsDB, err := m.credentialsRepository.FindById(context.Background(), id)
	if err != nil {
		return core.FullCredentials{}, errors.NewNotFoundError(err.Error())
	}

	var result core.FullCredentials
	tools.Mapper(credentialsDB, &result)

	return result, nil
}

func (m *MongoDbCredentialsService) UpdateCredentials(credentials core.FullCredentials) error {
	db := CredentialsDB{}
	tools.Mapper(&credentials, &db)
	_, err := m.credentialsRepository.UpdateById(context.Background(), &db, credentials.Id)
	return err
}

type MongoDbUserService struct {
	userRepository MongoDbUserRepository
}

func NewMongoDbUserService(userRepository MongoDbUserRepository) core.UserPersistenceService {
	return &MongoDbUserService{userRepository: userRepository}
}

func (m *MongoDbUserService) FindUserByUserId(id string) (core.User, error) {
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
