package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ag89201/go_final_project/app/model"

	"github.com/ag89201/go_final_project/app/domain"
	"github.com/golang-jwt/jwt"
)

func errorInternalResponse(w http.ResponseWriter,  err error){

	_, filename, line, _ := runtime.Caller(1)
	log.Errorf("[%s:%d] %s", filename, line, err.Error())
	if w != nil {
		http.Error(w, fmt.Errorf("internal server error").Error(), http.StatusInternalServerError)
	}
	
}

func errorResponse(w http.ResponseWriter, errMsg string, err error) {
	data, _ := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%s: %w", errMsg, err).Error()})
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(data)

	if err != nil {
		errorInternalResponse(nil, err)
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
		log.Error(err)
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

	id, err := model.Database.InsertTask(newTask)
	if err != nil {
		errorInternalResponse(w, err)
		return
	}

	data, err := json.Marshal(model.IdResponse{Id: id})
	if err != nil {
		errorInternalResponse(w,err)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		errorInternalResponse(nil, err)
		return
	}

}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []model.Task
	search := r.URL.Query().Get("search")

	if len(search) > 0 {
		date, err := time.Parse(model.SearchDateFormat, search)

		if err != nil {
			tasks, err = model.Database.GetTasksByTitleOrComment(search)
			if err != nil {
				errorInternalResponse(w,err)
				return
			}
		} else {
			// search by date
			tasks, err = model.Database.GetTasksByDate(date.Format(model.DateFormat))
			if err != nil {
				errorInternalResponse(w,err)
				return
			}
		}
	} else {
		var err error
		if tasks, err = model.Database.GetTasks(); err != nil {
			errorInternalResponse(w,err)
			return
		}
	}

	if tasks == nil {
		tasks = make([]model.Task, 0)
	}

	data, err := json.Marshal(model.TaskResponse{Tasks: tasks})
	if err != nil {
		errorInternalResponse(w,err)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)

	if err != nil {
		errorInternalResponse(nil, err)
		return
	}

}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		errorResponse(w, "invalid id", err)
		return
	}

	task, err := model.Database.GetTask(id)

	if err != nil {
		if err == sql.ErrNoRows {
		    errorResponse(w, "task was not found", err)
			return
		}
		errorInternalResponse(w,err)
		return		
	}

	data, err := json.Marshal(task)
	if err != nil {
		errorInternalResponse(w,err)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)

	if err != nil {
		errorInternalResponse(nil, err)
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

	rowsAffected, err := model.Database.UpdateTask(task)
	if err != nil {
	    errorInternalResponse(w,err)
		return
	}
	if rowsAffected == 0 {
		errorResponse(w, "task was not found", err)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		errorInternalResponse(w,err)
		return
	}
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		errorInternalResponse(nil,err)
		
		return
	}

}

func PostDoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		errorResponse(w, "invalid id", err)
	}

	task, err := model.Database.GetTask(id)
	
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, "task was not found", err)
			return
		}
		errorInternalResponse(w,err)
		return		    
		}		
	

	if len(task.Repeat) == 0 {
		rows,err := model.Database.DeleteTask(id)
		if err != nil {
			errorInternalResponse(w, err)
			return
		}
		if rows == 0 {
		    errorResponse(w, "task was not found", err)
			return
		}
	} else {
		task.Date, err = domain.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			errorResponse(w, "error getting next date", err)
			return
		}
	}
	_, err = model.Database.UpdateTask(task)
	if err != nil {
		errorResponse(w, "error updating task", err)
		return
	}

	data, err := json.Marshal(struct{}{})
	if err != nil {
		errorInternalResponse(w,err)
		return
	}
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		errorInternalResponse(nil,err)
		return
	}

}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		errorResponse(w, "invalid id", err)
		return
	}

	rows,err := model.Database.DeleteTask(id)
	if err != nil {
		errorInternalResponse(w, err)
		return
	}
	if rows == 0 {
		errorResponse(w, "task was not found", err)
		return
	}

	data, err := json.Marshal(struct{}{})
	if err != nil {
		errorInternalResponse(w,err)
		return
	}
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		errorInternalResponse(nil, err)
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
			errorInternalResponse(w,err)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(takIdData)

		if err != nil {
			errorInternalResponse(nil,err)
			return
		}
	} else {
		errData, err := json.Marshal(model.ErrorResponse{Error: "wrong password"})
		if err != nil {
			errorInternalResponse(w,err)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write(errData)
		if err != nil {
			errorInternalResponse(nil,err)
			return
		}
	}
}
