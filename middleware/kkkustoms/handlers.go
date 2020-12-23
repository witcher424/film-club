package main

import (
    "net/http"
    "strconv"
    "encoding/json"
)

func internalError(w http.ResponseWriter, err error) {
    logger.Println(err)
    http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func showDbQueryError(w http.ResponseWriter, dbErr *DbQueryError) {
    if dbErr.Type == DB_INTERNAL_ERROR {
        internalError(w, dbErr.Err)
    } else {
        http.Error(w, dbErr.Err.Error(), http.StatusBadRequest)
    }
}

func isValidLogin(login string) bool {
    for i := 0; i < len(login); i++ {
        if (login[i] < 'A' || login[i] > 'Z') &&
           (login[i] < 'a' || login[i] > 'z') &&
           (login[i] < '0' || login[i] > '9') &&
           login[i] != '_' {
               return false
           }
    }
    return true
}

func authHandler(
    hdl func (http.ResponseWriter, *http.Request, int),
) func (http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, req *http.Request) {
        login := req.FormValue("login")
        password := req.FormValue("password")

        userId, dbErr := checkAuth(login, sha256d([]byte(password)))
        if dbErr != nil {
            showDbQueryError(w, dbErr)
            return
        }

        if userId < 0 {
            http.Error(w, "Authentication is failed", http.StatusUnauthorized)
            return
        }

        hdl(w, req, userId)
    }
}

func gendersHandler(w http.ResponseWriter, req *http.Request) {
    genders, dbErr := getGenders()
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }

    err := json.NewEncoder(w).Encode(genders)
    if err != nil {
        logger.Println(err)
    }
}

func occupationsHandler(w http.ResponseWriter, req *http.Request) {
    occupations, dbErr := getOccupations()
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }

    err := json.NewEncoder(w).Encode(occupations)
    if err != nil {
        logger.Println(err)
    }
}

func genresHandler(w http.ResponseWriter, req *http.Request) {
    genres, dbErr := getGenres()
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }

    err := json.NewEncoder(w).Encode(genres)
    if err != nil {
        logger.Println(err)
    }
}

func registerHandler(w http.ResponseWriter, req *http.Request) {
    login := req.FormValue("login")
    password := req.FormValue("password")
    gender := req.FormValue("gender")
    occupation := req.FormValue("occupation")

    if len(login) < 3 || len(login) > 15 {
        http.Error(w, "login must be at least 3 and at most 15 characters long",
                   http.StatusBadRequest)
        return
    }

    if !isValidLogin(login) {
        http.Error(w, "login must contain only A-Za-z0-9_ characters",
                   http.StatusBadRequest)
        return
    }

    if len(password) < 8 {
        http.Error(w, "password must be at least 8 characters long",
                   http.StatusBadRequest)
        return
    }

    age, err := strconv.Atoi(req.FormValue("age"))
    if err != nil {
        http.Error(w, "age must be a number", http.StatusBadRequest)
        return
    }

    if age < 7 || age > 120 {
        http.Error(w, "age must be greater than 6 and lower than 121",
                   http.StatusBadRequest)
        return
    }

    dbErr := insertUser(login, sha256d([]byte(password)), age, gender, occupation)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }
}

func changePasswordHandler(w http.ResponseWriter, req *http.Request, userId int) {
    newPassword := req.FormValue("new_password")

    if len(newPassword) < 8 {
        http.Error(w, "new_password must be at least 8 characters long",
                   http.StatusBadRequest)
        return
    }

    dbErr := changePassword(userId, sha256d([]byte(newPassword)))
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }
}

func deleteAccountHandler(w http.ResponseWriter, req *http.Request, userId int) {
    dbErr := deleteAccount(userId)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }
}

func userIdHandler(w http.ResponseWriter, req *http.Request) {
    userId, err := strconv.Atoi(req.FormValue("user_id"))
    if err != nil {
        http.Error(w, "user_id must be a number", http.StatusBadRequest)
        return
    }

    user, dbErr := userById(userId)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }

    err = json.NewEncoder(w).Encode(user)
    if err != nil {
        logger.Println(err)
    }
}

func usersLoginHandler(w http.ResponseWriter, req *http.Request) {
    login := req.FormValue("login")

    users, dbErr := usersByLogin(login)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }

    err := json.NewEncoder(w).Encode(users)
    if err != nil {
        logger.Println(err)
    }
}

func movieIdHandler(w http.ResponseWriter, req *http.Request) {
    movieId, err := strconv.Atoi(req.FormValue("movie_id"))
    if err != nil {
        http.Error(w, "movie_id must be a number", http.StatusBadRequest)
        return
    }

    movie, dbErr := movieById(movieId)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }

    err = json.NewEncoder(w).Encode(movie)
    if err != nil {
        logger.Println(err)
    }
}

func moviesTitleHandler(w http.ResponseWriter, req *http.Request) {
    title := req.FormValue("title")

    movies, dbErr := moviesByTitle(title)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }

    err := json.NewEncoder(w).Encode(movies)
    if err != nil {
        logger.Println(err)
    }
}

func moviesUserHandler(w http.ResponseWriter, req *http.Request) {
    userId, err := strconv.Atoi(req.FormValue("user_id"))
    if err != nil {
        http.Error(w, "user_id must be a number", http.StatusBadRequest)
        return
    }

    movies, dbErr := moviesByUser(userId)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }

    err = json.NewEncoder(w).Encode(movies)
    if err != nil {
        logger.Println(err)
    }
}

func moviesGenreHandler(w http.ResponseWriter, req *http.Request) {
    genre := req.FormValue("genre")

    movies, dbErr := moviesByGenre(genre)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }

    err := json.NewEncoder(w).Encode(movies)
    if err != nil {
        logger.Println(err)
    }
}

func moviesTopHandler(w http.ResponseWriter, req *http.Request) {
    movies, dbErr :=  getTopMovies()
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }

    err := json.NewEncoder(w).Encode(movies)
    if err != nil {
        logger.Println(err)
    }
}

func movieRatingsHandler(w http.ResponseWriter, req *http.Request) {
    movieId, err := strconv.Atoi(req.FormValue("movie_id"))
    if err != nil {
        http.Error(w, "movie_id must be a number", http.StatusBadRequest)
        return
    }

    ratings, dbErr := movieRatings(movieId)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
        return
    }

    err = json.NewEncoder(w).Encode(ratings)
    if err != nil {
        logger.Println(err)
    }
}

func insertRatingHandler(w http.ResponseWriter, req *http.Request, userId int) {
    rating, err := strconv.Atoi(req.FormValue("rating"))
    if err != nil {
        http.Error(w, "rating must be a number", http.StatusBadRequest)
        return
    }

    movieId, err := strconv.Atoi(req.FormValue("movie_id"))
    if err != nil {
        http.Error(w, "movie_id must be a number", http.StatusBadRequest)
        return
    }

    dbErr := insertRating(userId, movieId, rating)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }
}

func insertCommentHandler(w http.ResponseWriter, req *http.Request, userId int) {
    movieId, err := strconv.Atoi(req.FormValue("movie_id"))
    if err != nil {
        http.Error(w, "movie_id must be a number", http.StatusBadRequest)
        return
    }

    content := req.FormValue("content")
    if len(content) < 10 {
        http.Error(w, "content must be at least 10 characters long",
                   http.StatusBadRequest)
        return
    }

    dbErr := insertComment(userId, movieId, content)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }
}

func deleteCommentHandler(w http.ResponseWriter, req *http.Request, userId int) {
    movieId, err := strconv.Atoi(req.FormValue("movie_id"))
    if err != nil {
        http.Error(w, "movie_id must be a number", http.StatusBadRequest)
        return
    }

    dbErr := deleteComment(userId, movieId)
    if dbErr != nil {
        showDbQueryError(w, dbErr)
    }
}
