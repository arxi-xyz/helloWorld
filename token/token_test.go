package token

import (
	"errors"
	"testing"
)

func TestHash(t *testing.T) {
	got := Hash("hello")

	const expected = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	if got.Hex() != expected {
		t.Fatalf(
			"Hash(hello) = %s, want %s",
			got.Hex(),
			expected,
		)
	}
}

func TestEqual(t *testing.T) {
	a := Hash("token")
	b := Hash("token")
	c := Hash("different-token")

	if !Equal(a, b) {
		t.Fatal("equal digests must compare equal")
	}

	if Equal(a, c) {
		t.Fatal("different digests must not compare equal")
	}
}

func TestCompare(t *testing.T) {
	digest := Hash("secret-token")

	matched, err := Compare("secret-token", digest)
	if err != nil {
		t.Fatalf("Compare() error = %v", err)
	}
	if !matched {
		t.Fatal("expected token to match")
	}

	matched, err = Compare("wrong-token", digest)
	if err != nil {
		t.Fatalf("Compare() error = %v", err)
	}
	if matched {
		t.Fatal("expected token not to match")
	}
}

func TestHexRoundTrip(t *testing.T) {
	original := Hash("some-token")

	encoded := original.Hex()

	decoded, err := ParseHex(encoded)
	if err != nil {
		t.Fatalf("ParseHex() error = %v", err)
	}

	if !Equal(original, decoded) {
		t.Fatalf(
			"round trip mismatch: got %s, want %s",
			decoded.Hex(),
			original.Hex(),
		)
	}
}

func TestParseHexRejectsShortDigest(t *testing.T) {
	_, err := ParseHex("abcd")

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrInvalidHexLength) {
		t.Fatalf(
			"expected ErrInvalidHexLength, got %v",
			err,
		)
	}
}

func TestParseHexRejectsLongDigest(t *testing.T) {
	input := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	_, err := ParseHex(input)

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrInvalidHexLength) {
		t.Fatalf(
			"expected ErrInvalidHexLength, got %v",
			err,
		)
	}
}

func TestParseHexRejectsInvalidHex(t *testing.T) {
	input := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"

	_, err := ParseHex(input)

	if err == nil {
		t.Fatal("expected invalid hex error")
	}
}
