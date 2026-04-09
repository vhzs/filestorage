# FileStorage

Сервис загрузки и хранения файлов.

## Стек

- Go 1.22+
- PostgreSQL (метаданные файлов, пользователи)
- Redis (кэш метаданных)
- Chi
- pgx
- golang-jwt (авторизация)
- goose (миграции)
- slog (логирование)
- cleanenv (конфиг)
- swaggo (Swagger-документация)
- Docker Compose

## Функциональность

- Регистрация и JWT-авторизация
- Загрузка файлов (POST /api/files)
- Скачивание файла по ID (GET /api/files/{id}/download)
- Список файлов пользователя (GET /api/files)
- Удаление файла (DELETE /api/files/{id})
- Квоты на объём хранилища для каждого пользователя
- Поиск по имени файла (GET /api/files?search=...)

## Структура

```
filestorage/
├── cmd/filestorage/     — точка входа
├── internal/
│   ├── config/          — загрузка конфига
│   ├── handler/         — HTTP-хендлеры
│   ├── service/         — бизнес-логика
│   ├── storage/         — работа с PostgreSQL и Redis
│   ├── auth/            — JWT, middleware авторизации
│   └── model/           — структуры данных
├── uploads/             — директория для хранения файлов
├── migrations/          — SQL-миграции (goose)
├── docs/                — сгенерированный Swagger
├── docker-compose.yml
├── Dockerfile
└── .env.example
```

## Запуск

```bash
docker-compose up --build
```
