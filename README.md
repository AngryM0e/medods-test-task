# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`


## Решения для тестового задания:
- Для структуры задачи добавлено поле даты завершения задачи end_date и поле повтора задачи repeat_type.

- Если у задачи есть повторяемость (заполненное поле repeat_type), то по её завершению автоматически создаётся её копия с статусом "new", скопированными полями title, description и repeat_type из исходной задачи и перерасчитанной датой завершению с использованием модуля scheduler.

- Правила повторения задач включают в себя: 

  1. daily N — каждые N дней (1–400)

  2. monthly day — по указанным дням месяца (поддержка -1 и -2 для последних дней)

  3. even / odd — чётные/нечётные дни

  4. specific [date1,date2,...] — ближайшая будущая дата из списка (формат YYYYMMDD)

- Добавлены тесты для всех типов повторений, включая граничные случаи (переход месяцев, расчёт от дат в прошлом).

