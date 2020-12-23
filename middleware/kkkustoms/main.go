package main

import (
    "context"
    "crypto/sha256"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jackc/pgx/pgxpool"
    "io/ioutil"
    "log"
    "net"
    "net/http/fcgi"
    "os"
)

const configPath = "/etc/kkkustoms.json"

var (
    dbpool *pgxpool.Pool
    logger *log.Logger
)

type Config struct {
    SocketPath string
    DbUrl      string
}

func sha256d(data []byte) []byte {
    result := sha256.Sum256(data)
    result = sha256.Sum256(result[:])
    return result[:]
}

func main() {
    // Init logger
    logger = log.New(os.Stderr, "", log.Flags())

    // Read config file
    buf, err := ioutil.ReadFile(configPath)
    if err != nil {
        logger.Fatalln(err)
    }

    var config Config
    err = json.Unmarshal(buf, &config)
    if err != nil {
        logger.Fatalln(err)
    }
    
    // Remove socket file
    os.Remove(config.SocketPath)

    // Init router
    router := mux.NewRouter()
    router.HandleFunc("/genders", gendersHandler)
    router.HandleFunc("/occupations", occupationsHandler)
    router.HandleFunc("/genres", genresHandler)
    router.HandleFunc("/register", registerHandler)
    router.HandleFunc("/change_password", authHandler(changePasswordHandler))
    router.HandleFunc("/delete_account", authHandler(deleteAccountHandler))
    router.HandleFunc("/user_id", userIdHandler)
    router.HandleFunc("/users_login", usersLoginHandler)
    router.HandleFunc("/movie_id", movieIdHandler)
    router.HandleFunc("/movies_title", moviesTitleHandler)
    router.HandleFunc("/movies_user", moviesUserHandler)
    router.HandleFunc("/movies_genre", moviesGenreHandler)
    router.HandleFunc("/movies_top", moviesTopHandler)
    router.HandleFunc("/movie_ratings", movieRatingsHandler)
    router.HandleFunc("/insert_rating", authHandler(insertRatingHandler))
    router.HandleFunc("/insert_comment", authHandler(insertCommentHandler))
    router.HandleFunc("/delete_comment", authHandler(deleteCommentHandler))

    // Init db pool
    dbpool, err = pgxpool.Connect(context.Background(), config.DbUrl)
    if err != nil {
        logger.Fatalln(err)
    }
    defer dbpool.Close()

    // Run fastcgi over unix domain socket
    ln, err := net.Listen("unix", config.SocketPath)
    if err != nil {
        logger.Fatalln(err)
    }
    defer ln.Close()

    // Run fcgi
    logger.Println("Server is started")
    logger.Fatalln(fcgi.Serve(ln, router))
}
