package user

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"serverAPI/service/auth"
	"serverAPI/types"
	"testing"
	"time"
)

func TestUserServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("valid payload", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "test",
			LastName:  "test",
			Email:     "valid@mail.com",
			Password:  "test",
		}
		marshalled, _ := json.Marshal(payload)
		request, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusCreated {
			t.Errorf(
				"excpected status code %d, but got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})
	t.Run("invalid email", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "test",
			LastName:  "test",
			Email:     "invalid",
			Password:  "test",
		}
		marshalled, _ := json.Marshal(payload)
		request, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"excpected status code %d, but got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})
}

func TestUserLoginHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)
	tests := []struct {
		name    string
		payload types.LoginUserPayload
		want    int
		wantErr bool
	}{
		{
			name: "test valid",
			payload: types.LoginUserPayload{
				Email:    "test@mail.com",
				Password: "test",
			},
			want:    http.StatusOK,
			wantErr: false,
		},
		{
			name: "test invalid",
			payload: types.LoginUserPayload{
				Email:    "test@mail.com",
				Password: "invalid",
			},
			want:    http.StatusUnauthorized,
			wantErr: false,
		},
		{
			name: "test not found",
			payload: types.LoginUserPayload{
				Email:    "invalid@mail.com",
				Password: "test",
			},
			want:    http.StatusBadRequest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshalled, _ := json.Marshal(tt.payload)
			request, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			recorder := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/login", handler.handleLogin)
			router.ServeHTTP(recorder, request)
			if recorder.Code != tt.want {
				t.Errorf(
					"excpected status code %d, but got %d",
					tt.want,
					recorder.Code,
				)
			}
		})
	}
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if email != "test@mail.com" {
		return nil, NotFoundError
	}
	password, _ := auth.HashPassword("test")
	return &types.User{
		ID:        1,
		FirstName: "test",
		LastName:  "test",
		Email:     "test@mail.com",
		Password:  password,
		Created:   time.Time{},
	}, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(*types.User) error {
	return nil
}
