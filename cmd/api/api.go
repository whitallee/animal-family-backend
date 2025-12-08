package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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
		// Allow requests from localhost:3000
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
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

	var headersOk = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	frontendURL, ok := os.LookupEnv("FRONTEND_URL")
	if !ok {
		log.Fatal("FRONTEND_URL is not set")
	}
	var originsOk = handlers.AllowedOrigins([]string{frontendURL})
	var methodsOk = handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, corsMiddleware(handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}
