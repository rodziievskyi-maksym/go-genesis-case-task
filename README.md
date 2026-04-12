# go-genesis-case-task

GitHub Release Notifier: A Go-based service that monitors repositories and sends email alerts for new releases.

🚀 GitHub Release Notifier
GitHub Release Notifier — це високонавантажений мікросервіс на Go, розроблений для моніторингу оновлень у GitHub репозиторіях та автоматичного сповіщення користувачів. Проект побудований за принципами Clean Architecture з фокусом на продуктивність та надійність.

✨ Key Features
Real-time Monitoring: Автоматичне сканування GitHub репозиторіїв на наявність нових тегів.

Smart Caching: Використання Redis (Decorator Pattern) для мінімізації запитів до GitHub API.

Rate Limit Protection: Проактивна обробка лімітів GitHub (сервіс "засинає" до моменту скидання ліміту).

Clean Architecture: Чіткий поділ на шари (Domain, UseCase, Infrastructure, Delivery).

Твій поточний README — це хороший фундамент, але щоб він виглядав як проект Senior рівня, нам треба додати трохи "структурного лиску": чіткі інструкції з тестування, Swagger та деталі про архітектурні рішення (наприклад, твій Redis Decorator).

🚀 GitHub Release Notifier
GitHub Release Notifier — це мікросервіс на Go, розроблений для моніторингу оновлень у GitHub репозиторіях та автоматичного сповіщення користувачів. Проект побудований за принципами Clean Architecture з фокусом на продуктивність та надійність.

✨ Key Features
Real-time Monitoring: Автоматичне сканування GitHub репозиторіїв на наявність нових тегів.

Smart Caching: Використання Redis (Decorator Pattern) для мінімізації запитів до GitHub API.

Rate Limit Protection: Проактивна обробка лімітів GitHub (сервіс "засинає" до моменту скидання ліміту).

Clean Architecture: Чіткий поділ на шари (Domain, UseCase, Infrastructure, Delivery).

Full Observability: Метрики Prometheus та структуроване логування через slog.

🛠 Technical Stack
Language: Go 1.26

Framework: Gin Gonic

Database: PostgreSQL 17 (pgx pool)

Caching: Redis 7

Auth: API Key Middleware

Monitoring: Prometheus

Docs: Swagger (OpenAPI)

🚀 Quick Start
1. Environment Configuration
Створіть .env файл у корені проекту:

Code snippet
SERVER_PORT=8080
DB_DSN=postgres://dev:pass@db:5432/genesis_task?sslmode=disable
GITHUB_TOKEN=your_token_here
SCANNER_INTERVAL=10m
API_KEY=your_secret_api_key

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_CACHE_TTL=10m

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_app_password

2. Execution via Docker
Запустіть весь стек однією командою:

docker-compose up --build

📖 API Documentation
Після запуску документація Swagger доступна за адресою:
👉 http://localhost:8080/swagger/index.html

Note: Всі запити (крім Swagger) потребують заголовок X-API-KEY: genesis-secret-key

🧪 Testing
Проект має високе покриття Unit-тестами (Table-driven tests) для бізнес-логіки.

Запуск усіх тестів:

go test ./... -v

go test -coverprofile=cover.out ./internal/usecase/...
go tool cover -func=cover.out

📂 Project Structure
cmd/ — Точка входу (ініціалізація DI контейтера).

internal/domain/ — Бізнес-сутності та інтерфейси.

internal/usecase/ — Реалізація бізнес-логіки (табл-тести тут).

internal/infrastructure/ — Зовнішні сервіси (GitHub API з Redis Decorator, DB, Email).

internal/worker/ — Scanner (фонoві задачі).

migrations/ — SQL файли для ініціалізації БД (Postgres /docker-entrypoint-initdb.d).