# FinGo 💰

**[🇷🇺 Русский](#fingo---трекер-личных-финансов) · [🇬🇧 English](#fingo---personal-finance-tracker)**

---

## FinGo 💰 — Трекер личных финансов

> Отслеживай расходы, устанавливай бюджеты и анализируй свои финансы.

### О проекте

FinGo — self-hosted веб-приложение для управления личными финансами. Позволяет вести несколько счетов, категоризировать транзакции, устанавливать месячные бюджеты, отслеживать регулярные платежи и визуализировать финансовую активность через отчёты и графики.

### Функциональность

- **Счета** — управление наличными, картами и накопительными счетами
- **Транзакции** — учёт доходов и расходов с категориями, фильтрами и фото чеков
- **Бюджеты** — месячные лимиты по категориям с прогресс-баром и уведомлениями
- **Регулярные платежи** — автоматизация подписок и периодических платежей
- **Аналитика** — динамика доходов/расходов, топ категорий, сравнение периодов
- **Экспорт** — CSV, PDF, Excel
- **Восстановление пароля** — сброс пароля по ссылке на email (через Resend)

### Стек технологий

| Слой           | Технология                              |
|----------------|-----------------------------------------|
| Язык           | Go 1.25                                 |
| HTTP Router    | [Gin](https://github.com/gin-gonic/gin) |
| База данных    | PostgreSQL 17                           |
| SQL            | [sqlc](https://sqlc.dev) + pgx          |
| Миграции       | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Конфиги        | [Viper](https://github.com/spf13/viper) |
| Аутентификация | PASETO токены                           |
| Хэширование    | bcrypt                                  |
| Email          | [Resend](https://resend.com)            |
| Фронтенд       | [Pencil](https://pencil.dev)            |
| Reverse proxy  | Nginx                                   |
| Контейнеры     | Docker / Docker Compose                 |

### Структура проекта

```
.
├── cmd/                  # точка входа
├── internal/
│   ├── domain/           # сущности и бизнес-логика
│   ├── repository/       # работа с БД (sqlc)
│   ├── service/          # бизнес-слой
│   └── handler/          # HTTP-хэндлеры (Gin)
├── pkg/                  # переиспользуемые утилиты
├── migrations/           # миграции БД
├── docs/                 # Swagger / OpenAPI
├── nginx/
│   ├── nginx.conf        # конфиг Nginx для Kubernetes
│   └── nginx-dev.conf    # конфиг Nginx для локального запуска
├── docker-compose.yml
├── Dockerfile.fingo
└── Dockerfile.nginx
```

### Быстрый старт

**Требования:** [Docker](https://www.docker.com/) и Docker Compose

```bash
cp .env.example .env              # настроить окружение
docker compose up --build -d
```

Приложение будет доступно по адресу `http://localhost:8001`.

**Локальная разработка** (требуется [Go 1.25+](https://go.dev/)):

```bash
make run       # запустить приложение
make migrate   # применить миграции БД
make test      # запустить тесты
make build     # собрать бинарник
```

### Переменные окружения

| Переменная             | По умолчанию            | Описание                              |
|------------------------|-------------------------|---------------------------------------|
| `APP_PORT`             | `8008`                  | Порт приложения                       |
| `APP_ENV`              | `local`                 | Окружение (`local`, `prod`)           |
| `DB_CONN_STRING`       | —                       | Строка подключения к PostgreSQL       |
| `TOKEN_SYMMETRIC_KEY`  | —                       | PASETO ключ (64 hex-символа)          |
| `TOKEN_DURATION`       | `24h`                   | Время жизни токена                    |
| `RESEND_API_KEY`       | —                       | API ключ [Resend](https://resend.com) |
| `APP_BASE_URL`         | `http://localhost:8008` | Базовый URL для ссылок в письмах      |

### Коды завершения

При аварийном завершении процесс возвращает уникальный код — это позволяет определить причину сбоя без просмотра кода.

| Код | Константа          | Причина                                                       |
|-----|--------------------|---------------------------------------------------------------|
| `0` | `exitOK`           | Успешное завершение (штатный shutdown по сигналу)             |
| `2` | `exitConfigError`  | Ошибка загрузки конфигурации (например, не задан `DB_CONN_STRING`) |
| `3` | `exitTokenError`   | Ошибка инициализации PASETO-токенера (неверный `TOKEN_SYMMETRIC_KEY`) |
| `4` | `exitDBConnect`    | Не удалось подключиться к базе данных (недоступен PostgreSQL) |
| `5` | `exitDBMigrate`    | Ошибка выполнения миграций БД                                 |
| `6` | `exitServerShutdown` | Graceful shutdown не завершился в отведённое время          |

> Код `1` зарезервирован системой (panic в `main`). Приложение его не использует.

### Лицензия

MIT

---

## FinGo 💰 — Personal Finance Tracker

> Track expenses, set budgets, and analyze your spending habits.

### About

FinGo is a self-hosted web application for personal finance management. It lets you manage multiple accounts, categorize transactions, set monthly budgets, track recurring payments, and visualize your financial activity through reports and charts.

### Features

- **Accounts** — manage cash, cards, and savings accounts
- **Transactions** — income/expense tracking with categories, filters, and receipt photos
- **Budgets** — set monthly limits per category with progress tracking and alerts
- **Recurring payments** — automate subscriptions and regular bills
- **Analytics** — income/expense dynamics, top spending categories, period-over-period comparison
- **Export** — CSV, PDF, Excel
- **Password recovery** — reset password via email link (powered by Resend)

### Tech Stack

| Layer          | Technology                              |
|----------------|-----------------------------------------|
| Language       | Go 1.25                                 |
| HTTP Router    | [Gin](https://github.com/gin-gonic/gin) |
| Database       | PostgreSQL 17                           |
| SQL            | [sqlc](https://sqlc.dev) + pgx          |
| Migrations     | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Config         | [Viper](https://github.com/spf13/viper) |
| Auth           | PASETO tokens                           |
| Password hash  | bcrypt                                  |
| Email          | [Resend](https://resend.com)            |
| Frontend       | [Pencil](https://pencil.dev)            |
| Reverse proxy  | Nginx                                   |
| Containerized  | Docker / Docker Compose                 |

### Project Structure

```
.
├── cmd/                  # entry point
├── internal/
│   ├── domain/           # entities and business logic
│   ├── repository/       # DB layer (sqlc)
│   ├── service/          # business layer
│   └── handler/          # HTTP handlers (Gin)
├── pkg/                  # shared utilities
├── migrations/           # DB migrations
├── docs/                 # Swagger / OpenAPI
├── nginx/
│   ├── nginx.conf        # Nginx config for Kubernetes
│   └── nginx-dev.conf    # Nginx config for local Docker Compose
├── docker-compose.yml
├── Dockerfile.fingo
└── Dockerfile.nginx
```

### Getting Started

**Prerequisites:** [Docker](https://www.docker.com/) & Docker Compose

```bash
cp .env.example .env              # configure environment
docker compose up --build -d
```

App will be available at `http://localhost:8001`.

**Local development** (requires [Go 1.25+](https://go.dev/)):

```bash
make run       # run the app
make migrate   # apply DB migrations
make test      # run tests
make build     # build binary
```

### Environment Variables

| Variable               | Default                 | Description                               |
|------------------------|-------------------------|-------------------------------------------|
| `APP_PORT`             | `8008`                  | Application port                          |
| `APP_ENV`              | `local`                 | Environment (`local`, `prod`)             |
| `DB_CONN_STRING`       | —                       | PostgreSQL connection string              |
| `TOKEN_SYMMETRIC_KEY`  | —                       | PASETO key (64 hex chars)                 |
| `TOKEN_DURATION`       | `24h`                   | Token TTL                                 |
| `RESEND_API_KEY`       | —                       | API key from [Resend](https://resend.com) |
| `APP_BASE_URL`         | `http://localhost:8008` | Base URL used in password reset emails    |

### Exit Codes

Each failure scenario returns a unique exit code so you can identify the root cause without reading the source.

| Code | Constant           | Reason                                                        |
|------|--------------------|---------------------------------------------------------------|
| `0`  | `exitOK`           | Clean exit (graceful shutdown on signal)                      |
| `2`  | `exitConfigError`  | Configuration load failed (e.g. `DB_CONN_STRING` not set)     |
| `3`  | `exitTokenError`   | PASETO token maker init failed (invalid `TOKEN_SYMMETRIC_KEY`) |
| `4`  | `exitDBConnect`    | Could not connect to the database (PostgreSQL unreachable)    |
| `5`  | `exitDBMigrate`    | Database migration failed                                     |
| `6`  | `exitServerShutdown` | Graceful shutdown did not complete in time                  |

> Code `1` is reserved by the Go runtime (unrecovered panic). The application never returns it explicitly.

### License

MIT
