package main

import (
    "context"
    "encoding/json"
    "github.com/jackc/pgx"
    "io/ioutil"
    "log"
    "os"
)

const configPath = "/etc/kkkollect.json"

type Config struct {
    DbUrl   string
    LogFile string
}

func main() {
    configBuf, err := ioutil.ReadFile(configPath)
    if err != nil {
        panic(err)
    }

    var config Config
    err = json.Unmarshal(configBuf, &config)
    if err != nil {
        panic(err)
    }

    logFd, err := os.OpenFile(config.LogFile, os.O_CREATE | os.O_APPEND | os.O_WRONLY,
                              0644)
    if err != nil {
        panic(err)
    }
    defer logFd.Close()
    logger := log.New(logFd, "", log.Flags())

    conn, err := pgx.Connect(context.Background(), config.DbUrl)
    if err != nil {
        logger.Fatalln(err)
    }
    defer conn.Close(context.Background())

    _, err = conn.Exec(context.Background(),
                       `UPDATE movies
                       SET mean_rating=
                       (SELECT ROUND(AVG(rating),2) FROM ratings WHERE movie=movie_id),
                       is_rating_updated=false
                       WHERE is_rating_updated=true`)
    if err != nil {
        logger.Fatalln(err)
    }
}
