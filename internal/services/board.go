package services

import (
	"errors"

	"kings-corner/internal/deck"
	"kings-corner/internal/game"
	"kings-corner/internal/storage"
)

type BoardService struct {
	boardRepository  storage.BoardRepository
	playerRepository storage.PlayerRepository
}

func NewBoardService(
	boardRepository storage.BoardRepository,
	playerRepository storage.PlayerRepository,
) *BoardService {
	return &BoardService{boardRepository, playerRepository}
}

func (bs *BoardService) Create() (*game.Board, error) {
	player := game.NewPlayer()

	d := deck.New()
	board := game.New(d)
	board.Join(player)

	sBoard, err := bs.boardRepository.Create(*board)
	if err != nil {
		return nil, err
	}

	_, err = bs.playerRepository.Create(player, sBoard.ID)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (bs *BoardService) Join(boardID string) (*game.Board, *game.Player, error) {
	board, err := bs.boardRepository.GetByID(boardID)
	if err != nil {
		return nil, nil, err
	}

	player := game.NewPlayer()
	board.Join(player)

	err = bs.boardRepository.Update(board.Board)
	if err != nil {
		return nil, nil, err
	}

	_, err = bs.playerRepository.Create(player, boardID)
	if err != nil {
		return nil, nil, err
	}

	return &board.Board, &player, nil
}

func (bs *BoardService) Run(boardID string) error {
	board, err := bs.boardRepository.GetByID(boardID)
	if err != nil {
		return err
	}

	return board.Run()
}

func (bs *BoardService) Play(boardID string, playerID string, turn game.Turn) error {
	board, err := bs.boardRepository.GetByID(boardID)
	if err != nil {
		return err
	}

	sPlayer, err := bs.playerRepository.GetByID(playerID)
	if err != nil {
		return err
	}

	if sPlayer.BoardID != board.ID {
		return errors.New("forbidden")
	}

	var player game.Player
	for i := range board.Players {
		if board.Players[i].ID() == sPlayer.ID {
			player = board.Players[i]
			break
		}
	}

	err = player.Play(turn)

	return err
}
