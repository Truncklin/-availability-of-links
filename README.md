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




## Обратная свзязь по заданию.
#### Я бы изменил логику приложения, чтобы у нас был запрос POST на запись ссылок которые необходимо обработать, и метод GET который отдает обработанные ссылки, более масштабируемо и безопасно.