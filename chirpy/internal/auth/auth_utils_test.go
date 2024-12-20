package auth

import (
	"net/http"
	"testing"
)

func TestStripAuthorizationHeader(t *testing.T) {
	tests := []struct {
		token      string
		tokenType  string
		authHeader string
		wantErr    bool
	}{
		{
			token:      "sometoken",
			tokenType:  "Bearer",
			authHeader: "Bearer sometoken",
			wantErr:    false,
		},
		{
			tokenType:  "Bearer",
			authHeader: "sometoken",
			wantErr:    true,
		},
		{
			tokenType:  "Bearer",
			authHeader: "Bearer",
			wantErr:    true,
		},
		{
			tokenType:  "ApiKey",
			token:      "this_is_a_key",
			authHeader: "ApiKey this_is_a_key",
			wantErr:    false,
		},
		{
			tokenType:  "ApiKey",
			authHeader: "ApiKeythis_is_a_key",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.authHeader, func(t *testing.T) {
			headers := http.Header{}
			headers.Add("Authorization", tt.authHeader)
			token, err := StripAuthorizationHeader(headers, tt.tokenType)
			if (err != nil) != tt.wantErr {
				t.Errorf("StripAuthorizationHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.token != token {
				t.Errorf("StripAuthorizationHeader() returned wrong token. got = %v, want = %v", token, tt.token)
			}
		})
	}
}
