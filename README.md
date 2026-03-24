# Request Managing App

Сервис для приёма и отбора заявок на проекты от внешних инициаторов.

## Технологии

- Go 1.24
- PostgreSQL 15
- JWT для аутентификации

## Быстрый старт

```bash
# Скопировать конфиг
cp .env.example .env

# Запустить PostgreSQL
make docker-up

# Выполнить миграции
make migrate-up

# Заполнить тестовыми данными
make seed

# Запустить сервер
make run