package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost         = bcrypt.DefaultCost
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
	emailRegex         = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func NewUser(params *CreateUserParams) (*User, error) {
	encryptedPassword, err := EncryptPassword(params.Password)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: encryptedPassword,
	}, nil
}

func (params CreateUserParams) Validate() map[string]string {
	errors := make(map[string]string)

	validateFirstName(errors, params.FirstName)
	validateLastName(errors, params.LastName)
	validatePassword(errors, params.Password)
	validateEmail(errors, params.Email)

	return errors
}

func (params UpdateUserParams) Validate() map[string]string {
	errors := make(map[string]string)

	if params.FirstName != "" {
		validateFirstName(errors, params.FirstName)
	}
	if params.LastName != "" {
		validateLastName(errors, params.LastName)
	}

	return errors
}

func validateFirstName(errors map[string]string, firstName string) {
	if len(firstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("first name must be at least %d characters long", minFirstNameLength)
	}
}

func validateLastName(errors map[string]string, lastName string) {
	if len(lastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("last name must be at least %d characters long", minLastNameLength)
	}
}

func validatePassword(errors map[string]string, password string) {
	if len(password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters long", minPasswordLength)
	}
}

func validateEmail(errors map[string]string, email string) {
	if res := isEmailValid(email); !res {
		errors["email"] = fmt.Sprintf("invalid email address")
	}
}

func isEmailValid(email string) bool {
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func EncryptPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), nil
}
