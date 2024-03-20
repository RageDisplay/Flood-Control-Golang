package floodcontrol

type Config struct {
	N int
	K int
}

func (c Config) DefaultConfig() *Config { // Функция, которая возвращает конфигурацию по умолчанию
	return &Config{
		N: 5, // Время ожидания в main 6 секунд
		K: 3,
	}
}
