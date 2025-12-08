# Для запуска приложения используем команду:
```bash
go run ./cmd/webserver --config configs/local.yaml
```
- Добавил путь к конфигу так как не удобно будет добавлять путь в окружение.(лишние команды в консоли)


## Реализация загрузки конфига 
#### Считывания конфига с помощью библиотеки cleanenv
```go func MustLoadConfig(configPath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Printf("Failed to read config file: %v", err)
		return nil, err
	}
	return &cfg, nil
}
```
#### Реализация считывания флага при запуске програмы 

```go func pafseFlags() string {
	var configPath string
	pflag.StringVar(&configPath, "config", "configs/local.yaml", "Path to config file")
	pflag.Parse()
	return configPath
}
```

## Реализация создание базы данных
#### Использовалась sqlite так как не надо использовать дополнительных приложений, для создании бд, реализована миграция которая с использование пакета "modernc.org/sqlite". Написан файл инициирования базы данных, считывается путь из конфига, который мы указали в параметрах при запуске

```go
func NewStorage(pathDb string) error {

	db, err := sql.Open("sqlite", pathDb)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	sqlComand, err := os.ReadFile("./internal/storage/migrations/001_init.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
		return err
	}

	_, err = db.Exec(string(sqlComand))
	if err != nil {
		log.Fatalf("Failed to execute SQL commands: %v", err)
		return err
	}

	return nil
}
```

## Реализация handlers

#### Что реализовано
- POST /api/links
#### Принимает список ссылок, создаёт batch, записывает их в таблицу links со статусом queued, проверяет каждую ссылку “на лету” и возвращает готовый результат в этом же ответе.
- GET /api/report?links_list=1,2,3
#### Принимает набор ID, собирает ссылки, формирует PDF и отдаёт файл пользователю.

#### Особенности
- Проверка ссылок выполняется сразу → POST возвращает итоговые статусы
- Все операции пишутся в SQLite
- Реализована обработка ошибок и логирование
- Генерация PDF — через gofpdf

## Реализация router
#### Используется chi:
- Подключена middleware-логика
- Сгруппированы маршруты /api/*
- Вынесены health-check и API-методы отдельно

## Реализация запуска сервера
#### Используется http.Server
- Реализовано корректное завершение работы:
- ловим os.Interrupt
- создаём контекст с таймаутом 15 секунд
- вызываем server.Shutdown(ctx)
#### Это позволяет корректно завершить все активные соединения.

## Обратная свзязь по заданию.
#### Я бы изменил логику приложения, чтобы у нас был запрос POST на запись ссылок которые необходимо обработать, и метод GET который отдает обработанные ссылки, более масштабируемо и безопасно.

## Тесты в postman

<img width="933" height="619" alt="image" src="https://github.com/user-attachments/assets/f2210788-2ff4-4169-8d74-055e3319b9af" />

<img width="936" height="696" alt="image" src="https://github.com/user-attachments/assets/90fa13b2-6f29-4526-8ac5-2f03bbe0d4fb" />


# ДАТА ЗАВЕРЕШНИЯ: 08.12.2025
