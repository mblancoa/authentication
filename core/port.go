package core

// CredentialsPersistenceService defines operations related to the user credentials persistence
type CredentialsPersistenceService interface {

	// ExistsCredentialsByIdAndPassword confirms if exits any credentials like the argument in database
	// credentials must be hashed
	ExistsCredentialsByIdAndPassword(credentials Credentials) (Credentials, bool)

	// InsertCredentials inserts a new user credentials record into database
	// credentials must be hashed
	InsertCredentials(credentials Credentials) (Credentials, error)

	// FindCredentialsById Returns the userFullCredentials found in db. If not found, userFullCredentials will be empty and error not nil
	// id must be hashed and the returned userFullCredentials will be hashed
	FindCredentialsById(id string) (FullCredentials, error)

	// UpdateCredentials modifies the userFullCredentials with the new data
	UpdateCredentials(credentials FullCredentials) error
}

// UserPersistenceService defines operations related to the users persistence
type UserPersistenceService interface {
	// FindUserById returns the user found in db. If not found, user will be empty and error not nil
	// id must be encrypted and the returned user will be encrypted
	FindUserById(id string) (User, error)

	// FindUserByEmail returns the user found in db. If not found, user will be empty and error not nil
	// email must be encrypted and the returned user will be encrypted
	FindUserByEmail(email string) (User, error)

	// FindUserByPhoneNumber returns the user found in db. If not found, user will be empty and error not nil
	// phoneNumber must be encrypted and the returned user will be encrypted
	FindUserByPhoneNumber(phoneNumber string) (User, error)

	// InsertUser inserts the user into the database
	InsertUser(user User) (User, error)

	// UpdateUser modifies the user with the new data
	UpdateUser(user User) (User, error)
}
