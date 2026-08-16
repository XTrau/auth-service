package password_test

import (
	"testing"

	"github.com/XTrau/auth-service/internal/auth/password"
)

func TestBcryptHasher(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		hasError bool
	}{
		{
			name:     "numbers",
			text:     "123321",
			hasError: false,
		},
		{
			name:     "text",
			text:     "aSDGsFAddGsa",
			hasError: false,
		},
		{
			name:     "text and numbers",
			text:     "asdDdsaA4634sdSDGif13i",
			hasError: false,
		},
		{
			name:     "too long password",
			text:     "01234asdDdsa567890123456789012asdDdsa3456789012asdDdsa34567asdDdsa8901234567890123asdDdsa4567890123456789012",
			hasError: true,
		},
	}

	hasher := password.NewBcryptHasher(10)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := hasher.Hash(tc.text)

			if err != nil {
				if tc.hasError {
					return
				}

				t.Fatalf("Error on hashing password: %v, password: %v", err.Error(), tc.text)
			}

			if !hasher.Compare(hashedPassword, tc.text) {
				t.Fatalf("False on compare password: %v", tc.text)
			}
		})
	}

}
