package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/anwarminst/students-api/internal/storage"
	"github.com/anwarminst/students-api/internal/types"
	"github.com/anwarminst/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		slog.Info("User Created successfully", slog.String("userId", fmt.Sprint(lastId)))
		if err != nil {
			slog.Info("my_error: ", slog.String("g", fmt.Sprint(err)))
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student by Id: ", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("id not converting into integer : ", slog.String("id ", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, e := storage.RetrieveById(intId)
		if e != nil {
			slog.Error("error getting user : ", slog.String("id ", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(e))
			return
		}
		response.WriteJson(w, http.StatusOK, student)

	}
}
