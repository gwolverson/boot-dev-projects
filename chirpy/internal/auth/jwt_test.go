package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userId1 := uuid.New()
	tokenSecret1 := "secret1"
	tokenSecret2 := "secret2"

	tests := []struct {
		name           string
		userId         uuid.UUID
		makeSecret     string // secret used to make the token
		validateSecret string // secret used to validate the token
		duration       time.Duration
		malformed      bool // if true, mess up the token string
		wantErr        bool
	}{
		{
			name:           "Valid token",
			userId:         userId1,
			makeSecret:     tokenSecret1,
			validateSecret: tokenSecret1,
			duration:       time.Hour,
			wantErr:        false,
		},
		{
			name:           "Wrong secret",
			userId:         userId1,
			makeSecret:     tokenSecret1,
			validateSecret: tokenSecret2, // different secret
			duration:       time.Hour,
			wantErr:        true,
		},
		{
			name:           "Expired token",
			userId:         userId1,
			makeSecret:     tokenSecret1,
			validateSecret: tokenSecret1,
			duration:       time.Nanosecond, // expires immediately
			wantErr:        true,
		},
		{
			name:           "Malformed token",
			userId:         userId1,
			makeSecret:     tokenSecret1,
			validateSecret: tokenSecret1,
			duration:       time.Hour,
			malformed:      true, // indicate we want to mess up the token
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt, _ := MakeJWT(tt.userId, tt.makeSecret, tt.duration)
			if tt.malformed {
				jwt = "not.a.validtoken" // deliberately malformed token
			}
			uuid, validateErr := ValidateJWT(jwt, tt.validateSecret)
			if (validateErr != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", validateErr, tt.wantErr)
			}
			if !tt.wantErr && uuid != tt.userId {
				t.Errorf("ValidateJWT() returned wrong UUID. got = %v, want = %v", uuid, tt.userId)
			}
		})
	}
}
