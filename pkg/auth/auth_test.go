package auth

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashPass(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "digits",
			password: "43621967",
		},
		{
			name:     "alphas",
			password: "TUiUKJZy",
		},
		{
			name:     "random",
			password: "7#M1*Z0O",
		},
		{
			name:     "long random",
			password: "&SA%l0O8l5(H8u7#",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed := HashPass(tt.password)
			require.NotEqual(t, tt.password, hashed)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		wantErr bool
	}{
		{
			name:    "alphas",
			login:   "Login",
			wantErr: false,
		},
		{
			name:    "alphasAndDigits",
			login:   "Login1234",
			wantErr: false,
		},
		{
			name:    "email",
			login:   "example@company.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.login)
			require.NoError(t, err)
			actual, err := ParseToken(token)
			require.NoError(t, err)
			require.Equal(t, tt.login, actual)
		})
	}
}

func TestParseToken(t *testing.T) {
	tests := []struct {
		name     string
		rawToken string
		want     string
		wantErr  bool
	}{
		{
			name:     "positive",
			rawToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOm51bGwsImxvZ2luIjoiZXhhbXBsZUBjb21wYW55LmNvbSJ9.L6jChBPp83QP02J9MCpy1WAF3tRQIU06BfzjFUMCIhA",
			want:     "example@company.com",
			wantErr:  false,
		},
		{
			name:     "negative",
			rawToken: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOm51bGwsImxvZ2luIjoiZXhhbXBsZUBjb21wYW55LmNvbSJ9.A0OEPfu3bZnl71Y7oFjxWiPs1eKwOjreb72eZC_3La08IsIf5wrYjqHNB-EuuyHX",
			want:     "example@company.com",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			login, err := ParseToken(tt.rawToken)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, login)
			}
		})
	}
}
