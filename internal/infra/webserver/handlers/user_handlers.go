package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github/GuilhermeHermes/GO_API/internal/dto"
	"github/GuilhermeHermes/GO_API/internal/entity"
	"github/GuilhermeHermes/GO_API/internal/infra/database"

	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB        database.UserDB
	Jwt           *jwtauth.JWTAuth
	JwtExpiration int64
}

func NewUserHandler(db database.UserDB, jwt *jwtauth.JWTAuth, jwtExpiration int64) *UserHandler {
	return &UserHandler{
		UserDB:        db,
		Jwt:           jwt,
		JwtExpiration: jwtExpiration,
	}
}

func (h *UserHandler) GetJwt(w http.ResponseWriter, r *http.Request) {
	var user dto.GetJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	existingUser, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Usar o método CheckPassword para validar a senha criptografada
	if !existingUser.CheckPassword(user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	_, tokenString, err := h.Jwt.Encode(map[string]interface{}{
		"sub": existingUser.ID,
		"exp": time.Now().Add(time.Duration(h.JwtExpiration) * time.Second).Unix(),
	})
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	accesstoken := map[string]string{
		"token": tokenString,
	}

	json.NewEncoder(w).Encode(accesstoken)
}

// CreateUser cria um novo usuário
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userReq dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validações básicas
	if strings.TrimSpace(userReq.Username) == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(userReq.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(userReq.Password) == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Criar usuário
	user, err := entity.NewUser(userReq.Username, userReq.Email, userReq.Password, userReq.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.UserDB.Create(user); err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remover senha da resposta
	userResponse := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

// GetUserByEmail busca um usuário por email
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Remover senha da resposta
	userResponse := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

// GetUserByID busca um usuário por ID
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL (assumindo padrão /users/{id})
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	id := parts[len(parts)-1]

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Remover senha da resposta
	userResponse := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

// UpdateUser atualiza um usuário
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	id := parts[len(parts)-1]

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Verificar se o usuário existe
	existingUser, err := h.UserDB.FindByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updateReq dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Atualizar os campos se fornecidos
	if updateReq.Username != "" {
		existingUser.Username = updateReq.Username
	}
	if updateReq.Role != "" {
		existingUser.Role = updateReq.Role
	}

	if err := h.UserDB.Update(existingUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remover senha da resposta
	userResponse := map[string]interface{}{
		"id":         existingUser.ID,
		"username":   existingUser.Username,
		"email":      existingUser.Email,
		"role":       existingUser.Role,
		"created_at": existingUser.CreatedAt,
		"updated_at": existingUser.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

// DeleteUser deleta um usuário
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	id := parts[len(parts)-1]

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Verificar se o usuário existe
	_, err := h.UserDB.FindByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := h.UserDB.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CheckUserExists verifica se um usuário existe por email
func (h *UserHandler) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	exists, err := h.UserDB.Exists(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]bool{
		"exists": exists,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
