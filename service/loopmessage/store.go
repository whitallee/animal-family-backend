package loopmessage

import (
	"database/sql"
	"fmt"

	"github.com/whitallee/animal-family-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ReceiveLoopMessage(loopMessage types.InboundLoopMessagePayload) error {
	fmt.Println(loopMessage.Text)
	return nil
}

func (s *Store) SendLoopMessage(loopMessage types.SentLoopMessagePayload) error {
	fmt.Println(loopMessage.Text)
	return nil
}
