package domain

import "golang.org/x/crypto/bcrypt"

func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

func ValidatePassword(pass string) error {
	if len(pass) < 8 {
		return ErrWeakPassword
	}
	return nil
}

func ValidateAndHashPassword(plain string) (string, error) {
	if err := ValidatePassword(plain); err != nil {
		return "", err
	}

	return HashPassword(plain)
}

func checkPassword(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
