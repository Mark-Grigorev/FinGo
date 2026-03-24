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

### Стек технологий

| Слой           | Технология                              |
|----------------|-----------------------------------------|
| Язык           | Go 1.26                                 |
| HTTP Router    | [Gin](https://github.com/gin-gonic/gin) |
| База данных    | PostgreSQL 17                           |
| SQL            | [sqlc](https://sqlc.dev) + pgx          |
| Миграции       | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Конфиги        | [Viper](https://github.com/spf13/viper) |
| Аутентификация | PASETO токены                           |
| Хэширование    | bcrypt                                  |
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
├── nginx/                # конфиг Nginx
├── docker-compose.yml
└── Dockerfile
```

### Быстрый старт

**Требования:** [Docker](https://www.docker.com/) и Docker Compose

```bash
cp .env.example .env        # настроить окружение
docker compose up --build
```

Приложение будет доступно по адресу `http://localhost`.

**Локальная разработка** (требуется [Go 1.26+](https://go.dev/)):

```bash
make run       # запустить приложение
make migrate   # применить миграции БД
make test      # запустить тесты
make build     # собрать бинарник
```

### Переменные окружения

| Переменная          | По умолчанию   | Описание                 |
|---------------------|----------------|--------------------------|
| `APP_PORT`          | `8008`         | Порт приложения          |
| `DB_HOST`           | `localhost`    | Хост PostgreSQL          |
| `DB_PORT`           | `5432`         | Порт PostgreSQL          |
| `DB_USER`           | `fingo`        | Пользователь БД          |
| `DB_PASSWORD`       | —              | Пароль БД                |
| `DB_NAME`           | `fingo`        | Имя базы данных          |

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

### Tech Stack

| Layer          | Technology                              |
|----------------|-----------------------------------------|
| Language       | Go 1.26                                 |
| HTTP Router    | [Gin](https://github.com/gin-gonic/gin) |
| Database       | PostgreSQL 17                           |
| SQL            | [sqlc](https://sqlc.dev) + pgx          |
| Migrations     | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Config         | [Viper](https://github.com/spf13/viper) |
| Auth           | PASETO tokens                           |
| Password hash  | bcrypt                                  |
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
├── nginx/                # Nginx config
├── docker-compose.yml
└── Dockerfile
```

### Getting Started

**Prerequisites:** [Docker](https://www.docker.com/) & Docker Compose

```bash
cp .env.example .env        # configure environment
docker compose up --build
```

App will be available at `http://localhost`.

**Local development** (requires [Go 1.26+](https://go.dev/)):

```bash
make run       # run the app
make migrate   # apply DB migrations
make test      # run tests
make build     # build binary
```

### Environment Variables

| Variable            | Default        | Description              |
|---------------------|----------------|--------------------------|
| `APP_PORT`          | `8008`         | Application port         |
| `DB_HOST`           | `localhost`    | PostgreSQL host          |
| `DB_PORT`           | `5432`         | PostgreSQL port          |
| `DB_USER`           | `fingo`        | Database user            |
| `DB_PASSWORD`       | —              | Database password        |
| `DB_NAME`           | `fingo`        | Database name            |

### License

MIT
