# backend-bootcamp-task

[![ci.yml](https://github.com/dvkonovalov/backend-bootcamp-task/actions/workflows/ci.yml/badge.svg)](https://github.com/dvkonovalov/backend-bootcamp-task/actions/workflows/ci.yml)

# Go Service with PostgreSQL

Этот проект содержит сервис Go, которая взаимодействует с базой данных PostgreSQL. Обе службы упакованы с помощью Docker и управляются с помощью Docker Compose.
Сервис выполняет функционал в соответствие с поставленным техническим заданием.

## Запуск

Выполните следующие действия, чтобы запустить проект с помощью Docker:

### 1. Клонируйте репозиторий

Клонируйте репозиторий на свой локальный компьютер:

```bash
git clone https://github.com/dvkonovalov/backend-bootcamp-task.git
cd backend-bootcamp-task
```

### 2. Создайте и запустите контейнеры

Используйте Docker Compose для создания и запуска службы Go и базы данных PostgreSQL:

```bash
docker-compose up --build
```

### 3. Получение доступа к сервису
   Как только контейнеры будут запущены, сервис станет доступен по адресу:

```bash
http://localhost:8080
```

### 4. Остановка контейнеров
Чтобы остановить контейнеры, используйте:

```bash
docker-compose down
```

## Вопросы к условию и их рещения
### 1. Работать с сервисом могут несколько модераторов. При этом конкретную квартиру может проверять только один модератор. Перед началом работы нужно перевести квартиру в статус on moderate — тем самым запретив брать её на проверку другим модераторам. Как хранить связь квартиры на модерации и конкретного модератора, работающего с ней?
Решение: Создать таблицу Moderation с полями - id (уникальный номер записи), flat_id(уникальный идекнтификатор квартиры, которую модерируют) и username(username модератора, взявшего в работу квартиру)

### 2. Когда/как часто делать расслыку пользователям по их подписке?
Решение: отправка пользователям информации о том, что появилась новая квартира в доме сразу после модерации и присвоения квартире статуса "approved"

### 3. Как хранить подписки пользователей?
Решение: Создать таблицу CREATE TABLE house_subscriptions с полями - id (уникальный номер записи), house_id (номер дома для подписки),
user_email (email пользователя, куда отправлять уведомления).