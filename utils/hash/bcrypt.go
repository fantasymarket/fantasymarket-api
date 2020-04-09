package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// While these might seem unnececary,
// we created these functions in case:
//	1. we change to a different hashing algorithm
//	2. we want to increase the cost in the future

// GeneratePasswordHash generates a bcrypt hash (salted)
func GeneratePasswordHash(password []byte) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword(password, 12)
	return hash, err
}

// CompareHashAndPassword compares a bcrypt hashed password with its possible
// plaintext equivalent.
func CompareHashAndPassword(hashedPassword []byte, password []byte) (err error) {
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err
}
