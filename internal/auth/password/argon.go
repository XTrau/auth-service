package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	timeParams   uint32
	memoryParams uint32
	keyLength    uint32
	threads      uint8
	saltLength   int
}

func NewArgon2Params(timeParams, memoryParams, keyLength uint32, threads uint8, saltLength int) Argon2Params {
	return Argon2Params{
		timeParams:   timeParams,
		memoryParams: memoryParams,
		keyLength:    keyLength,
		threads:      threads,
		saltLength:   saltLength,
	}
}

// RFC 9106 profiles
func Argon2DefaultParams() Argon2Params {
	return Argon2Params{
		timeParams:   3,
		memoryParams: 64 * 1024,
		threads:      4,
		keyLength:    32,
		saltLength:   16,
	}
}

// Hasher implementation Argon2id
type Argon2Hasher struct {
	Argon2Params
}

func NewArgon2Hasher(params Argon2Params) *Argon2Hasher {
	return &Argon2Hasher{params}
}

func (h *Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Compute raw hash bytes via IDKey (Argon2id)
	hash := argon2.IDKey([]byte(password), salt, h.timeParams, h.memoryParams, h.threads, h.keyLength)

	// Encode to base64 for string formatting
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format matching standard PHC format
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, h.memoryParams, h.timeParams, h.threads, b64Salt, b64Hash)

	return encoded, nil
}

func (h *Argon2Hasher) Compare(hash, password string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}

	var memory, time, parallelism uint32

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	originalHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// Re-hash input with identical parameters and salt
	newHash := argon2.IDKey([]byte(password), salt, time, memory, byte(parallelism), uint32(len(originalHash)))

	// Use constant-time comparison to thwart timing attacks
	if subtle.ConstantTimeCompare(originalHash, newHash) == 1 {
		return true
	}

	return false
}
