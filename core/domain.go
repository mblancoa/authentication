package core

import (
	"fmt"
)

type CredentialState string

const (
	Active  CredentialState = "Active"
	Blocked                 = "Blocked"
)

type User struct {
	ID          string `crypt:"true"`
	Email       string `crypt:"true"`
	PhoneNumber string `crypt:"true"`
}
type UserCredentials struct {
	ID       string `hash:"true"`
	Password string `hash:"true"`
	State    CredentialState
}

type UserFullCredentials struct {
	UserCredentials
	LastPasswords []string
}

type ConfirmationCredentials struct {
	User
	Password string
	Otp      string
}

type StrongCustomerAuthentication struct {
	Otp string
}

func (u User) String() string {
	return fmt.Sprintf("User: {\n\tID: %s,\n\tEmail: %s,\n\tPhoneNumber: %s,\n}", u.ID, u.Email, u.PhoneNumber)
}
func (u UserCredentials) String() string {
	return fmt.Sprintf("UserCredentials: {\n\tID: %s,\n\tPassword: %s,\n\tState: %s,\n}", u.ID, u.Password, u.State)
}
