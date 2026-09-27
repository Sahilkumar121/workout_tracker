package handler

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestRegisterUser(t *testing.T) {

	validate := validator.New()
	tests := []struct {
		name    string
		user    RegisterUser
		wantErr bool
	}{
		{
			name: "Valid register",
			user: RegisterUser{
				FirstName: "test",
				LastName:  "user",
				Email:     "testuser@gmail.com",
				Password:  "testdummyuser123",
			},
			wantErr: false,
		},
		{
			name: "Invalid first name (Too short)",
			user: RegisterUser{
				FirstName: "te",
				LastName:  "user",
				Email:     "testuser@gmail.com",
				Password:  "testdummyuser123",
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			user: RegisterUser{
				FirstName: "test",
				LastName:  "user",
				Email:     "not-a-email",
				Password:  "testdummyuser123",
			},
			wantErr: true,
		},
		{
			name: "Invalid password (Too short)",
			user: RegisterUser{
				FirstName: "test",
				LastName:  "user",
				Email:     "testuser@gmail.com",
				Password:  "testuser",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got %v \n", tt.wantErr, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {

	validate := validator.New()

	tests := []struct {
		name    string
		user    LoginUser
		wantErr bool
	}{
		{
			name: "Valid login",
			user: LoginUser{
				Email:    "testuser@gmail.com",
				Password: "testuserpassword123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: LoginUser{
				Email:    "not-a-email",
				Password: "testuserpassword123",
			},
			wantErr: true,
		},
		{
			name: "invalid password (Too short)",
			user: LoginUser{
				Email:    "testuser@gmail.com",
				Password: "testuser",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v \n", tt.wantErr, err)
			}
		})
	}
}
