# Financial Tracker API 📊

![GitHub repo size](https://img.shields.io/github/repo-size/nemopss/financial-tracker?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nemopss/financial-tracker)
![GitHub contributors](https://img.shields.io/github/contributors/nemopss/financial-tracker?style=flat-square)
![GitHub issues](https://img.shields.io/github/issues/nemopss/financial-tracker?style=flat-square)

## 📜 Описание

**Financial Tracker API** — это RESTful API для отслеживания личных финансов, предоставляющее возможности для управления категориями, транзакциями и получения аналитики. Проект построен на **Go (Gin)** и документирован с использованием **Swagger**.

---

## 🚀 Возможности API

- Регистрация пользователей и авторизация через JWT
- CRUD-операции для категорий и транзакций
- Аналитика доходов и расходов
- Swagger UI для удобной документации и тестирования

---

## 📂 Структура проекта

```plaintext
├── cmd/
│   └── main.go                  # Точка входа в приложение
├── config/
│   └── config.go                # Загрузка и управление конфигурацией
├── docs/
│   ├── docs.go                  # Автогенерация Swagger-документации
│   ├── swagger.json             # Swagger-спецификация (JSON)
│   └── swagger.yaml             # Swagger-спецификация (YAML)
├── internal/
│   ├── handlers/                # Обработчики HTTP-запросов
│   │   ├── analytics.go         # Аналитика
│   │   ├── auth.go              # Аутентификация
│   │   ├── category.go          # Категории
│   │   └── transaction.go       # Транзакции
│   ├── middleware/              # Middleware
│   │   └── auth.go              # JWT-проверка авторизации
│   ├── repository/              # Логика работы с БД
│   │   ├── analytics.go         # SQL для аналитики
│   │   ├── category.go          # SQL для категорий
│   │   ├── transactions.go      # SQL для транзакций
│   │   ├── user.go              # SQL для пользователей
│   │   ├── db.go                # Подключение к БД
│   │   └── repository.go        # Интерфейс репозитория
│   ├── response/                # Унификация ответ от сервера
│   │   └── response.go          # Success и Error ответы
├── migrations/                  # SQL-скрипты для миграции базы данных
│   ├── 20241122150001_create_users.sql
│   ├── 20241122150002_create_categories.sql
│   └── 20241122150003_create_transactions.sql
├── .env                         # Конфигурация среды
├── go.mod                       # Зависимости Go
└── README.md                    # Документация проекта
```

---

## 🛠 Технологии

- **Go**: язык программирования
- **Gin**: веб-фреймворк
- **Swagger**: генерация и документация API
- **PostgreSQL**: база данных
- **JWT**: аутентификация
- **SQLBoiler** (возможно, в будущем): ORM

---

## 📖 Документация

Документация доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

---

## 🏁 Как запустить проект

### 1. Клонирование репозитория

```bash
git clone https://github.com/nemopss/financial-tracker.git
cd financial-tracker
```

### 2. Настройка переменных окружения

Создайте файл `.env` на основе схемы:

```plaintext
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=financial_tracker
JWT_SECRET=your_secret_key
PORT=8080
```

### 3. Установка `goose`

Убедитесь, что `goose` установлен. Если нет, выполните:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 4. Выполнение миграций

Примените миграции для настройки базы данных:

```bash
goose -dir migrations postgres "user=your_user password=your_password dbname=financial_tracker sslmode=disable" up
```

- **`-dir migrations`** — папка с миграциями.
- **Подключение к PostgreSQL**: Убедитесь, что данные соответствуют вашему `.env`.

### 5. Запуск сервера

После успешных миграций запустите сервер:

```bash
go run cmd/main.go
```

---

## ✅ Уже сделано

- [x] Перенос проекта на Gin
- [x] Подключение Swagger UI
- [x] Реализация CRUD для категорий и транзакций
- [x] Аналитика доходов и расходов
- [x] Поддержка JWT-аутентификации
- [x] SQL миграции для базы данных

---

## 📝 Задачи на будущее

- [ ] 🛠 Переделать юнит-тесты под Gin
- [ ] ⚙️ Добавить интеграционные тесты
- [ ] 🧬 Настроить CI/CD
- [ ] 💀 Реализовать фронтенд-составляющую
- [ ] 🐳 Добавить деплой
- [ ] 🔗 Реализовать версионирование API
- [ ] 🔐 Улучшить безопасность (Rate Limiting, CORS)
