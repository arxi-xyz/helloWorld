package token

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	Size    = sha256.Size     // 32 bytes
	HexSize = sha256.Size * 2 // 64 hex characters
)

var (
	ErrInvalidHexLength = errors.New("tokenhash: invalid hex digest length")
)

type Digest [Size]byte

func Hash(token string) Digest {
	return Digest(sha256.Sum256([]byte(token)))
}

func Equal(d1, d2 Digest) bool {
	// return d1 == d2
	return subtle.ConstantTimeCompare(d1[:], d2[:]) == 1
}

func Compare(token string, hexDigest Digest) (bool, error) {
	if len(hexDigest) >= HexSize {
		return false, ErrInvalidHexLength
	}

	digest := Hash(token)
	return Equal(digest, hexDigest), nil
}

func (d Digest) String() string {
	return d.Hex()
}

func (d Digest) Hex() string {
	return hex.EncodeToString(d[:])
}

func ParseHex(s string) (Digest, error) {
	var digest Digest

	if len(s) != HexSize {
		return digest, fmt.Errorf(
			"%w: got %d characters, want %d",
			ErrInvalidHexLength,
			len(s),
			HexSize,
		)
	}

	if _, err := hex.Decode(digest[:], []byte(s)); err != nil {
		return Digest{}, fmt.Errorf(
			"tokenhash: invalid hex digest: %w",
			err,
		)
	}

	return digest, nil
}
