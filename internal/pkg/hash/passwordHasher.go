package hash

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct {
}

func (p *PasswordHasher) Hash(data []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(data, 12)
	return hash, err
}

func (p *PasswordHasher) Verify(data []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, data)
	return err == nil
}

func NewPasswordHasher() Hasher {
	return &PasswordHasher{}
}
