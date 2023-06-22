package storage

import "log"

type inMemory struct {
	*inMemoryBoardRepository
	*inMemoryPlayerRepository
}

type inMemoryBoardRepository struct {
	boards map[string]*Board
}

type inMemoryPlayerRepository struct {
	players map[string]*Player
}

func newInMemory() *inMemory {
	return &inMemory{
		inMemoryBoardRepository:  &inMemoryBoardRepository{},
		inMemoryPlayerRepository: &inMemoryPlayerRepository{},
	}
}

func (in *inMemory) Run() {
	in.boards = make(map[string]*Board)
	in.players = make(map[string]*Player)

	log.Println("in-memory database running...")
}

func (in *inMemory) BoardRepository() BoardRepository {
	return in.inMemoryBoardRepository
}

func (in *inMemory) PlayerRepository() PlayerRepository {
	return in.inMemoryPlayerRepository
}
