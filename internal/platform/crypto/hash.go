package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

// --- Password hashing (argon2id) ---

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	argonSaltLen = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	enc := base64.RawStdEncoding.EncodeToString(salt)
	kenc := base64.RawStdEncoding.EncodeToString(key)
	return fmt.Sprintf("argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads, enc, kenc), nil
}

func VerifyPassword(encoded, password string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 5 {
		return false, errors.New("crypto: malformed hash")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	got := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return hmac.Equal(got, want), nil
}

// --- Random ID / token generation ---

func RandomID(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

func NewID() string { return RandomID(12) }

func NewToken(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// Random6 returns a random 6-digit number as a string.
func Random6() int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}
	return n.Int64()
}

// HashToken produces a stable hash for storing token secrets at rest.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawStdEncoding.EncodeToString(sum[:])
}

// --- TOTP (RFC 6238) ---

// GenerateTOTPSecret returns a random base32 TOTP secret (no padding).
func GenerateTOTPSecret() string {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

// TOTPCode computes the current 6-digit code for a secret.
func TOTPCode(secret string, t int64) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha1.New, key)
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(t))
	mac.Write(msg[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	value := (int(sum[offset])&0x7f)<<24 |
		(int(sum[offset+1])&0xff)<<16 |
		(int(sum[offset+2])&0xff)<<8 |
		(int(sum[offset+3]) & 0xff)
	code := value % 1000000
	return fmt.Sprintf("%06d", code), nil
}

// VerifyTOTP checks a code against the current ±1 windows (30s steps).
func VerifyTOTP(secret, code string) bool {
	now := time.Now().Unix() / 30
	for _, delta := range []int64{0, -1, 1} {
		want, err := TOTPCode(secret, now+delta)
		if err == nil && hmac.Equal([]byte(want), []byte(code)) {
			return true
		}
	}
	return false
}

// --- backup codes ---

func GenerateBackupCodes(n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 6)
		if _, err := rand.Int(rand.Reader, big.NewInt(10)); err == nil {
			for j := range b {
				n, _ := rand.Int(rand.Reader, big.NewInt(10))
				b[j] = byte(n.Int64() + '0')
			}
		}
		out = append(out, fmt.Sprintf("%s-%s", string(b[:3]), string(b[3:])))
	}
	return out
}
