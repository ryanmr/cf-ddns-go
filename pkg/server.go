package pkg

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type IndexData struct {
	CurrentIp          string
	CurrentTime        string
	PreviousIp         string
	PreviousUpdateTime string
	Now                int64
}

type HealthData struct {
	Project string `json:"project"`
	Time    int64  `json:"time"`
}

func Serve() {

	InitState()
	go func() {
		InitCron()
	}()

	r := chi.NewRouter()
	r.Use(RequestLogger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static"))))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		data := HealthData{
			Project: "cf-ddns-go",
			Time:    time.Now().UnixMilli(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		state.Mutex.Lock()
		currentIp := *state.CurrentIp
		previousIp := *state.PreviousIp
		previousUpdateTime := *state.UpdateTime
		data := IndexData{
			Now:                time.Now().UnixMilli(),
			CurrentTime:        time.Now().String(),
			CurrentIp:          currentIp,
			PreviousIp:         previousIp,
			PreviousUpdateTime: previousUpdateTime.String(),
		}
		state.Mutex.Unlock()

		serveTemplate(w, "./index.html.tmpl", data)
	})

	r.Post("/force-update", func(w http.ResponseWriter, r *http.Request) {

		// this is pretend
		UpdateState(fmt.Sprintf("%d", requestID))

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Info().Str("port", port).Msgf("Starting the server at http://localhost:%s", port)
	http.ListenAndServe(":"+port, r)
}

func serveTemplate(w http.ResponseWriter, path string, data any) {
	tmpl, err := template.ParseFiles(filepath.Join("./templates", path))
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
