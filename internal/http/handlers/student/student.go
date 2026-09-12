package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ashish-goyals/students-api/internal/storage"
	"github.com/ashish-goyals/students-api/internal/types"
	"github.com/ashish-goyals/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		slog.Info("Creating a student")

		// Request validation
		if err := validator.New().Struct(student); err != nil {
			var validateErrs validator.ValidationErrors
			if errors.As(err, &validateErrs) {
				response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
				return
			}
			// response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		lastId , err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("User Created succesfully", slog.String("userId",fmt.Sprint(lastId)))
if err != nil{
response.WriteJson(w, http.StatusInternalServerError, err)
return
}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id":lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){

		id := r.PathValue("id")

		slog.Info("getting a students by id", slog.String("id", id))
        inId, err := strconv.ParseInt(id, 10, 64)
		
		if err != nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentbyId(inId)

		if err != nil{

			slog.Error("error Geting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		
			return 
		}

		response.WriteJson(w, http.StatusOK,student)
	}
}