package common

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
)

// HashPassword return bcrypt hashed password using the default cost(10).
func HashPassword(password string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w:%w", ErrPasswordHashingFailed, err)
	}

	return string(bs), nil

}

// CheckPassword checks if plainPassword matches hashedPassword.
func CheckPassword(plainPassword, hashedPassword string) error {
	log.Println([]byte(plainPassword))
	log.Println(hashedPassword)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	patterns := map[string]string{
		"special": `[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`,
		"lower":   `[a-z]`,
		"upper":   `[A-Z]`,
		"number":  `[0-9]`,
	}

	for name, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Errorf("failed to compile %s pattern: %w", name, err)
			return false
		}
		if !re.MatchString(password) {
			return false
		}
	}
	return true
}
