package user

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
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
		utils.WriteError(w, http.StatusInternalServerError, err)
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
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	}
}
