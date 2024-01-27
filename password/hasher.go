package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordHasher struct {
	Time       uint32
	Memory     uint32
	Threads    uint8
	KeyLength  uint32
	SaltLength uint32
}

func (ph *PasswordHasher) GenerateFromPassword(password string) (hash string, err error) {
	salt := make([]byte, ph.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, ph.Time, ph.Memory, ph.Threads, ph.KeyLength)

	return fmt.Sprintf("argon2$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", ph.Memory, ph.Time, ph.Threads,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hashed)), nil
}

func (ph *PasswordHasher) CompareHashAndPassword(hash, password string) (err error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return errors.New("error: hashed password is incorrectly formatted")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	decodedHashed, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return err
	}

	comparisonHashed := argon2.IDKey([]byte(password), salt, ph.Time, ph.Memory, ph.Threads, ph.KeyLength)
	if subtle.ConstantTimeCompare(decodedHashed, comparisonHashed) != 1 {
		return errors.New("passwords do not match")
	}

	return nil
}

func NewDjangoArgon2PasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		Time:       2,
		Memory:     102400,
		Threads:    8,
		KeyLength:  32,
		SaltLength: 16,
	}
}
