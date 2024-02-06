package core

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

const (
	EmailNotificationEmailConfirmationTemplate     = "emailEmailConfirmation.txt"
	SMSNotificationPhoneNumberConfirmationTemplate = "smsPhoneConfirmation.txt"
	EmailNotificationOtp                           = "emailOtp.txt"
	SMSNotificationOtp                             = "smsOtp.txt"
)

type CredentialState string

const (
	Active  CredentialState = "Active"
	Blocked                 = "Blocked"
)

type User struct {
	UserId      string `crypt:"true"`
	Email       string `crypt:"true"`
	PhoneNumber string `crypt:"true"`
	Roles       []string
}
type Credentials struct {
	UserId   string `hash:"true"`
	Password string `hash:"true"`
	State    CredentialState
}

type FullCredentials struct {
	Credentials
	LastPasswords []string
}

type CustomClaims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
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
	roles := "nil"
	if u.Roles != nil {
		roles = strings.Join(u.Roles, ",")
	}
	return fmt.Sprintf("User: {\n\tID: %s,\n\tEmail: %s,\n\tPhoneNumber: %s,\n\tRoles: %s,\n}", u.UserId, u.Email, u.PhoneNumber, roles)
}
func (u Credentials) String() string {
	return fmt.Sprintf("Credentials: {\n\tID: %s,\n\tPassword: %s,\n\tState: %s,\n}", u.UserId, u.Password, u.State)
}
