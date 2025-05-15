package user

import (
	"net/http"
	"os"
	"portarius/internal/user/domain"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserHandler struct {
	repo domain.IUserRepository
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewUserHandler(repo domain.IUserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user with the provided registration data.
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "User registration data"
// @Success 201 {object} domain.User "Created user (password omitted)"
// @Failure 400
// @Failure 500
// @Router /auth/register [post]
func (c *UserHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := c.repo.Create(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	user.Password = ""
	ctx.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Authenticates user and returns JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "User login data"
// @Success 200 {object} map[string]interface{} "JWT token and user info"
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /auth/login [post]
func (c *UserHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// GetAll godoc
// @Summary Get all users with pagination
// @Description Retrieves paginated list of users (passwords omitted).
// @Tags Users
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {array} domain.User "List of users"
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /users/ [get]
func (c *UserHandler) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	users, err := c.repo.GetAll(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	ctx.JSON(http.StatusOK, users)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Retrieves user details by ID (password omitted).
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.User "User found"
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /users/{id} [get]
func (c *UserHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user, err := c.repo.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	user.Password = ""
	ctx.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary Update a user
// @Description Updates the user identified by ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.User true "User data"
// @Success 200 {object} domain.User "Updated user (password omitted)"
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /users/{id} [put]
func (c *UserHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uint(id)
	if err := c.repo.Update(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	ctx.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary Delete a user by ID
// @Description Deletes user with the specified ID.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /users/{id} [delete]
func (c *UserHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuário excluído com sucesso"})
}
