## NewsTG
Автоматизированный агрегатор новостей для Telegram.

Собирает контент из внешних источников, сохраняет в PostgreSQL и пересылает в каналы. Написан на Go, полностью контейнеризирован.

## Архитектура

Проект развернут в **Docker** и состоит из трех независимых сервисов:
- **PostgreSQL**: Хранилище данных.
- **Parser**: Модуль сбора новостей (API/HTTP).
- **Poster**: Модуль публикации новостей (Telegram).

Стек:
- **Language:** Golang
- **Database:** PostgreSQL
- **Infrastructure:** Docker/Docker-compose
- **Libraries:**
    - `pgx`                 (Database)
    - `go-telegram-bot-api` (Telegram bot)
    - `godotenv`            (Env variables)

Запуск
``` bash
git clone https://github.com/KEKACIK/NewsTG.git

docker compose up --force-recreate --build -d
```
