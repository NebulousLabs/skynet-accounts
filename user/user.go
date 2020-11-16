package user

import (
	"crypto/sha256"
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// EmailValidatorRegEx checks if a string is a valid email address.
	// See https://emailregex.com/
	EmailValidatorRegEx = regexp.MustCompile(`^[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,64}$`)

	// ErrInvalidEmail is returned when we encounter an invalid email value.
	ErrInvalidEmail = errors.New("invalid email")
)

type (
	// Email is an email.
	Email string

	// Hash represents a 256bit hash value.
	Hash [32]byte

	// User represents a Skynet user.
	User struct {
		// ID is a hexadecimal string representation of the MongoDB id assigned
		// to this user object. It is auto-generated by Mongo on insert.
		ID        primitive.ObjectID `bson:"_id"`
		FirstName string             `bson:"firstName" json:"firstName"`
		LastName  string             `bson:"lastName" json:"lastName"`
		Email     Email              `bson:"email" json:"email"`
		password  Hash               `bson:"password"`
	}
)

// Validate validates an email address.
func (e Email) Validate() bool {
	return EmailValidatorRegEx.MatchString(string(e))
}

// SetPassword sets the user's password.
func (u *User) SetPassword(pw string) error {
	// TODO Implement
	u.password = sha256.Sum256([]byte(pw))
	return nil
}
