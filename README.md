# Request Managing App

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org/)
[![Test](https://github.com/matveevaolga/request-managing-app/actions/workflows/test.yml/badge.svg)](https://github.com/matveevaolga/request-managing-app/actions/workflows/test.yml)
[![Lint](https://github.com/matveevaolga/request-managing-app/actions/workflows/lint.yml/badge.svg)](https://github.com/matveevaolga/request-managing-app/actions/workflows/lint.yml)
[![Docker](https://github.com/matveevaolga/request-managing-app/actions/workflows/docker.yml/badge.svg)](https://github.com/matveevaolga/request-managing-app/actions/workflows/docker.yml)

Сервис для приёма и отбора заявок на проекты от внешних инициаторов.

## Технологии

- Go 1.25
- PostgreSQL 15
- JWT для аутентификации
- Docker, Docker Compose
- GitHub Actions (CI/CD)
- golangci-lint

## Архитектура

Проект построен на принципах Clean Architecture:

- **internal/domain** - бизнес-сущности и интерфейсы репозиториев
- **internal/service** - сценарии использования (бизнес-логика)
- **internal/repository** - реализация работы с PostgreSQL
- **internal/transport** - HTTP handlers, DTO, middleware
- **migrations** - миграции базы данных
- **seeds** - автоматическое заполнение тестовыми данными

## Быстрый старт

### Локальная разработка

1. Скопировать и заполнить конфиг:
   ```bash
   cp .env.example .env
   ```

2. Запустить все сервисы (PostgreSQL, миграции, seed, приложение):
   ```bash
   make docker-up
   ```

Сервер будет доступен по адресу: `http://localhost:8000`

## Команды Makefile

- **make run** - запустить сервер локально (требует запущенной БД)
- **make build** - собрать бинарный файл
- **make test** - запустить тесты с race detector
- **make test-coverage** - запустить тесты с отчетом о покрытии
- **make docker-up** - запустить все контейнеры в правильном порядке:
  1. PostgreSQL (ждет готовности)
  2. Миграции базы данных
  3. Заполнение тестовыми данными (seed)
  4. Приложение на порту 8000
- **make docker-down** - остановить и удалить контейнеры (данные сохраняются)
- **make docker-logs** - просмотреть логи Docker
- **make migrate-up** - применить миграции
- **make migrate-down** - откатить миграции
- **make migrate-create** - создать новую миграцию

## Тестовые пользователи

Для входа в систему доступны следующие учетные записи:

- Логин: `admin1`, пароль: `admin1` - роль ADMIN
- Логин: `admin2`, пароль: `admin2` - роль ADMIN
- Логин: `user1`, пароль: `user1` - роль USER
- Логин: `user2`, пароль: `user2` - роль USER

**Пользователи с ролью ADMIN (администраторы):**
- Авторизоваться в системе
- Просматривать список всех заявок
- Просматривать детальную информацию о любой заявке
- Принимать или отклонять заявки
- Просматривать типы проектов

**Пользователи с ролью USER (обычные пользователи):**
- Авторизоваться в системе
- Просматривать типы проектов
- Подавать новые заявки

**Неавторизованные пользователи (без токена):**
- Просматривать типы проектов
- Подавать новые заявки
- **Не могут** авторизоваться (для входа нужны учетные данные)
- **Не имеют доступа** к администрированию заявок
- **Не могут** просматривать чужие заявки

## API Endpoints

### Публичные эндпоинты (без аутентификации)

- **GET /health** - проверка работоспособности сервиса
- **POST /login** - авторизация (получение JWT токена)
- **GET /project/type** - список типов проектов
- **POST /project/application/external** - создание новой заявки

### Защищенные эндпоинты (требуют роль ADMIN)

- **GET /project/application/external/list** - список заявок с фильтрацией
- **GET /project/application/external/{applicationId}** - детальная информация о заявке
- **POST /project/application/external/{applicationId}/accept** - принять заявку
- **POST /project/application/external/{applicationId}/reject** - отклонить заявку (с указанием причины)

### Параметры фильтрации для списка заявок

- **active** - только заявки в статусе PENDING (boolean)
- **search** - поиск по названию проекта или ФИО (string)
- **projectTypeId** - фильтр по типу проекта (integer)
- **sortByDateUpdated** - сортировка по дате обновления (ASC/DESC)
- **limit** - количество записей на странице (по умолчанию 20)
- **offset** - смещение для пагинации (по умолчанию 0)

## Примеры запросов

### 1. Авторизация
```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"login":"admin1","password":"admin1"}'
```

### 2. Создание заявки
```bash
curl -X POST http://localhost:8000/project/application/external \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Jane Smith",
    "email": "jane@example.com",
    "organisationName": "Startup Inc",
    "projectName": "Mobile App",
    "typeId": 1,
    "expectedResults": "iOS and Android MVP",
    "isPayed": false
  }'
```

### 3. Получение списка заявок
```bash
TOKEN="ваш_jwt_токен"
curl -X GET "http://localhost:8000/project/application/external/list?active=true&limit=10" \
  -H "X-API-TOKEN: $TOKEN"
```

### 4. Принятие заявки
```bash
curl -X POST http://localhost:8000/project/application/external/1/accept \
  -H "X-API-TOKEN: $TOKEN"
```

### 5. Отклонение заявки
```bash
curl -X POST http://localhost:8000/project/application/external/1/reject \
  -H "Content-Type: application/json" \
  -H "X-API-TOKEN: $TOKEN" \
  -d '{"reason": "Не соответствует критериям"}'
```

## CI/CD

Проект настроен на автоматические проверки при каждом push:

- **Test** - запуск тестов с PostgreSQL. Выполняется при push в ветки: `main`, `develop`, `feature/*`, `ci/cd-integration`, а также при pull request в `main` и `develop`
- **Lint** - статический анализ кода через golangci-lint. Выполняется при push в ветки: `main`, `develop`, `feature/*`, `ci/cd-integration`, а также при pull request в `main` и `develop`
- **Docker** - сборка и публикация Docker образа. Выполняется при push в ветки: `main`, `develop`, `feature/*`, `ci/cd-integration`, а также при создании тегов `v*`. Образ публикуется в Docker Hub под именем `{username}/request-managing-app` с тегами: версия (semver), название ветки и короткий хеш коммита.

### Локальное тестирование workflows

Для локального запуска GitHub Actions workflows используйте [act](https://github.com/nektos/act).

1. Скопировать файл с секретами:
   ```bash
   cp .env.secrets.example .env.secrets
   ```

2. Заполнить секреты в `.env.secrets`:
   ```
   DOCKER_USERNAME=your-dockerhub-username
   DOCKER_PASSWORD=your-dockerhub-token
   ```

3. Запустить конкретный workflow:
   ```bash
   act -j test --secret-file .env.secrets
   act -j lint --secret-file .env.secrets
   act -j build --secret-file .env.secrets
   ```

**Примечание:** `.env.secrets` добавлен в `.gitignore` и не должен попадать в репозиторий.

## Структура проекта

- **cmd/server/** - точка входа (main.go)
- **internal/config/** - конфигурация из переменных окружения
- **internal/domain/** - бизнес-сущности и интерфейсы репозиториев
- **internal/repository/** - реализация работы с PostgreSQL
- **internal/service/** - бизнес-логика (сценарии использования)
- **internal/transport/handler/** - HTTP handlers и DTO
- **internal/transport/middleware/** - middleware (аутентификация, логирование)
- **internal/logger/** - настройка структурированного логирования (slog)
- **migrations/** - SQL миграции для создания схемы базы данных
- **seeds/** - автоматическое заполнение базы тестовыми данными
- **.github/workflows/** - CI/CD пайплайны (GitHub Actions)
- **docker-compose.yaml** - конфигурация для запуска всех сервисов
- **Dockerfile** - инструкция для сборки Docker образа
- **Makefile** - автоматизация команд (запуск, тесты, миграции)
- **README.md** - документация проекта
