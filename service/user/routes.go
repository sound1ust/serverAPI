package user

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"serverAPI/service/auth"
	"serverAPI/types"
	"serverAPI/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/test", h.handleTest).Methods("GET")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v: %v", r.Method, r.RequestURI)
	response := []byte("Hello, world!")
	_, err := w.Write(response)
	if err != nil {
		log.Printf("Error while handling /test request: %v", err)
	}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	if err := utils.Validate.Struct(payload); err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		utils.WriteError(w, http.StatusBadRequest, errs)
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if errors.Is(err, NotFoundError) {
		utils.WriteError(
			w, http.StatusBadRequest,
			fmt.Errorf("invalid password or email"),
		)
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("server error occured"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid password or email"))
		return
	}
	if err = utils.WriteJSON(w, http.StatusOK, nil); err != nil {
		http.Error(w, "server error occurred", http.StatusInternalServerError)
	}
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errs))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(
			w, http.StatusConflict,
			fmt.Errorf("user with email '%s' already exists", payload.Email),
		)
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("server error occured"))
	}

	err = h.store.CreateUser(&types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	if err = utils.WriteJSON(w, http.StatusCreated, nil); err != nil {
		http.Error(w, "server error occurred", http.StatusInternalServerError)
	}
}
