package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SecretKey = []byte("SECRET")

func CreateToken(email string, id int, role string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["email"] = email
    claims["id"] = id
    claims["role"] = role
    claims["email"] = email
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    tokenString, err := token.SignedString(SecretKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authorization := c.GetHeader("Authorization")
		if authorization == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

		tokenString := strings.Split(authorization, " ")[1]

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return SecretKey, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims"})
            return
        }

        email := claims["email"].(string)
        id := claims["id"].(float64)
        role := claims["email"].(string)
        c.Set("email", email)
        c.Set("id", id)
        c.Set("role", role)
        c.Next()
    }
}