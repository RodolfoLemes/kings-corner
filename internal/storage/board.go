package storage

import (
	"errors"

	"kings-corner/internal/game"

	"github.com/rs/xid"
)

type Board struct {
	ID string
	game.Board
}

type BoardRepository interface {
	Create(b game.Board) (*Board, error)
	GetByID(boardID string) (*Board, error)
}

func (in *inMemoryBoardRepository) Create(b game.Board) (*Board, error) {
	id := xid.New().String()

	board := &Board{id, b}

	in.boards[id] = board

	return board, nil
}

func (in *inMemoryBoardRepository) GetByID(boardID string) (*Board, error) {
	board, found := in.boards[boardID]

	if !found {
		return nil, errors.New("board not found")
	}

	return board, nil
}
