package main

import (
	"github.com/recoilme/sniper"

	. "aoanima.ru/logger"
)

func initSniper() {
	s, _ := sniper.Open(sniper.Dir("sniper"))
	s.Set([]byte("hello"), []byte("go"), 0)
	res, _ := s.Get([]byte("hello"))
	Инфо("  %+v \n", res)
	s.Close()

}

func Сохранить(данныеДляСохранения chan []byte) {
	for данные := range данныеДляСохранения {
		s, _ := sniper.Open(sniper.Dir("sniper"))
		s.Set([]byte("hello"), данные, 0)
		s.Close()
	}

}
