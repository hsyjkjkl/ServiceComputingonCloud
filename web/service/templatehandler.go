package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func templateHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.HTML(w, http.StatusOK, "index", struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "17343038", Content: "Hello Web!"})
	}
}
