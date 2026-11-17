package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"golang.org/x/crypto/argon2"
)

type parms struct {
	time    uint32
	mermory uint32
	threds  uint8
	keyLen  uint32
	saltLen uint32
}

func HashPassword(plainPassword string) (string, error) {

	p := parms{
		time:    3,
		mermory: 64 * 1024,
		threds:  4,
		keyLen:  32,
		saltLen: 16,
	}

	salt := make([]byte, p.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	hash := argon2.IDKey([]byte(plainPassword), salt, p.time, p.mermory, p.threds, p.keyLen)

	// $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodeHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.mermory,
		p.time,
		p.threds,
		b64Salt,
		b64Hash,
	)

	return encodeHash, nil
}

func VerifyPassword(plainPassword, hashPassword string) (bool, error) {

	// $argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s
	parts := strings.Split(hashPassword, "$")

	if len(parts) != 6 {
		return false, httpx.ErrPasswordLen
	}

	if parts[1] != "argon2id" {
		return false, httpx.ErrInvalidPasswordHash
	}

	// v=%d
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, httpx.ErrArgonVersion
	}

	// m=%d,t=%d,p=%
	p := parms{}
	n, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.mermory, &p.time, &p.threds)

	if err != nil || n != 3 {
		return false, httpx.ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodeHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(plainPassword), salt, p.time, p.mermory, p.threds, uint32(len(decodeHash)))

	return subtle.ConstantTimeCompare(decodeHash, computedHash) == 1, nil
}
