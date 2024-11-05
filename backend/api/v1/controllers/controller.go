package controllers

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
	ErrConflict     = errors.New("conflict")
)

// Create a helper method to extract user ID from claims
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

// Create a centralized method for parsing pagination and search parameters
func ParsePaginationParams(r *http.Request) (page, pageSize int, searchQuery string, err error) {
	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery = r.URL.Query().Get("search")
	page, err = strconv.Atoi(page1)
	if err != nil {
		return 0, 0, "", errors.New("invalid page number")
	}
	pageSize, err = strconv.Atoi(pageSize1)
	if err != nil {
		return 0, 0, "", errors.New("invalid page size")
	}
	return page, pageSize, searchQuery, nil
}

// Create a centralized error handling method
func HandleControllerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.ErrorResponse(w, err.Error())
	case errors.Is(err, ErrUnauthorized):
		response.ErrorResponse(w, err.Error())
	default:
		response.ErrorResponse(w, err.Error())
	}
}
