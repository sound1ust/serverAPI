package user

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"serverAPI/types"
	"testing"
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

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(*types.User) error {
	return nil
}
