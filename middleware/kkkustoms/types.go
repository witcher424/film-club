package main

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
