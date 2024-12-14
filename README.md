## Аутентификация JWT

Сервис аутентификации JWT со следующими требованиями и исполнением:
1. При аутентификации выдается пара токенов, Access&Refresh
2. Алгоритм шифрования Access токена SHA512, в базе не хранится.
3. Тип Refresh токена - JWT c полями payload'a: "uuid", "addrIP", "expiredAt", "createdAt".
4. В базе хранится время генерации рефреш токена и его хеш последних 50 байтов.
5. Осуществляется сверка IP адресов, откуда осуществляется запрос на аутентификацию и рефреш токенов. В случае отличия - отправляется сообщение на почту пользователя
6. Запуск приложения, БД, миграции осуществлен в контейнерах с помощью docker compose.

### Выполненные эндпоинты
* GET    /api/v1/auth/signin/:uuid - аутентификация с параметром uuid пользователя
* POST   /api/v1/auth/refresh - в теле запроса передается refresh токен

### Запуск приложения

1. Склонировать репозиторий
2. Изменить название файла ".env.example" на ".env"
3. Заполнить его соответствующими данными
4. выполнить ```docker compose up```