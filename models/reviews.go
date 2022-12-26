package models

import "time"

type Review struct {
	Id      int       `json:"id" db:"id"`
	Date    time.Time `json:"date" db:"date"`
	UserId  int       `json:"userId" db:"user_id"`
	MovieId int       `json:"movieId" db:"movie_id"`
	Review  string    `json:"review" db:"review"`
	Rate    int       `json:"rate" db:"rate"`
}

type ReviewRequest struct {
	MovieId int    `json:"movieId" db:"movie_id"`
	Review  string `json:"review" db:"review"`
	Rate    int    `json:"rate" db:"rate"`
}

type ReviewsCountAndRaters struct {
	Count    int `db:"count"`
	RatesSum int `db:"total"`
}
