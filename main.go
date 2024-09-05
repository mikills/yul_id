package main

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	ErrorInvalidInput = errors.New("input should be exactly four alphabetic characters")
)

const (
	prefixLen    = 4 // prefixLen represents the length of the prefix in a YULID
	separatorLen = 1 // separator is denoted by a hyphen '-' in a YULID
	minSuffixLen = 4 // minSuffixLen represents the minimum length of the random part in a YULID after the separator
	maxSuffixLen = 6 // maxSuffixLen represents the maximum length of the random part in a YULID after the separator
	alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// YULID represents a distinct, human-readable identifier for a Yul customer.
// This identifier is designed to be shared easily and combines a user-specific prefix with a random, alphanumeric suffix.
//
// The format of a YULID is: [4-character prefix] + '-' + [4-6 character random string].
//
// Example:
// For a user with the full name "John Doe", their YULID might be "JNDE-ED24HS".
// - "JNDE" is the prefix derived from the user's name.
// - "ED24HS" is the random alphanumeric suffix generated for uniqueness.
type YULID [prefixLen + separatorLen + maxSuffixLen]byte

// String implements the Stringer interface for YULID
func (yd YULID) String() string {
	for i, b := range yd {
		if b == 0 {
			return string(yd[:i])
		}
	}
	return string(yd[:])
}

func New(prefix string) (YULID, error) {
	var yulid YULID
	final := make([]byte, prefixLen+separatorLen+maxSuffixLen)
	if len(prefix) != prefixLen {
		return YULID{}, ErrorInvalidInput
	}

	// append string to final
	for i, r := range prefix {
		if !isAlphanumeric(r) {
			return YULID{}, ErrorInvalidInput
		}
		final[i] = byte(r)
	}

	// append separator
	final[prefixLen] = '-'

	// append random part
	randomPart := generateSuffix()

	// append random part to final
	for i, b := range randomPart {
		final[prefixLen+separatorLen+i] = b
	}

	// copy final to YULID
	copy(yulid[:], final)

	return yulid, nil
}

func generateSuffix() []byte {
	// set up random part
	randomPart := make([]byte, maxSuffixLen)
	max := big.NewInt(int64(len(alphanumeric)))

	// generate random alphanumeric characters
	for i := range randomPart {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		randomPart[i] = alphanumeric[n.Int64()]
	}

	return randomPart

}

func isAlphanumeric(b rune) bool {
	return (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// Validate checks if a YULID is correctly formatted
func Validate(id YULID) error {
	// Ensure length is correct
	ydLen := len(id.String())
	if ydLen < prefixLen+separatorLen+minSuffixLen || ydLen > prefixLen+separatorLen+maxSuffixLen {
		return errors.New("YULID has an invalid length")
	}

	// Check that the prefix is alphanumeric
	for i := 0; i < prefixLen; i++ {
		if !isAlphanumeric(rune(id[i])) {
			return errors.New("YULID has an invalid prefix")
		}
	}

	// Check that the separator is a hyphen
	if id[prefixLen] != '-' {
		return errors.New("YULID separator is invalid")
	}

	// Check that the suffix part is alphanumeric
	for i := prefixLen + separatorLen; i < ydLen; i++ {
		if !isAlphanumeric(rune(id[i])) {
			return errors.New("YULID random part contains invalid characters")
		}
	}

	return nil
}
