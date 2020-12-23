# kkkustoms

Промежуточный слой между QT приложением и базой данной.

## Запросы

На все запросы возвращается ответ в формате json. Список запросов приведён ниже:

| Запрос           | Аутентификация | Аргументы                           | Результат                                                              |
|------------------|----------------|-------------------------------------|------------------------------------------------------------------------|
| /genders         | -              | -                                   | Возвращает массив валидных значений графы пол                          |
| /occupations     | -              | -                                   | Возвращает массив валидных значений графы занятость                    |
| /genres          | -              | -                                   | Возвращает список жанров                                               |
| /register        | -              | login, password, gender, occupation | Создаёт учётную запись пользователя                                    |
| /change_password | +              | new_password                        | Устанавливает пароль пользователя в значение new_password              |
| /delete_account  | +              | -                                   | Удаляет аккаунт пользователя                                           |
| /user_id         | -              | user_id                             | Возвращает информацию о пользователе с id равным user_id               |
| /users_login     | -              | login                               | Возвращает список пользователей, чьи логины имеют префикс login        |
| /movie_id        | -              | movie_id                            | Возвращает информацию о фильме с id равным movie_id                    |
| /movies_title    | -              | title                               | Возвращает список фильмов, чьи названия имеют префикс title            |
| /movies_user     | -              | user_id                             | Возвращает список фильмов, оценённых пользователем с id равным user_id |
| /movies_genre    | -              | genre                               | Возвращает список фильмов в жанре genre                                |
| /movies_top      | -              | -                                   | Возвращает топ 100 фильмов                                             |
| /movie_ratings   | -              | movie_id                            | Возвращает оценки и комментарии к фильму с id равным movie_id          |
| /insert_rating   | +              | movie_id, rating                    | Устанавливает оценку rating фильму с id равным movie_id                |
| /insert_comment  | +              | movie_id, content                   | Добавляет комментарий content к фильму с id равным movie_id            |
| /delete_comment  | +              | movie_id                            | Удаляет комментарий к фильму с id равным movie_id                      |

## Структуры

Ниже приведён список структур, которые будут преобразовываться в json.

```go
type User struct {
    UserId     int
    Login      string
    phash      []byte
    Age        int
    Gender     string
    Occupation string
}

type Movie struct {
    MovieId     int
    Title       string
    Description string
    PosterPath  string
    Genres      []string
    Rating      float32
    Comment     string
}

type Rating struct {
    UserId  int
    Login   string
    Rating  int
    Comment string
}
```