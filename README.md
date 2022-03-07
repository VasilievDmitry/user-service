# UserService

[![Build Status](https://github.com/lotproject/user-service/workflows/Build/badge.svg?branch=develop)](https://github.com/lotproject/user-service/actions)
[![Code Coverage](https://codecov.io/gh/lotproject/user-service/branch/develop/graph/badge.svg?token=PXQHLW26AY)](https://codecov.io/gh/lotproject/user-service)

Сервис предназначен для регистрации, авторизации и управление данными пользователя.

### Переменные окружения

| Name                        | Required | Default        | Description                                                                                 |
|:----------------------------|:---------|:---------------|:--------------------------------------------------------------------------------------------|
| DEVELOP_MODE                |          | false          | Активация настроек для локальной разработки (отключение csrf)                               |
| METRICS_PORT                |          | 8086           | Порт приложения для сервиса метрики                                                         |
| METRICS_READ_TIMEOUT        |          | 60             | Таймаут для получения метрик                                                                |
| METRICS_READ_HEADER_TIMEOUT |          | 60             | Таймаут для получения метрик                                                                |
| LOG_FILE_PATH               |          | ./logs/log.txt | Путь до лог-файла                                                                           |
| LOG_LEVEL                   |          | error          | Уровень логирования                                                                         |
| LOG_TO_FILE_ENABLED         |          | false          | Активация логирования в файл                                                                |
| MYSQL_DSN                   | *        |                | DSN строка для подключения к БД MySQL (параметр parseTime=true обязателен)                  |
| MIGRATIONS_LOCK_TIMEOUT     |          | 120            | Таймаут для выполнения скриптов миграции БД                                                 |
| BCRYPT_COST                 |          | 10             | Сложность алгоритма шифрования токенов                                                      |
| REFRESH_TOKEN_LIFETIME      |          | 365            | Время жизни refresh-токена (дней)                                                           |
| ACCESS_TOKEN_LIFETIME       |          | 3              | Время жизни access-токена (часов)                                                           |
| ACCESS_TOKEN_SECRET         | *        |                | Серкертный ключ для подписи токенов                                                         |
| ACCESS_TOKEN_SIGNING_METHOD |          | HS256          | Алгоритм хэширования токенов                                                                |
| MICRO_TRANSPORT             |          | http           | Транспорт общения микросервисов. Доступные значения: http, grpc                             |
| MICRO_REGISTRY              |          | mdns           | Discovery-сервис для микросервисов. Доступные занчения: etcd, mdns                          |
| MICRO_REGISTRY_ADDRESS      |          |                | Хост и порт discovery-сервиса. Можно использовать множественные значения с разделителем `;` |
