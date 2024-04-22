package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	model "github.com/ag89201/go_final_project/app"
	"github.com/ag89201/go_final_project/app/db"
	"github.com/ag89201/go_final_project/app/domain"
	"github.com/go-chi/chi/v5"
)

const (
	mountEndpoint      = "/"
	jsonMimeType       = "application/json; charset=UTF-8"
	nextDatePattern    = "/api/nextdate"
	apiTaskPattern     = "/api/task"
	apiTasksPattern    = "/api/tasks"
	apiTaskPatternDone = "/api/task/done"
	contentTypeHeader  = "Content-Type"
)

func errorResponse(w http.ResponseWriter, errMsg string, err error) {
	data, err := json.Marshal(model.ErrorResponse{Error: errMsg})
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(data)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(model.DateFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := domain.GetNextDate(now, date, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))

	if err != nil {
		http.Error(w, "error writing response", http.StatusBadRequest)
	}
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask model.Task
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if err := newTask.CheckCorrectData(); err != nil {

		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%s", err).Error()})
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, fmt.Errorf("%w", err).Error(), http.StatusBadRequest)
			return
		}

		return
	}

	id, err := db.Database.InsertTask(newTask)
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(model.IdResponse{Id: id})
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Database.GetTasks()
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	if tasks == nil {
		tasks = make([]model.Task, 0)
	}

	data, err := json.Marshal(model.TaskResponse{Tasks: tasks})
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		// http.Error(w, fmt.Errorf("не указан идентификатор").Error(), http.StatusBadRequest)
		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("не указан идентификатор").Error()})
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)
		return
	}

	task, err := db.Database.GetTask(id)

	if err != nil {
		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("задача не найдена").Error()})
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)

		fmt.Printf("%s\n", data)

		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%s\n", data)

}

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("error reading request body").Error()})
		if err != nil {
			log.Fatal(err)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	var task model.Task
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {

		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("error parsing JSON").Error()})
		if err != nil {
			log.Fatal(err)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("the ID must be a number").Error()})
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}
		return
	}

	if err := task.CheckCorrectData(); err != nil {
		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%s", err).Error()})
		if err != nil {
			log.Fatal(err)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	rowsAffected, err := db.Database.UpdateTask(task)
	if rowsAffected == 0 || err != nil {
		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%s", err).Error()})
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set(contentTypeHeader, jsonMimeType)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
			return
		}
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%s\n", data)
}

// func PutTaskHandlerDone(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if err != nil {
// 		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("не указан идентификатор").Error()})
// 		if err != nil {
// 			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
// 			return
// 		}
// 		w.Header().Set(contentTypeHeader, jsonMimeType)
// 		w.WriteHeader(http.StatusBadRequest)
// 		_, err = w.Write(data)
// 	}

// 	task, err := db.Database.GetTask(id)
// 	if err != nil {
// 		data, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("задача не найдена").Error()})
// 		if err != nil {
// 			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
// 			return
// 		}
// 		w.Header().Set(contentTypeHeader, jsonMimeType)
// 		w.WriteHeader(http.StatusBadRequest)
// 		_, err = w.Write(data)
// 	}

// }

func StartServer(port string, webDir string) error {

	r := chi.NewRouter()

	r.Mount(mountEndpoint, http.FileServer(http.Dir(webDir)))
	r.Get(nextDatePattern, NextDateHandler)
	r.Post(apiTaskPattern, PostTaskHandler)
	r.Get(apiTasksPattern, GetTasksHandler)
	r.Get(apiTaskPattern, GetTaskHandler)
	r.Put(apiTaskPattern, PutTaskHandler)
	// r.Post(apiTaskPatternDone, PutTaskHandlerDone)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		return err
	}

	return nil
}
