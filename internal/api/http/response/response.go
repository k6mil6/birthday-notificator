package response

import (
	"github.com/go-chi/render"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	resp := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    statusCode,
			"message": message,
		},
	}

	render.Status(r, statusCode)
	render.JSON(w, r, resp)
}

func HandleSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}
