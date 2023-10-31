package ConnQuic

import (
	"errors"

	quic "github.com/quic-go/quic-go"
)

type ОчередьПотоков struct {
	Потоки chan quic.Stream
}

func НоваяОчередьПотоков(размер int) *ОчередьПотоков {
	return &ОчередьПотоков{
		Потоки: make(chan quic.Stream, размер),
	}
}
func (о *ОчередьПотоков) Взять(поток quic.Stream) (quic.Stream, error) {
	select {
	case поток := <-о.Потоки:
		return поток, nil
	default:
		return nil, errors.New("Нет свободных потоков")
	}

}

func (о *ОчередьПотоков) Вернуть(поток quic.Stream) {
	select {
	case о.Потоки <- поток:
	default:
		// Если канал полон, просто закрываем поток
		поток.Close()
	}
}
