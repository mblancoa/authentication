package adapter

import "github.com/mblancoa/authentication/core"

type MongoDbUserCredentialsService struct {
	userCredentialsRepository *MongoDbUserCredentialsRepository
}

func NewMongoDbUserCredentialsService(userCredentialsRepository *MongoDbUserCredentialsRepository) core.UserCredentialsPersistenceService {
	return &MongoDbUserCredentialsService{userCredentialsRepository: userCredentialsRepository}
}

func (m *MongoDbUserCredentialsService) ExistsUserCredentialsByIdAndPassword(credentials core.UserCredentials) (core.UserCredentials, bool) {
	return m.ExistsUserCredentialsByIdAndPassword(credentials)
}

func (m *MongoDbUserCredentialsService) InsertUserCredentials(credentials core.UserCredentials) (core.UserCredentials, error) {
	return m.InsertUserCredentials(credentials)
}

func (m *MongoDbUserCredentialsService) FindUserCredentialsById(credentials core.UserFullCredentials) (core.UserFullCredentials, error) {
	return m.userCredentialsRepository.FindUserCredentialsById(credentials)
}

func (m *MongoDbUserCredentialsService) UpdateUserCredentials(credentials core.UserFullCredentials) error {
	return m.userCredentialsRepository.UpdateUserCredentials(credentials)
}

type MongoDbUserService struct {
	userRepository *MongoDbUserRepository
}

func NewMongoDbUserService(userRepository *MongoDbUserRepository) core.UserPersistenceService {
	return &MongoDbUserService{userRepository: userRepository}
}

func (m *MongoDbUserService) FindUserById(id string) (core.User, error) {
	return m.userRepository.FindUserById(id)
}

func (m *MongoDbUserService) FindUserByEmail(email string) (core.User, error) {
	return m.userRepository.FindUserByEmail(email)
}

func (m *MongoDbUserService) FindUserByPhoneNumber(phoneNumber string) (core.User, error) {
	return m.userRepository.FindUserByPhoneNumber(phoneNumber)
}

func (m *MongoDbUserService) UpdateUser(user core.User) (core.User, error) {
	return m.userRepository.UpdateUser(user)
}
