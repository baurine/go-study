package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var (
	ErrMissingHeader = errors.New("The length of the `Authorization` header is zero.")
)

type Context struct {
	ID       uint64
	Username string
}

func Sign(c Context, secret string) (tokenString string, err error) {
	if secret == "" {
		secret = viper.GetString("jwt_secret")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       c.ID,
		"username": c.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	})
	tokenString, err = token.SignedString([]byte(secret))
	return
}

func ParseRequest(c *gin.Context) (*Context, error) {
	header := c.Request.Header.Get("Authorization")
	secret := viper.GetString("jwt_secret")
	if len(header) == 0 {
		return &Context{}, ErrMissingHeader
	}

	var t string
	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, secret)
}

func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	token, err := jwt.Parse(tokenString, secretFunc(secret))

	// Parse error.
	if err != nil {
		return ctx, err

		// Read the token if it's valid.
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.ID = uint64(claims["id"].(float64))
		ctx.Username = claims["username"].(string)
		return ctx, nil

		// Other errors.
	} else {
		return ctx, err
	}
}

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}
