# Для запуска приложения используем команду:
- Добавил путь к конфигу так как не удобно будет добавлять путь в окружение.(лишние команды в консоли)
## ```go run ./cmd/webserver --config configs/local.yaml```

### реализация загрузки конфига 
```func MustLoadConfig(configPath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Printf("Failed to read config file: %v", err)
		return nil, err
	}
	return &cfg, nil
}```




