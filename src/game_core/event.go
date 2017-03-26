package game_core

type IGameEvent interface {
	OnGameEvent(buff []byte)
}
