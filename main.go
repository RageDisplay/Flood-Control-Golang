package main

import (
	"context"
	floodcontrol "task/func"
	"time"
)

func main() {
	flooduse := floodcontrol.ExFloodControl(floodcontrol.Config{}.DefaultConfig()) // Создаем объект Config c нулевыми значениями, применяем DefaultConfig для установки конфигурации и передаем его в функцию ExFloodControl

	for i := 0; i < 6; i++ { //Создаем 6 запросов для проверки
		allow, err := flooduse.Check(context.Background(), int64(i)) // Проверка запроса
		if err != nil {
			panic(err)
		}
		if allow { // Если true, то запрос разрешен
			println("OK")
		} else { // Если false, то запрос заблокирован
			println("BLOCKED")
		}
	}

	time.Sleep(6 * time.Second) // Ждём 6 секунд, чтобы все запросы попали в период проверки

	for i := 0; i < 2; i++ { //Создаем 2 запроса для проверки
		allow, err := flooduse.Check(context.Background(), int64(i)) // Проверка запроса
		if err != nil {
			panic(err)
		}
		if allow { // Если true, то запрос разрешен
			println("OK")
		} else { // Если false, то запрос заблокирован
			println("BLOCKED")
		}
	}
}
