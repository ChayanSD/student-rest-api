package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ChayanSD/student-rest-api/internal/types"
	"github.com/ChayanSD/student-rest-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest , response.GeneralError(err))
		}

		//Validate the request
		if errs := validator.New().Struct(student); errs != nil {
			validateErrs := errs.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest , response.ValidationError(validateErrs))
			return
		} 

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "Ok"})

	}
}


