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