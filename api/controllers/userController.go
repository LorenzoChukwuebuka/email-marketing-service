package controllers

import (
	"fmt"
	"net/http"
)




func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}
