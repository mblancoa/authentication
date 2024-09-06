package mongodb

import (
	"context"
	"errors"
	"github.com/mblancoa/authentication/core/domain"
	"github.com/mblancoa/authentication/core/ports"
	errors2 "github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbCredentialsService struct {
	credentialsRepository MongoDbCredentialsRepository
}

func NewMongoDbCredentialsService(credentialsRepository MongoDbCredentialsRepository) ports.CredentialsPersistenceService {
	return &MongoDbCredentialsService{credentialsRepository: credentialsRepository}
}

func (m *MongoDbCredentialsService) InsertCredentials(credentials domain.Credentials) (domain.Credentials, error) {
	ctx := context.Background()
	credentialsDB := CredentialsDB{}
	tools.Mapper(&credentials, &credentialsDB)

	_, err := m.credentialsRepository.InsertOne(ctx, &credentialsDB)
	if err != nil {
		return domain.Credentials{}, errors2.NewGenericErrorByCause("Error inserting credentials", err)
	}

	result := domain.Credentials{}
	tools.Mapper(&credentialsDB, &result)
	return result, nil
}

func (m *MongoDbCredentialsService) FindCredentialsById(id string) (domain.FullCredentials, error) {
	credentialsDB, err := m.credentialsRepository.FindById(context.Background(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.FullCredentials{}, errors2.NewNotFoundError(err.Error())
		} else {
			return domain.FullCredentials{}, errors2.NewGenericError(err.Error())
		}
	}

	var result domain.FullCredentials
	tools.Mapper(credentialsDB, &result)

	return result, nil
}

func (m *MongoDbCredentialsService) UpdateCredentials(credentials domain.FullCredentials) error {
	db := CredentialsDB{}
	tools.Mapper(&credentials, &db)
	db.Id = ""

	_, err := m.credentialsRepository.UpdateById(context.Background(), &db, credentials.Id)
	if err != nil {
		return errors2.NewGenericError(err.Error())
	}
	return nil
}

type MongoDbUserService struct {
	userRepository MongoDbUserRepository
}

func NewMongoDbUserService(userRepository MongoDbUserRepository) ports.UserPersistenceService {
	return &MongoDbUserService{userRepository: userRepository}
}

func (m *MongoDbUserService) FindUserById(id string) (domain.User, error) {
	userDB, err := m.userRepository.FindById(context.Background(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, errors2.NewNotFoundError(err.Error())
		} else {
			return domain.User{}, errors2.NewGenericError(err.Error())
		}
	}

	var result domain.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) FindUserByEmail(email string) (domain.User, error) {
	userDB, err := m.userRepository.FindByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, errors2.NewNotFoundError(err.Error())
		} else {
			return domain.User{}, errors2.NewGenericError(err.Error())
		}
	}

	var result domain.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) FindUserByPhoneNumber(phoneNumber string) (domain.User, error) {
	userDB, err := m.userRepository.FindByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, errors2.NewNotFoundError(err.Error())
		} else {
			return domain.User{}, errors2.NewGenericError(err.Error())
		}
	}

	var result domain.User
	tools.Mapper(userDB, &result)

	return result, nil
}

func (m *MongoDbUserService) InsertUser(user domain.User) (domain.User, error) {
	db := UserDB{}
	tools.Mapper(&user, &db)
	_, err := m.userRepository.InsertOne(context.Background(), &db)
	if err != nil {
		return domain.User{}, errors2.NewGenericErrorByCause("Error inserting credentials", err)
	}

	result := domain.User{}
	tools.Mapper(&db, &result)
	return result, nil
}

func (m *MongoDbUserService) UpdateUser(user domain.User) error {
	db := UserDB{}
	tools.Mapper(&user, &db)
	db.Id = ""
	_, err := m.userRepository.UpdateById(context.Background(), &db, user.Id)

	if err != nil {
		return errors2.NewGenericError(err.Error())
	}
	return nil
}
