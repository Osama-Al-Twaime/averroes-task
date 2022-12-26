package models

import "time"

type Watchlist struct {
	Id      int       `json:"id" db:"id"`
	Date    time.Time `json:"date" db:"date"`
	UserId  int       `json:"userId" db:"user_id"`
	MovieId int       `json:"movieId" db:"movie_id"`
}

type AddToWatchlistRequest struct {
	MovieId int `json:"movieId" db:"movie_id"`
}
