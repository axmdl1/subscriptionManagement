# Subscription Management Service

REST-сервис для агрегации и расчёта стоимости онлайн-подписок пользователей.

Реализовано в рамках тестового задания Junior Golang Developer.

---

## Возможности

- CRUDL операции над подписками
- Подсчёт суммарной стоимости подписок за период
- Фильтрация по пользователю и названию сервиса
- Swagger документация
- PostgreSQL + миграции
- Конфигурация через env
- Структурированное логирование (zap)
- Запуск одной командой через Docker Compose

---

## Стек технологий

- Go 1.24
- Gin
- PostgreSQL
- pgx
- golang-migrate
- swaggo / Swagger
- zap logger
- Docker / Docker Compose

---

## Архитектура

Проект построен по принципам слоистой архитектуры:

 - cmd/app — точка входа
 - internal/config — конфигурация
 - internal/handler — HTTP слой
 - internal/service — бизнес-логика
 - internal/repository— работа с БД
 - internal/model — доменные модели
 - internal/utils — вспомогательные функции
 - internal/logger — логирование
 - migrations — SQL миграции
 - docs — swagger документация

## Запуск приложения

Приложение полностью контейнеризировано и запускается одной командой.

### Требования
- Docker >= 20
- Docker Compose (входит в Docker Desktop)

---

### 1. Склонировать репозиторий

``` git clone https://github.com/axmdl1/subscriptionManagement.git ```
 
``` cd subscriptionManagement```


---

### 2. Запуск

``` docker compose up --build ```


При первом запуске произойдёт:

1. Поднимется PostgreSQL
2. Приложение дождётся готовности БД
3. Автоматически применятся миграции
4. Запустится HTTP сервер

---

### 3. Проверка работы

После запуска доступны:

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

---

### Остановка

``` docker compose down ```

Удалить данные БД:

``` docker compose down -v ```
