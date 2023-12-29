package ConnQuic

import (
	"testing"

	// . "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

func main() {
	TestВставить(nil)
}

func TestВставить(t *testing.T) {
	бд := База{}

	бд.Подключиться("./test.db")

	объект := map[string]interface{}{
		"Логин": "логин_клиента",
		"Имя":   "Саня",
		"Адрес": map[string]interface{}{
			"Страна":   "Россия",
			"Город":    "Москва",
			"Улица":    "Льва Толстого",
			"Дом":      "16",
			"Квартира": "2",
		},
	}

	ошибка := бд.Вставить("test", "key_test", объект, []string{"Адрес.Город"})
	if ошибка.Код != 0 {
		Ошибка("  %+v \n", ошибка)
	}

}
