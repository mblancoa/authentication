package core

type UserCredentialsRepository interface {

	// ExistsUserCredentialsByIdAndPassword confirms if exits any credentials like the argument in database
	// credentials must be hashed
	ExistsUserCredentialsByIdAndPassword(credentials UserCredentials) (UserCredentials, bool)

	// InsertUserCredentials inserts a new user credentials record into database
	// credentials must be hashed
	InsertUserCredentials(credentials UserCredentials)

	// FindUserCredentialsById Returns the userFullCredentials found in db. If not found, userFullCredentials will be empty and error not nil
	// id must be hashed and the returned userFullCredentials will be hashed
	FindUserCredentialsById(credentials UserFullCredentials)

	// UpdateUserCredentials modifies the userFullCredentials with the new data
	UpdateUserCredentials(credentials UserFullCredentials)
}

type UserRepository interface {
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
