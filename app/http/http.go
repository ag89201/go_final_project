package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	model "github.com/ag89201/go_final_project/app"
	"github.com/ag89201/go_final_project/app/db"
	"github.com/ag89201/go_final_project/app/domain"

	"github.com/ag89201/go_final_project/app/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

const (
	mountEndpoint      = "/"
	jsonMimeType       = "application/json; charset=UTF-8"
	nextDatePattern    = "/api/nextdate"
	apiTaskPattern     = "/api/task"
	apiTasksPattern    = "/api/tasks"
	apiTaskPatternDone = "/api/task/done"
	apiSigninPattern   = "/api/signin"
	contentTypeHeader  = "Content-Type"
)

func errorResponse(w http.ResponseWriter, errMsg string, err error) {
	data, _ := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%s: %w", errMsg, err).Error()})
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(data)

	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
		return
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(model.DateFormat, r.FormValue("now"))
	if err != nil {
		errorResponse(w, "error parsing date", err)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := domain.GetNextDate(now, date, repeat)

	if err != nil {
		errorResponse(w, "error getting next date", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))

	if err != nil {
		errorResponse(w, "error writing response", err)
		return
	}
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask model.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		errorResponse(w, "error reading request body", err)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		errorResponse(w, "Error parsing JSON", err)
		return
	}

	if err := newTask.CheckCorrectData(); err != nil {
		errorResponse(w, "invalid data", err)
		return
	}

	id, err := db.Database.InsertTask(newTask)
	if err != nil {
		errorResponse(w, "error inserting task", err)
		return
	}

	data, err := json.Marshal(model.IdResponse{Id: id})
	if err != nil {
		errorResponse(w, "error marshaling response", err)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		errorResponse(w, "error writing response", err)
		return
	}

}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Database.GetTasks()
	if err != nil {
		errorResponse(w, "error getting tasks", err)
		return
	}

	if tasks == nil {
		tasks = make([]model.Task, 0)
	}

	data, err := json.Marshal(model.TaskResponse{Tasks: tasks})
	if err != nil {
		errorResponse(w, "error marshaling response", err)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)

	if err != nil {
		errorResponse(w, "error writing response", err)
		return
	}

}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		errorResponse(w, "invalid id", err)
		return
	}

	task, err := db.Database.GetTask(id)

	if err != nil {
		errorResponse(w, "error getting task", err)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		errorResponse(w, "error marshaling response", err)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)

	if err != nil {
		errorResponse(w, "error writing response", err)
		return
	}

}

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		errorResponse(w, "error reading request body", err)
		return
	}

	var task model.Task
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		errorResponse(w, "Error parsing JSON", err)
		return
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		errorResponse(w, "invalid id", err)
		return
	}

	if err := task.CheckCorrectData(); err != nil {
		errorResponse(w, "invalid data", err)
		return
	}

	rowsAffected, err := db.Database.UpdateTask(task)
	if rowsAffected == 0 || err != nil {
		errorResponse(w, "error updating task", err)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		errorResponse(w, "error marshaling response", err)
		return
	}
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		errorResponse(w, "error writing response", err)
		return
	}

}

func PostDoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		errorResponse(w, "invalid id", err)
	}

	task, err := db.Database.GetTask(id)
	if err != nil {
		errorResponse(w, "error getting task", err)
	}

	if len(task.Repeat) == 0 {
		err = db.Database.DeleteTask(id)
		if err != nil {
			errorResponse(w, "error deleting task", err)
			return
		}
	} else {
		task.Date, err = domain.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			errorResponse(w, "error getting next date", err)
			return
		}
	}
	_, err = db.Database.UpdateTask(task)
	if err != nil {
		errorResponse(w, "error updating task", err)
		return
	}

	data, err := json.Marshal(struct{}{})
	if err != nil {
		errorResponse(w, "error marshaling response", err)
		return
	}
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		errorResponse(w, "error writing response", err)
		return
	}

}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		errorResponse(w, "invalid id", err)
		return
	}
	err = db.Database.DeleteTask(id)
	if err != nil {
		errorResponse(w, "error deleting task", err)
		return
	}
	data, err := json.Marshal(struct{}{})
	if err != nil {
		errorResponse(w, "error marshaling response", err)
		return
	}
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		errorResponse(w, "error writing response", err)
		return
	}
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		errorResponse(w, "error reading request body", err)
		return
	}

	var signin model.Sign
	if err := json.Unmarshal(buf.Bytes(), &signin); err != nil {
		errorResponse(w, "Error parsing JSON", err)
		return
	}

	envPass := os.Getenv("TODO_PASSWORD")
	if signin.Password == envPass {
		jwtInstance := jwt.New(jwt.SigningMethodHS256)
		token, err := jwtInstance.SignedString([]byte(envPass))
		if err != nil {
			errorResponse(w, "error signing token", err)
		}

		takIdData, err := json.Marshal(model.AuthToken{Token: token})
		if err != nil {
			errorResponse(w, "error marshaling response", err)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(takIdData)

		if err != nil {
			errorResponse(w, "error writing response", err)
			return
		}
	} else {
		errData, err := json.Marshal(model.ErrorResponse{Error: "wrong password"})
		if err != nil {
			errorResponse(w, "error marshaling response", err)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write(errData)
		if err != nil {
			errorResponse(w, "error writing response", err)
			return
		}
	}
}

func StartServer(port string, webDir string) error {

	r := chi.NewRouter()

	r.Mount(mountEndpoint, http.FileServer(http.Dir(webDir)))
	r.Get(nextDatePattern, middleware.Auth(NextDateHandler))
	r.Post(apiTaskPattern, middleware.Auth(PostTaskHandler))
	r.Get(apiTasksPattern, middleware.Auth(GetTasksHandler))
	r.Get(apiTaskPattern, middleware.Auth(GetTaskHandler))
	r.Put(apiTaskPattern, middleware.Auth(PutTaskHandler))
	r.Post(apiTaskPatternDone, middleware.Auth(PostDoneTaskHandler))
	r.Delete(apiTaskPattern, middleware.Auth(DeleteTaskHandler))
	r.Post(apiSigninPattern, SigninHandler)

	// Start server
	log.Info("Starting server...")
	log.Info("Server listening on port: ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		return err
	}

	return nil
}
