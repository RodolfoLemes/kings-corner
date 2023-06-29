package storage

import (
	"errors"

	"kings-corner/internal/game"
)

type Board struct {
	game.Board
}

type BoardRepository interface {
	Create(b game.Board) (*Board, error)
	Update(b game.Board) error
	GetByID(boardID string) (*Board, error)
}

func (in *inMemoryBoardRepository) Create(b game.Board) (*Board, error) {
	board := &Board{b}

	in.boards[b.ID] = board

	return board, nil
}

func (in *inMemoryBoardRepository) Update(b game.Board) error {
	_, found := in.boards[b.ID]

	if !found {
		return errors.New("board not found")
	}

	in.boards[b.ID] = &Board{b}

	return nil
}

func (in *inMemoryBoardRepository) GetByID(boardID string) (*Board, error) {
	board, found := in.boards[boardID]

	if !found {
		return nil, errors.New("board not found")
	}

	return board, nil
}
