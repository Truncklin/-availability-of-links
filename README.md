# Для запуска приложения используем команду:
```bash
go run ./cmd/webserver --config configs/local.yaml
```
- Добавил путь к конфигу так как не удобно будет добавлять путь в окружение.(лишние команды в консоли)


### Реализация загрузки конфига 
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





