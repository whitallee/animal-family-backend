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
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(subrouter)

	speciesStore := species.NewStore(s.db)
	speciesService := species.NewHandler(speciesStore)
	speciesService.RegisterRoutes(subrouter)

	habitatStore := habitat.NewStore(s.db)
	habitatService := habitat.NewHandler(habitatStore)
	habitatService.RegisterRoutes(subrouter)

	enclosureStore := enclosure.NewStore(s.db)
	enclosureService := enclosure.NewHandler(enclosureStore)
	enclosureService.RegisterRoutes(subrouter)

	animalStore := animal.NewStore(s.db)
	animalService := animal.NewHandler(animalStore)
	animalService.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
