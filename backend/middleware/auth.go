package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		if len(authHeader) > 1000 {
			c.JSON(http.StatusRequestHeaderFieldsTooLarge, gin.H{"error": "Header muito grande"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token malformado"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		parts := strings.Split(tokenString, ".")
		if len(parts) != 3 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token malformado"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		expFloat, ok := claims["exp"].(float64)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		if time.Now().Unix() > int64(expFloat) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Set("exp", claims["exp"])

		c.Next()
	}
}

func TemplateAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		publicRoutes := []string{
			"/login",
			"/register",
			"/static",
			"/api/auth/login",
			"/api/auth/register",
		}

		for _, route := range publicRoutes {
			if strings.HasPrefix(c.Request.URL.Path, route) {
				c.Next()
				return
			}
		}

		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {

			c.HTML(http.StatusOK, "login.html", gin.H{
				"Title":        "Login",
				"FlashMessage": "Por favor, faça login para continuar",
				"FlashType":    "warning",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
			c.Abort()
			return
		}

		c.Next()
	}
}
