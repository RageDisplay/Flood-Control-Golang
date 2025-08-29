package floodcontrol

import (
	"context"
	"time"
)

type StatStorage struct { // Общее хранилище данных о запросах пользователей
	req map[int64]time.Time
}

type FloodControlStruct struct { 
	storage  *StatStorage // Создаем объект StatStorage
	config   *Config      // Конфигурация
	increase chan int64   // Канал для увеличения количества запросов в очередь
}

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

func ExFloodControl(config *Config) FloodControl { // Функция, которая реализует FloodControl
	flooduse := &FloodControlStruct{
		storage:  &StatStorage{req: make(map[int64]time.Time)}, // Создаем объект StatStorage
		config:   config,
		increase: make(chan int64), // Создаем канал для увеличения количества запросов в очередь
	}
	go flooduse.increaser() // Запускаем процесс, который увеличивает количество запросов в очередь
	return flooduse
}

func (flooduse FloodControlStruct) increaser() { // Процесс, который увеличивает количество запросов в очередь
	for userID := range flooduse.increase {
		flooduse.storage.req[userID] = time.Now() // Добавляем пользователя в очередь
	}
}

func (flooduse FloodControlStruct) clean(sec time.Time) { // Удаляет запросы, до истечения времени проверки
	for user, lastReq := range flooduse.storage.req {
		if sec.Sub(lastReq) > time.Duration(flooduse.config.N)*time.Second { // Если запрос не прошел проверку, то удаляем его из очереди
			delete(flooduse.storage.req, user)
		}
	}
}

func (flooduse FloodControlStruct) Check(ctx context.Context, userID int64) (bool, error) { // Проверка запроса
	sec := time.Now()
	flooduse.clean(sec) // Удаляем запросы, до истечения времени проверки

	count := len(flooduse.storage.req) // Получаем количество запросов в очереди
	if count >= flooduse.config.K {    // Если количество запросов в очереди больше максимального, то возвращаем false
		return false, nil
	}
	flooduse.increase <- userID // Добавляем пользователя в очередь
	return true, nil            // Возвращаем true
}

