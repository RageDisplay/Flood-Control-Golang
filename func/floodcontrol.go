package floodcontrol

import (
	"context"
	"time"
)

type StatStorage struct { // Общее хранилище данных о запросах пользователей
	req map[int64]time.Time
}

type FloodControlStruct struct { // Структура, реализующая интерфейс FloodControl
	storage  *StatStorage
	N        int
	K        int
	increase chan int64
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

func ExFloodControl(N, K int) FloodControl { // Функция, которая реализует FloodControl
	flooduse := &FloodControlStruct{
		storage:  &StatStorage{req: make(map[int64]time.Time)}, // Создаем объект StatStorage
		N:        N,
		K:        K,
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
		if sec.Sub(lastReq) > time.Duration(flooduse.N)*time.Second { // Если запрос не прошел проверку, то удаляем его из очереди
			delete(flooduse.storage.req, user)
		}
	}
}

func (flooduse FloodControlStruct) Check(ctx context.Context, userID int64) (bool, error) { // Проверка запроса
	sec := time.Now()
	flooduse.clean(sec) // Удаляем запросы, до истечения времени проверки

	count := len(flooduse.storage.req) // Получаем количество запросов в очереди
	if count >= flooduse.K {           // Если количество запросов в очереди больше максимального, то возвращаем false
		return false, nil
	}
	flooduse.increase <- userID // Добавляем пользователя в очередь
	return true, nil            // Возвращаем true
}
