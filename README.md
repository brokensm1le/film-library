# Film Library
___________________________

REST API приложение, который предоставляет управления БД фильмов.

## System Design

Само приложение дожидается пока БД поднимется на 5432 порту, после чего поднимается в Docker контейнере на порту 8080.

Чтобы **поднять приложение** с БД(Имплементацию всех make-запросов можно найти в Makefile):
```
 make buildrun
```

Чтобы останосить приложение с БД:
```
 make stop
```

Чтобы запустить unit tests:
```
 make test
```

Чтобы посмотреть покрытие unit test-ами:
```
 make cover
```

Для генерации mock-аных частей кода(для mock тестов):
```
 make gen
```

## API:

 Весь API описан в [/docs]{https://github.com/brokensm1le/film-library/tree/master/docs} и также приложен [postman collection]{https://github.com/brokensm1le/film-library/blob/master/New%20Collection.postman_collection.json}.
