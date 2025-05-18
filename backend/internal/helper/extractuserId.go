package helper

import(
	"net/http"
	"errors"
"github.com/golang-jwt/jwt"
)

func ExtractUserId(r *http.Request) (string, error) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid authentication claims")
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("invalid user ID in claims")
	}
	return userId, nil
}