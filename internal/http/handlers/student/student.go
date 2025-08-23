package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alimulhuq/students-api/internal/storage"
	"github.com/alimulhuq/students-api/internal/types"
	"github.com/alimulhuq/students-api/internal/types/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.Writejson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.Writejson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("Creating a student detail")

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.Writejson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User Created Successfully", slog.String("User ID:", fmt.Sprint(lastId)))

		if err != nil {
			response.Writejson(w, http.StatusInternalServerError, err)
			return
		}

		response.Writejson(w, http.StatusCreated, map[string]any{"success": "OK", "id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			response.Writejson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intID)

		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.Writejson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.Writejson(w, http.StatusOK, student)
	}
}
