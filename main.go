package main

import (
	"context"
	floodcontrol "task/func"
	"time"
)

func main() {
	flooduse := floodcontrol.ExFloodControl(5, 3) // Проверка за последние 5 секунд, максимум 3 запроса

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
