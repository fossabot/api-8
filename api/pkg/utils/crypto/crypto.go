package crypto

import "golang.org/x/crypto/bcrypt"

func CreateStringHash(plain string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	return string(bytes), err
}

func ValidateHash(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
