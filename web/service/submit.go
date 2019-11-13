package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func submitHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.HTML(w, http.StatusOK, "form", struct {
			NAME       string `json:"name"`
			UNIVERSITY string `json:"university"`
		}{NAME: req.FormValue("username"), UNIVERSITY: req.FormValue("university")})
	}
}
