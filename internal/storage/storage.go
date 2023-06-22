package storage

type Storage interface {
	Run()

	BoardRepository() BoardRepository
	PlayerRepository() PlayerRepository
}

func New() Storage {
	return newInMemory()
}
