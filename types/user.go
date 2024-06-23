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

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := make(map[string]string)

	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("first name must be at least %d characters long", minFirstNameLength)
	}

	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("last name must be at least %d characters long", minLastNameLength)
	}

	if len(params.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters long", minPasswordLength)
	}

	if res := isEmailValid(params.Email); !res {
		errors["email"] = fmt.Sprintf("invalid email address")
	}

	return errors
}

func isEmailValid(email string) bool {
	return regexp.MustCompile(emailRegex).MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
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

func EncryptPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), nil
}
