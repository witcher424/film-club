package main

import (
    "context"
    "errors"
    "github.com/jackc/pgx"
)

const (
    DB_INTERNAL_ERROR = iota
    DB_USER_ERROR     = iota
)

var (
    loginExistsErr   = dbQueryErrorStr("Login is already exists", DB_USER_ERROR)
    unknownGenderErr = dbQueryErrorStr("Unknown gender", DB_USER_ERROR)
    unknownOccupErr  = dbQueryErrorStr("Unknown occupation", DB_USER_ERROR)
    unknownUserErr   = dbQueryErrorStr("Unknown user", DB_USER_ERROR)
    unknownMovieErr  = dbQueryErrorStr("Unknown movie", DB_USER_ERROR)
)

type DbQueryError struct {
    Err  error
    Type int
}

func dbQueryErrorStr(msg string, typ int) *DbQueryError {
    return &DbQueryError {
        Err:  errors.New(msg),
        Type: typ,
    }
}

func dbQueryErrorErr(err error, typ int) *DbQueryError {
    return &DbQueryError {
        Err:  err,
        Type: typ,
    }
}

func checkAuth(login string, phash []byte) (int, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return -1, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    var userId int
    err = conn.QueryRow(context.Background(),
                        "SELECT user_id FROM users WHERE login=$1 AND phash=$2",
                        login,
                        phash).Scan(&userId)

    if err == pgx.ErrNoRows {
        return -1, nil
    } else if err != nil {
        return -1, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }

    return userId, nil
}

func getGenders() ([]string, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(), "SELECT title FROM genders")
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var genders []string
    var gender string
    for rows.Next() {
        err = rows.Scan(&gender)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        genders = append(genders, gender)
    }

    return genders, nil
}

func getOccupations() ([]string, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(), "SELECT title FROM occupations")
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var occupations []string
    var occupation string
    for rows.Next() {
        err = rows.Scan(&occupation)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        occupations = append(occupations, occupation)
    }

    return occupations, nil
}

func getGenres() ([]string, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(), "SELECT title FROM genres")
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var genres []string
    var genre string
    for rows.Next() {
        err = rows.Scan(&genre)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        genres = append(genres, genre)
    }

    return genres, nil
}

func insertUser(login string, phash []byte, age int, gender,
                occupation string) *DbQueryError {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    var dummyId, genderId, occupationId int

    // Check if login is already exists
    err = conn.QueryRow(context.Background(), "SELECT user_id FROM users WHERE login=$1",
                        login).Scan(&dummyId)
    if err != nil {
        if err != pgx.ErrNoRows {
            return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
    } else {
        return loginExistsErr
    }

    // Check if gender is valid
    err = conn.QueryRow(context.Background(),
                        "SELECT gender_id FROM genders WHERE title=$1",
                        gender).Scan(&genderId)
    if err == pgx.ErrNoRows {
        return unknownGenderErr
    } else if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }

    // Check if occupation is valid
    err = conn.QueryRow(context.Background(),
                        "SELECT occupation_id FROM occupations WHERE title=$1",
                        occupation).Scan(&occupationId)
    if err == pgx.ErrNoRows {
        return unknownOccupErr
    } else if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }

    // Insert new user
    _, err = conn.Exec(context.Background(),
                       `INSERT INTO users (login, phash, age, gender, occupation)
                       VALUES ($1, $2, $3, $4, $5)`, login, phash, age, genderId,
                       occupationId)
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    return nil
}

func changePassword(userId int, phash []byte) *DbQueryError {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    _, err = conn.Exec(context.Background(),
                       "UPDATE users SET phash=$1 WHERE user_id=$2",
                       phash, userId)
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    return nil
}

func deleteAccount(userId int) *DbQueryError {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    _, err = conn.Exec(context.Background(), "DELETE FROM users WHERE user_id=$1", userId)
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    return nil
}

func unmarshalUser(row pgx.Row, u *User) error {
    return row.Scan(&u.UserId, &u.Login, &u.Age, &u.Gender, &u.Occupation)
}

func userById(userId int) (*User, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    user := &User{}
    row := conn.QueryRow(context.Background(),
                         `SELECT user_id, login, age, G.title, O.title
                         FROM users
                         INNER JOIN genders as G ON gender = gender_id
                         INNER JOIN occupations as O ON occupation = occupation_id
                         WHERE user_id = $1`,
                         userId)
    err = unmarshalUser(row, user)
    if err == pgx.ErrNoRows {
        return nil, unknownUserErr
    } else if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }

    return user, nil
}

func usersByLogin(login string) ([]User, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
                            `SELECT user_id, login, age, G.title, O.title
                            FROM users
                            INNER JOIN genders as G ON gender = gender_id
                            INNER JOIN occupations as O ON occupation = occupation_id
                            WHERE login ILIKE $1 || '%'`,
                            login)
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var users []User
    var user User
    for rows.Next() {
        err = unmarshalUser(rows, &user)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        users = append(users, user)
    }

    return users, nil
}

func unmarshalMovie(row pgx.Row, m *Movie) error {
    err := row.Scan(&m.MovieId, &m.Title, &m.Description, &m.Rating, &m.PosterPath)
    if err != nil {
        return err
    }

    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return err
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
                            `SELECT title
                            FROM movies_genres
                            INNER JOIN genres ON genre = genre_id
                            WHERE movie=$1`, m.MovieId)
    if err != nil {
        return err
    }
    defer rows.Close()

    var genre string
    m.Genres = nil
    for rows.Next() {
        err = rows.Scan(&genre)
        if err != nil {
            return err
        }

        m.Genres = append(m.Genres, genre)
    }

    return err
}

func loadRatingAndComment(userId int, m *Movie) error {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return err
    }
    defer conn.Release()

    var rating int
    err = conn.QueryRow(context.Background(),
                        "SELECT rating FROM ratings WHERE user_id=$1 AND movie=$2",
                        userId,
                        m.MovieId).Scan(&rating)
    if err != nil {
        return err
    }

    m.Rating = float32(rating)

    err = conn.QueryRow(context.Background(),
                        "SELECT content FROM comments WHERE user_id=$1 AND movie=$2",
                        userId,
                        m.MovieId).Scan(&m.Comment)

    if err != nil && err != pgx.ErrNoRows {
        return err
    }
    return nil
}

func movieById(movieId int) (*Movie, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    movie := &Movie{}
    row := conn.QueryRow(context.Background(),
                         `SELECT movie_id, title, description, mean_rating, poster_path
                         FROM movies
                         WHERE movie_id = $1`,
                         movieId)
    err = unmarshalMovie(row, movie)
    if err == pgx.ErrNoRows {
        return nil, unknownMovieErr
    } else if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    return movie, nil
}

func moviesByTitle(title string) ([]Movie, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
                            `SELECT movie_id, title, description, mean_rating, poster_path
                            FROM movies
                            WHERE title ILIKE $1 || '%'`,
                            title)
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var movies []Movie
    var movie Movie
    for rows.Next() {
        err = unmarshalMovie(rows, &movie)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        movies = append(movies, movie)
    }

    return movies, nil
}

func moviesByUser(userId int) ([]Movie, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
                            `SELECT movie_id, title, description,
                            mean_rating, poster_path
                            FROM ratings
                            INNER JOIN movies ON movie = movie_id
                            WHERE user_id=$1`,
                            userId)
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var movies []Movie
    var movie Movie
    for rows.Next() {
        err = unmarshalMovie(rows, &movie)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }

        err = loadRatingAndComment(userId, &movie)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }

        movies = append(movies, movie)
    }

    return movies, nil
}

func moviesByGenre(genre string) ([]Movie, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
                            `SELECT movie_id, movies.title, description,
                            mean_rating, poster_path
                            FROM movies_genres
                            INNER JOIN genres ON genre = genre_id
                            INNER JOIN movies ON movie = movie_id
                            WHERE genres.title=$1`,
                            genre)
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var movies []Movie
    var movie Movie
    for rows.Next() {
        err = unmarshalMovie(rows, &movie)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        movies = append(movies, movie)
    }

    return movies, nil
}

func getTopMovies() ([]Movie, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
                            `SELECT movie_id, title, description,
                            mean_rating, poster_path
                            FROM movies ORDER BY mean_rating LIMIT 100`)
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var movies []Movie
    var movie Movie
    for rows.Next() {
        err = unmarshalMovie(rows, &movie)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        movies = append(movies, movie)
    }

    return movies, nil
}

func unmarshalRating(row pgx.Row, movieId int, r *Rating) error {
    err := row.Scan(&r.UserId, &r.Login, &r.Rating)
    if err != nil {
        return err
    }

    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return err
    }
    defer conn.Release()

    err = conn.QueryRow(context.Background(),
                        "SELECT content FROM comments WHERE user_id=$1 AND movie=$2",
                        r.UserId, movieId).Scan(&r.Comment)
    if err != nil && err != pgx.ErrNoRows {
        return err
    }
    return nil
}

func movieRatings(movieId int) ([]Rating, *DbQueryError) {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
                            `SELECT U.user_id, login, rating
                            FROM ratings as R
                            INNER JOIN users as U ON U.user_id = R.user_id
                            WHERE movie=$1`,
                            movieId)
    if err == pgx.ErrNoRows {
        return nil, nil 
    } else if err != nil {
        return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer rows.Close()

    var ratings []Rating
    var rating Rating
    for rows.Next() {
        err = unmarshalRating(rows, movieId, &rating)
        if err != nil {
            return nil, dbQueryErrorErr(err, DB_INTERNAL_ERROR)
        }
        ratings = append(ratings, rating)
    }

    return ratings, nil
}

func insertRating(userId, movieId, rating int) *DbQueryError {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    _, err = conn.Exec(context.Background(),
                       `INSERT INTO ratings (user_id, movie, rating)
                       VALUES ($1, $2, $3)
                       ON CONFLICT (user_id, movie) DO UPDATE
                       SET rating = excluded.rating`,
                       userId, movieId, rating)
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }

    return nil
}

func insertComment(userId, movieId int, content string) *DbQueryError {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    _, err = conn.Exec(context.Background(),
                       `INSERT INTO comments (user_id, movie, content)
                       VALUES ($1, $2, $3)
                       ON CONFLICT (user_id, movie) DO UPDATE
                       SET content = excluded.content`,
                       userId, movieId, content)
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    return nil
}

func deleteComment(userId, movieId int) *DbQueryError {
    conn, err := dbpool.Acquire(context.Background())
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    defer conn.Release()

    _, err = conn.Exec(context.Background(),
                       "DELETE FROM comments WHERE user_id=$1 AND movie=$2",
                       userId, movieId)
    if err != nil {
        return dbQueryErrorErr(err, DB_INTERNAL_ERROR)
    }
    return nil
}
