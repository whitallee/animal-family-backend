package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/animal"
	"github.com/whitallee/animal-family-backend/service/enclosure"
	"github.com/whitallee/animal-family-backend/service/habitat"
	"github.com/whitallee/animal-family-backend/service/species"
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

	animalStore := animal.NewStore(s.db)
	animalHandler := animal.NewHandler(animalStore, userStore)
	animalHandler.RegisterRoutes(subrouter)

	enclosureStore := enclosure.NewStore(s.db)
	enclosureHandler := enclosure.NewHandler(enclosureStore, userStore)
	enclosureHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
