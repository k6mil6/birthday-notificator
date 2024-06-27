# Birthday-notificator

Сервис для рассылки уведомлений пользователям об предстоящих днях рождения тех людей, на которых
подписан сам пользователь.


# Шаги для тестирования
1. зайти в консоль, запуллить репозиторий командой
```
git pull https://github.com/k6mil6/film-bot.git
```
2. перейти в корень репозитория, запустить тесты командой
```
go test ./...
```
3. для запуска проекта необходимо иметь установленный [docker(docker-compose)](https://www.docker.com/products/docker-desktop/)
4. при необходимости настроить http порт в ./config/config.yaml и в docker-compose.yml
5. поднять проект командной
```
docker-compose up -d
```

После перечисленных шагов сервис готов к работе, можно обратится к папке ./examples чтобы найти примеры взаимодействия 
с api

## Система каталогов
```
│   .gitignore
│   .golangci.yml - настройки линтера
│   docker-compose.dev.yml - для запуска бд
│   docker-compose.yml - для запуска всего проекта
│   Dockerfile - докерфайл собирает сам проект
│   go.mod
│   go.sum
│   README.md
│   Taskfile.yml
│   TODO.md
│
├───cmd
│   └───api
│           main.go - место откуда все запускается
│
├───config
│       config.yaml - конфиг для работы в докере
│       dev-config.yaml - конфиг для запуска проекта без докера
│
├───examples
│   │   README.md - ридми для примеров
│   │
│   ├───http - примеры запросов в файлах .http
│   │   └───user
│   │       ├───change_email
│   │       │       change_email.http
│   │       │       description.md
│   │       │
│   │       ├───change_notification_date
│   │       │       change_notification_date.http
│   │       │       description.md
│   │       │
│   │       ├───login
│   │       │       description.md
│   │       │       login.http
│   │       │
│   │       ├───register
│   │       │       description.md
│   │       │       register.http
│   │       │
│   │       ├───subscribe
│   │       │       description.md
│   │       │       subscribe.http
│   │       │
│   │       ├───subscritpitons
│   │       │       description.md
│   │       │       subscriptions.http
│   │       │
│   │       ├───unsubscribe
│   │       │       description.md
│   │       │       unsubscribe.http
│   │       │
│   │       └───user
│   │               descritption.md
│   │               user.http
│   │
│   └───postman - содержит коллекцию постман
│           birthday-notificator.postman_collection.json
│
├───internal
│   ├───api
│   │   └───http
│   │       │   server.go
│   │       │
│   │       ├───handlers - папка содержащая хендлеры для эндпоинтов по имени
│   │       │   └───user
│   │       │       │   user.go
│   │       │       │
│   │       │       ├───change
│   │       │       │   ├───email
│   │       │       │   │       email.go
│   │       │       │   │
│   │       │       │   └───notification_date
│   │       │       │           notification_date.go
│   │       │       │
│   │       │       ├───login
│   │       │       │       login.go
│   │       │       │
│   │       │       ├───register
│   │       │       │       register.go
│   │       │       │
│   │       │       ├───subscribe
│   │       │       │       subscribe.go
│   │       │       │
│   │       │       ├───subscriptions
│   │       │       │       subscriptions.go
│   │       │       │
│   │       │       └───unsubscribe
│   │       │               unsubscribe.go
│   │       │
│   │       ├───middleware - 
│   │       │   ├───identity
│   │       │   │       identity.go - middleware для получения id пользователя из jwt токена
│   │       │   │
│   │       │   └───logger
│   │       │           logger.go - логирует запросы к апи
│   │       │
│   │       └───response
│   │               response.go - содержит формат респонсов апи
│   │
│   ├───app
│   │   │   app.go - здесь собираются все сервисы приложения
│   │   │
│   │   └───api
│   │           api.go - здесь собирается апи
│   │
│   ├───config
│   │       config.go - файл для преобразования файлов конфигурации
│   │
│   ├───lib
│   │   ├───birthday
│   │   │       birthday.go - определяет когда следующий день рождения у пользователя
│   │   │       birthday_test.go
│   │   │
│   │   ├───email
│   │   │   │   email.go - пакет для работы с электронной почтой (валидация, отправка)
│   │   │   │   email_test.go
│   │   │   │
│   │   │   └───text
│   │   │           text.go - темплейты хэдеров и бади письма
│   │   │
│   │   ├───jwt
│   │   │       jwt.go - пакет для генерации jwt токенов и получения из них id пользователя
│   │   │       jwt_test.go
│   │   │
│   │   ├───logger - содержит файлы для настройки slog
│   │   │   │   logger.go
│   │   │   │
│   │   │   └───handlers
│   │   │       └───slogpretty
│   │   │               slogpretty.go
│   │   │
│   │   └───notification
│   │       └───offset
│   │               offset.go - пакет для преобразования структуры offset
│   │               offset_test.go
│   │
│   ├───model
│   │       model.go - содержит все сущности, которые используются в приложении
│   │
│   ├───service
│   │   ├───auth
│   │   │       auth.go - сервисный слой для авторизации пользователя
│   │   │       errors.go
│   │   │
│   │   ├───notification
│   │   │       notification.go - сервисный слой для отправки уведомлений пользователям
│   │   │
│   │   └───user
│   │       └───interaction
│   │               errors.go
│   │               interaction.go - сервисный слой для всех действий пользователя
│   │
│   └───storage 
│       └───postgres
│           │   postgres.go
│           │
│           ├───subscriptions - содержит слой для общения с бд таблицей subscriptions
│           │       errors.go
│           │       subscriptions.go
│           │
│           └───users - содержит слой для общения с бд таблицей users
│                   errors.go
│                   users.go
│
└───migrations - содержит все миграции в проекте
        000001_create_users_table.down.sql
        000001_create_users_table.up.sql
        000002_create_subscriptions_table.down.sql
        000002_create_subscriptions_table.up.sql
        000003_insert_dummy_users.down.sql
        000003_insert_dummy_users.up.sql

```