package jwt

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type (
	JWTConfig struct {
		Skipper    Skipper
		SigningKey interface{}
	}

	Skipper      func(c echo.Context) bool
	jwtExtractor func(echo.Context) (string, error)
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	ErrJWTInvalid = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt")
)

var JWTSecret = []byte("secret")

func JWT() echo.MiddlewareFunc {
	c := JWTConfig{}
	c.SigningKey = JWTSecret
	return JWTWithConfig(c)
}

func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	extractor := jwtHeader("Authorization", "Bearer")
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			auth, err := extractor(ctx)
			if err != nil {
				log.Default().Println("err: ", err)
				if config.Skipper != nil {
					if config.Skipper(ctx) {
						return hf(ctx)
					}
				}
				return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
			}
			token, err := jwt.Parse(auth, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Default().Fatal("unexpected signing method: ", t.Header["alg"])
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return config.SigningKey, nil
			})
			if err != nil {
				log.Default().Println("err: ", err)
				return ctx.JSON(http.StatusForbidden, "Forbidden")
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := uint64(claims["userId"].(float64))
				ctx.Set("userId", userID)
				return hf(ctx)
			} else {
				return ctx.JSON(http.StatusForbidden, "Forbidden")
			}
		}
	}
}

func jwtHeader(header string, authScheme string) jwtExtractor {
	return func(ctx echo.Context) (string, error) {
		auth := ctx.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

func GenerateJWT(id uint) string {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString([]byte(JWTSecret))
	return t
}
