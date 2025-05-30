package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/animal"
	"github.com/whitallee/animal-family-backend/service/enclosure"
	"github.com/whitallee/animal-family-backend/service/habitat"
	"github.com/whitallee/animal-family-backend/service/loopmessage"
	"github.com/whitallee/animal-family-backend/service/species"
	"github.com/whitallee/animal-family-backend/service/task"
	"github.com/whitallee/animal-family-backend/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	speciesStore := species.NewStore(s.db)
	speciesHandler := species.NewHandler(speciesStore, userStore)
	speciesHandler.RegisterRoutes(subrouter)

	habitatStore := habitat.NewStore(s.db)
	habitatHandler := habitat.NewHandler(habitatStore, userStore)
	habitatHandler.RegisterRoutes(subrouter)

	enclosureStore := enclosure.NewStore(s.db)
	enclosureHandler := enclosure.NewHandler(enclosureStore, userStore)
	enclosureHandler.RegisterRoutes(subrouter)

	animalStore := animal.NewStore(s.db)
	animalHandler := animal.NewHandler(animalStore, userStore, enclosureStore)
	animalHandler.RegisterRoutes(subrouter)

	taskStore := task.NewStore(s.db)
	taskHandler := task.NewHandler(taskStore, userStore, animalStore, enclosureStore)
	taskHandler.RegisterRoutes(subrouter)

	loopMessageStore := loopmessage.NewStore(s.db)
	loopMessageHandler := loopmessage.NewHandler(loopMessageStore)
	loopMessageHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
