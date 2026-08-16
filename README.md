# Auth Service

Микросервис аутентификации и авторизации пользователей на основе JWT с использованием RSA-ключей. Сервис предоставляет функциональность регистрации, входа и обновления токенов доступа.

## Описание

Auth Service — это REST API для управления аутентификацией пользователей. Сервис использует:
- **JWT (JSON Web Tokens)** с RSA-криптографией для генерации и валидации токенов доступа
- **PostgreSQL** для хранения данных пользователей
- **Swagger/OpenAPI** для документирования API
- **Docker** для контейнеризации и развертывания

### Основные возможности:
- ✅ Регистрация новых пользователей
- ✅ Аутентификация с использованием username/password
- ✅ Выдача JWT access и refresh токенов
- ✅ Обновление access токена через refresh токен
- ✅ Безопасное хранение паролей (поддержка bcrypt и argon2)
- ✅ Автоматические миграции базы данных
- ✅ API документация через Swagger UI

## 🛠 Технологический стек

| Компонент | Версия | Описание |
|-----------|--------|---------|
| **Go** | 1.26.4 | Язык программирования |
| **PostgreSQL** | - | База данных |
| **JWT** | v5.3.1 | Аутентификация токенов |
| **pgx** | v5.10.0 | PostgreSQL драйвер |
| **Swagger** | v1.8.1 | Документация API |

## 🔧 Переменные окружения (.env)

Создайте файл `.env` в корне проекта с следующими переменными:

```env
# Database Configuration
DB_USER=postgres
DB_PASS=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=auth_service

# JWT RSA Keys
PRIVATE_KEY_PATH=./secrets/private.key
PUBLIC_KEY_PATH=./secrets/public.key
```

### Описание переменных:

| Переменная | Тип | Обязательная | Описание |
|-----------|-----|------------|---------|
| `DB_USER` | string | ✅ | Имя пользователя PostgreSQL |
| `DB_PASS` | string | ✅ | Пароль пользователя PostgreSQL |
| `DB_HOST` | string | ✅ | Хост сервера PostgreSQL (localhost, IP или доменное имя) |
| `DB_PORT` | int | ✅ | Порт PostgreSQL (по умолчанию 5432) |
| `DB_NAME` | string | ✅ | Название базы данных |
| `PRIVATE_KEY_PATH` | string | ✅ | Путь к приватному RSA ключу (PEM формат) |
| `PUBLIC_KEY_PATH` | string | ✅ | Путь к публичному RSA ключу (PEM формат) |

**Примечание:** Если файл `.env` не найден, сервис будет использовать системные переменные окружения.

## 🗝 Генерация RSA ключей

Для работы сервиса необходимо сгенерировать пару RSA ключей:

```bash
# Генерация приватного ключа (4096 бит)
openssl genrsa -out secrets/private.key 4096

# Извлечение публичного ключа из приватного
openssl rsa -in secrets/private.key -pubout -out secrets/public.key
```

## 🚀 Запуск

### Локальный запуск

```bash
# Установка зависимостей
go mod download

# Запуск сервиса
go run cmd/main.go
```

### Запуск с Docker

```bash
# Сборка образа
docker build -t auth-service:latest .

# Запуск контейнера
docker-compose up -d
```

## 📚 API Документация

После запуска сервиса, документация Swagger доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

### Основные эндпоинты:

- `POST /auth/register` — Регистрация пользователя
- `POST /auth/login` — Вход в систему
- `POST /auth/refresh` — Обновление токена доступа

## 🗄 Структура проекта

```
.
├── cmd/                          # Entry point
│   └── main.go
├── internal/                      # Внутренний код приложения
│   ├── app/                       # Инициализация приложения
│   │   ├── app.go
│   │   ├── config.go
│   │   └── migrations.go
│   ├── auth/                      # Логика аутентификации
│   │   ├── jwt/                   # JWT обработка
│   │   │   ├── config.go
│   │   │   ├── decoder.go
│   │   │   ├── encoder.go
│   │   │   ├── generator.go
│   │   │   └── errors.go
│   │   └── password/              # Хеширование паролей
│   │       ├── argon.go
│   │       ├── bcrypt.go
│   │       ├── argon_test.go
│   │       └── bcrypt_test.go
│   ├── database/                  # Работа с БД
│   │   └── postgres.go
│   ├── domain/                    # Доменные модели
│   │   ├── auth.go
│   │   └── user.go
│   ├── dto/                       # Data Transfer Objects
│   │   └── user.go
│   ├── handlers/                  # HTTP обработчики
│   │   ├── auth.go
│   │   └── errors.go
│   ├── repositories/              # Работа с БД
│   │   └── user.go
│   └── usecases/                  # Бизнес-логика
│       ├── auth.go
│       ├── login.go
│       ├── refresh.go
│       ├── register.go
│       ├── user.go
│       └── errors.go
├── migrations/                    # SQL миграции
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── secrets/                       # RSA ключи
│   ├── private.key
│   └── public.key
├── docs/                          # Swagger документация
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── Dockerfile                     # Docker конфигурация
├── docker-compose.yaml            # Docker Compose конфигурация
├── go.mod                         # Go модули
└── README.md                      # Этот файл
```

## 🔐 Безопасность

- Пароли хешируются с использованием bcrypt или argon2
- JWT токены подписываются приватным RSA ключом
- Refresh токены хранятся в виде secure-only cookies
- Публичный ключ используется для валидации токенов в других сервисах
