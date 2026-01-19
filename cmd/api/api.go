package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/config"
	"github.com/whitallee/animal-family-backend/service/animal"
	"github.com/whitallee/animal-family-backend/service/enclosure"
	"github.com/whitallee/animal-family-backend/service/habitat"
	"github.com/whitallee/animal-family-backend/service/loopmessage"
	"github.com/whitallee/animal-family-backend/service/notification"
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

	// Create notification store and sender
	notificationStore := notification.NewStore(s.db)
	notificationSender := notification.NewNotificationSender(
		notificationStore,
		config.Envs.VAPIDPublicKey,
		config.Envs.VAPIDPrivateKey,
		config.Envs.VAPIDSubject,
	)
	notificationHandler := notification.NewHandler(notificationStore, notificationSender)
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
	var originsOk = handlers.AllowedOrigins([]string{frontendURL})
	var methodsOk = handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, handlers.CORS(headersOk, originsOk, methodsOk)(router))
}
