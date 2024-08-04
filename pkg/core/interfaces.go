package core

type DBStorage interface {
	PingDB() string
}
