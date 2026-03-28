# Request Managing App

[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/)
[![Test](https://github.com/matveevaolga/request-managing-app/actions/workflows/test.yml/badge.svg)](https://github.com/matveevaolga/request-managing-app/actions/workflows/test.yml)
[![Lint](https://github.com/matveevaolga/request-managing-app/actions/workflows/lint.yml/badge.svg)](https://github.com/matveevaolga/request-managing-app/actions/workflows/lint.yml)
[![Docker](https://github.com/matveevaolga/request-managing-app/actions/workflows/docker.yml/badge.svg)](https://github.com/matveevaolga/request-managing-app/actions/workflows/docker.yml)
[![codecov](https://codecov.io/gh/matveevaolga/request-managing-app/branch/main/graph/badge.svg)](https://codecov.io/gh/matveevaolga/request-managing-app)

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