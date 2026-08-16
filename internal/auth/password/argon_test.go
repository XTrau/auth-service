package password_test

import (
	"testing"

	"github.com/XTrau/auth-service/internal/auth/password"
)

func TestArgon2Hasher(t *testing.T) {
	testCases := []struct {
		name string
		text string
	}{
		{
			name: "numbers",
			text: "123321",
		},
		{
			name: "text",
			text: "aSDGsFAddGsa",
		},
		{
			name: "text and numbers",
			text: "asdDdsaA4634sdSDGif13i",
		},
		{
			name: "long password",
			text: "01234asdDdsa567890123456789012asdDdsa3456789012asdDdsa34567asdDdsa8901234567890123asdDdsa4567890123456789012",
		},
	}

	hasher := password.NewArgon2Hasher(password.Argon2DefaultParams())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedText, err := hasher.Hash(tc.text)
			if err != nil {
				t.Fatalf("Error on hashing text: %v", err.Error())
			}

			if hasher.Compare(hashedText, tc.text) != true {
				t.Fatalf("False on compare text: %v, hash: %v", tc.text, hashedText)
			}
		})
	}
}
