package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/config"
	"github.com/whitallee/animal-family-backend/docs"
	"github.com/whitallee/animal-family-backend/service/animal"
	"github.com/whitallee/animal-family-backend/service/enclosure"
	"github.com/whitallee/animal-family-backend/service/habitat"
	"github.com/whitallee/animal-family-backend/service/loopmessage"
	"github.com/whitallee/animal-family-backend/service/notification"
	"github.com/whitallee/animal-family-backend/service/species"
	"github.com/whitallee/animal-family-backend/service/task"
	"github.com/whitallee/animal-family-backend/service/user"
	"github.com/whitallee/animal-family-backend/utils"
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

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/health", s.handleHealth).Methods("GET")
	router.HandleFunc("/openapi.json", s.handleOpenAPISpec).Methods("GET")
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// v2 is additive: it runs alongside v1 so the frontend can migrate
	// incrementally. Services are ported to v2 one at a time; those without a
	// RegisterV2Routes call below are still v1-only.
	v2 := router.PathPrefix("/api/v2").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	speciesStore := species.NewStore(s.db, config.Envs.OpenAIAPIKey, config.Envs.S3AssetsBucket, config.Envs.AWSRegion)
	speciesHandler := species.NewHandler(speciesStore, userStore)
	speciesHandler.RegisterRoutes(subrouter)
	speciesHandler.RegisterV2Routes(v2)

	habitatStore := habitat.NewStore(s.db)
	habitatHandler := habitat.NewHandler(habitatStore, userStore)
	habitatHandler.RegisterRoutes(subrouter)
	habitatHandler.RegisterV2Routes(v2)

	enclosureStore := enclosure.NewStore(s.db)
	enclosureHandler := enclosure.NewHandler(enclosureStore, userStore)
	enclosureHandler.RegisterRoutes(subrouter)

	animalStore := animal.NewStore(s.db)
	animalHandler := animal.NewHandler(animalStore, userStore, enclosureStore)
	animalHandler.RegisterRoutes(subrouter)

	// Create notification store and sender
	notificationStore := notification.NewStore(s.db)
	notificationSender := notification.NewNotificationSender(
		notificationStore,
		config.Envs.VAPIDPublicKey,
		config.Envs.VAPIDPrivateKey,
		config.Envs.VAPIDSubject,
	)
	notificationHandler := notification.NewHandler(notificationStore, userStore, notificationSender)
	notificationHandler.RegisterRoutes(subrouter)

	taskStore := task.NewStore(s.db)
	taskHandler := task.NewHandler(taskStore, userStore, animalStore, enclosureStore, notificationSender)
	taskHandler.RegisterRoutes(subrouter)

	loopMessageStore := loopmessage.NewStore(s.db)
	loopMessageHandler := loopmessage.NewHandler(loopMessageStore)
	loopMessageHandler.RegisterRoutes(subrouter)

	var headersOk = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	frontendURL, ok := os.LookupEnv("FRONTEND_URL")
	if !ok {
		log.Fatal("FRONTEND_URL is not set")
	}
	// Allow both localhost and production frontend
	allowedOrigins := []string{frontendURL}
	productionURL, ok := os.LookupEnv("PRODUCTION_FRONTEND_URL")
	if ok && productionURL != "" {
		allowedOrigins = append(allowedOrigins, productionURL)
	}
	var originsOk = handlers.AllowedOrigins(allowedOrigins)
	var methodsOk = handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, handlers.CORS(headersOk, originsOk, methodsOk)(router))
}

// handleOpenAPISpec serves the v2 API contract. It is public and
// unauthenticated so the frontend's client generator can point straight at a
// running backend instead of needing a copy of the file.
func (s *APIServer) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(docs.OpenAPISpec); err != nil {
		log.Printf("failed to write openapi spec: %v", err)
	}
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Ping(); err != nil {
		utils.WriteError(w, http.StatusServiceUnavailable, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
