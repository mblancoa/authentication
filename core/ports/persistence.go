package ports

import (
	"github.com/mblancoa/authentication/core/domain"
)

var PersistenceContext *persistenceContext = &persistenceContext{}

type persistenceContext struct {
	CredentialsPersistenceService CredentialsPersistenceService
	UserPersistenceService        UserPersistenceService
}

// CredentialsPersistenceService defines operations related to the user credentials persistence
type CredentialsPersistenceService interface {

	// InsertCredentials inserts a new user credentials record into database
	// credentials must be hashed
	InsertCredentials(credentials domain.Credentials) (domain.Credentials, error)

	// FindCredentialsById Returns the userFullCredentials found in db. If not found, userFullCredentials will be empty and error not nil
	// id must be hashed and the returned userFullCredentials will be hashed
	FindCredentialsById(id string) (domain.FullCredentials, error)

	// UpdateCredentials modifies the userFullCredentials with the new data
	UpdateCredentials(credentials domain.FullCredentials) error
}

// UserPersistenceService defines operations related to the users persistence
type UserPersistenceService interface {
	// FindUserById returns the user found in db. If not found, user will be empty and error not nil
	// id must be encrypted and the returned user will be encrypted
	FindUserById(id string) (domain.User, error)

	// FindUserByEmail returns the user found in db. If not found, user will be empty and error not nil
	// email must be encrypted and the returned user will be encrypted
	FindUserByEmail(email string) (domain.User, error)

	// FindUserByPhoneNumber returns the user found in db. If not found, user will be empty and error not nil
	// phoneNumber must be encrypted and the returned user will be encrypted
	FindUserByPhoneNumber(phoneNumber string) (domain.User, error)

	// InsertUser inserts the user into the database
	InsertUser(user domain.User) (domain.User, error)

	// UpdateUser modifies the user with the new data
	UpdateUser(user domain.User) error
}
