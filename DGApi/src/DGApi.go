package main

import (
	. "aoanima.ru/DGApi"
	. "aoanima.ru/Logger"
	. "aoanima.ru/QErrors"
)

func main() {
	граф, закрыть := ДГраф()
	defer закрыть()

	for i := 0; i < 1000; i++ {
		go func() {
			Даннные := ДанныеЗапроса{
				Запрос: "set ",
				Данные: make(map[string]string),
				}
			результат, статус := Изменить(Даннные, граф)
			if статус.Код != Ок {
				Инфо(" %+v \n", статус.Текст)
			}
			Инфо(" результат %+v \n", результат)
		}()
	}
	

}
