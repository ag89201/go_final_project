package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func Start(port string, webDir string) error {

	r := chi.NewRouter()

	r.Mount(mountEndpoint, http.FileServer(http.Dir(webDir)))
	r.Get(nextDatePattern, Auth(NextDateHandler))
	r.Post(apiTaskPattern, Auth(PostTaskHandler))
	r.Get(apiTasksPattern, Auth(GetTasksHandler))
	r.Get(apiTaskPattern, Auth(GetTaskHandler))
	r.Put(apiTaskPattern, Auth(PutTaskHandler))
	r.Post(apiTaskPatternDone, Auth(PostDoneTaskHandler))
	r.Delete(apiTaskPattern, Auth(DeleteTaskHandler))
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
