package password_test

import (
	"fmt"
	"testing"

	"github.com/XTrau/auth-service/internal/auth/password"
)

func TestBcryptHasher(t *testing.T) {
	testCases := []struct {
		text string
	}{
		{text: "123321"},
		{text: "asddsa"},
		{text: "@"},
		{text: "ieurtoeuriojgoreoiiv"},
		{text: "0123456789012345678901234567890123456789012345678901234567890123456789012"},
	}

	hasher := password.NewBcryptHasher(10)

	for i, tc := range testCases {
		testName := fmt.Sprintf("Test #%d", i+1)
		t.Run(testName, func(t *testing.T) {
			hashedPassword, err := hasher.Hash(tc.text)

			if err != nil {
				t.Fatalf("Error on hashing password: %v, password: %v", err.Error(), tc.text)
			}

			if !hasher.Compare(hashedPassword, tc.text) {
				t.Fatalf("False on compare password: %v", tc.text)
			}
		})
	}

}
