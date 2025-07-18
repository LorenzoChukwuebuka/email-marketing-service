package helper

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
)

func ExtractUserId(r *http.Request) (string, string, error) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid authentication claims")
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return "", "", errors.New("invalid user ID in claims")
	}
	companyId, ok := claims["company_id"].(string)
	if !ok {
		return "", "", errors.New("invalid company ID in claims")
	}
	return userId, companyId, nil
}


func ExtractAdminDetails(r *http.Request) (string, string, error) {
    claims, ok := r.Context().Value("adminclaims").(jwt.MapClaims)
    if !ok {
        return  "", "", errors.New("invalid authentication claims")
    }
    
    userId, ok := claims["userId"].(string)
    if !ok {
        return  "", "", errors.New("invalid user ID in claims")
    }
    
    adminType, ok := claims["type"].(string)
    if !ok {
        return  "", "", errors.New("invalid admin type in claims")
    }
    
    return userId, adminType, nil
}