package password

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//go:embed common-passwords.txt
var f embed.FS

type PasswordValidator struct {
	passwords map[string]struct{}
}

func NewPasswordValidator() (*PasswordValidator, error) {
	passwordListPath := "common-passwords.txt"

	file, err := f.Open(passwordListPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	passwords := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passwords[strings.TrimSpace(strings.ToLower(scanner.Text()))] = struct{}{}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return &PasswordValidator{
		passwords: passwords,
	}, nil
}

func (v *PasswordValidator) MinimumLengthValidator(password string, minLen int) error {
	if len(password) < minLen {
		return fmt.Errorf("password is too short, should be at least %d characters long", minLen)
	}
	return nil
}

func (v *PasswordValidator) UserAttributeSimilarityValidator(password string, username string, email string) error {
	if strings.Contains(username, password) || strings.Contains(email, password) {
		return errors.New("password is too similar to the user attribute")
	}
	return nil
}

func (v *PasswordValidator) NumericPasswordValidator(password string) error {
	floatVal, err := strconv.ParseFloat(password, 64)
	if err == nil && floatVal != 0 {
		return errors.New("the password cannot be entirely numeric")
	}
	return nil
}

func (v *PasswordValidator) CommonPasswordValidator(password string) error {
	if _, exists := v.passwords[strings.ToLower(strings.TrimSpace(password))]; exists {
		return errors.New("password is too common")
	}
	return nil
}
