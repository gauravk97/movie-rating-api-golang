// model.go

package main

import (
	"database/sql"
)

type movie struct {
	ID            int                   `json:"id"`
	Name          string                `json:"movie_name"`
	AverageRating float64               `json:"average_rating"`
	NumRatings    int                   `json:"number_of_ratings"`
	MovieRatings  []user_movie_rating   `json:"movie_ratings"`
	MovieComments []user_movie_comments `json:"movie_comments"`
}

type user struct {
	ID           int                   `json:"id"`
	Name         string                `json:"user_name"`
	Role         int                   `json:"user_role"`
	UserRatings  []user_movie_rating   `json:"user_ratings"`
	UserComments []user_movie_comments `json:"user_comments"`
}

type user_movie_rating struct {
	UserID  int `json:"user_id"`
	MovieID int `json:"movie_id"`
	Rating  int `json:"movie_rating"`
}

type user_movie_comments struct {
	CommentID int    `json:"comment_id"`
	UserID    int    `json:"user_id"`
	MovieID   int    `json:"movie_id"`
	Comment   string `json:"movie_comments"`
}

type user_activity struct {
	UserID  int      `json:"user_id"`
	MovieID int      `json:"movie_id"`
	Rating  int      `json:"movie_rating"`
	Comment []string `json:"movie_comments"`
}

type movie_details struct {
	ID            int                   `json:"id"`
	Name          string                `json:"movie_name"`
	AverageRating float64               `json:"average_rating"`
	NumRatings    int                   `json:"number_of_ratings"`
	UserComments  []user_movie_comments `json:"user_comments"`
}

func (p *user_movie_comments) addUserCommentOnMovie(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO user_movie_comments(comment_id, movie_id, user_id, comment) VALUES($1, $2, $3, $4) RETURNING id",
		p.CommentID, p.MovieID, p.UserID, p.Comment).Scan(&p.CommentID)
	if err != nil {
		return err
	}
	return nil
}

func (p *user_movie_rating) addUserRatingOnMovie(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO user_movie_rating(movie_id, user_id, rating) VALUES($1, $2, $3) RETURNING movie_id",
		p.MovieID, p.UserID, p.Rating).Scan(&p.MovieID)
	if err != nil {
		return err
	}
	return nil
}

func searchMovie(db *sql.DB, name string) (movie_details, error) {
	var m movie_details
	err := db.QueryRow(
		"SELECT movie_id, movie_name FROM movies WHERE movie_name CONTAINS $1 LIMIT 1",
		name).Scan(&m.ID, &m.Name)

	if err != nil {
		return movie_details{}, err
	}

	user_movie_rating_rows, err := db.Query(
		"SELECT movie_id, rating FROM user_movie_rating WHERE movie_id=$1",
		m.ID)
	if err != nil {
		return movie_details{}, err
	}
	defer user_movie_rating_rows.Close()

	user_movie_ratings := []user_movie_rating{}

	for user_movie_rating_rows.Next() {
		var r user_movie_rating
		if err := user_movie_rating_rows.Scan(&r.MovieID, &r.Rating); err != nil {
			return movie_details{}, err
		}
		user_movie_ratings = append(user_movie_ratings, r)
	}

	user_movie_comment_rows, err := db.Query(
		"SELECT comment_id, comment FROM user_movie_comments WHERE movie_id=$1",
		m.ID)
	if err != nil {
		return movie_details{}, err
	}
	defer user_movie_comment_rows.Close()

	user_movie_comments_list := []user_movie_comments{}

	for user_movie_comment_rows.Next() {
		var r user_movie_comments
		if err := user_movie_comment_rows.Scan(&r.CommentID, &r.Comment); err != nil {
			return movie_details{}, err
		}
		user_movie_comments_list = append(user_movie_comments_list, r)
	}

	m.NumRatings = len(user_movie_ratings)
	m.AverageRating = 0

	for _, user_movie_rating := range user_movie_ratings {
		m.AverageRating += float64(user_movie_rating.Rating)
	}
	m.AverageRating /= float64(m.NumRatings)
	m.UserComments = user_movie_comments_list

	return m, nil
}

func getUserActivity(db *sql.DB, user_id int) ([]user_activity, error) {

	user_movie_rating_rows, err := db.Query(
		"SELECT user_id, movie_id, rating FROM user_movie_rating WHERE user_id=$1",
		user_id)
	if err != nil {
		return nil, err
	}
	defer user_movie_rating_rows.Close()

	user_activies := []user_activity{}

	for user_movie_rating_rows.Next() {
		var r user_activity
		if err := user_movie_rating_rows.Scan(&r.UserID, &r.MovieID, &r.Rating); err != nil {
			return nil, err
		}
		user_activies = append(user_activies, r)
	}

	user_movie_comment_rows, err := db.Query(
		"SELECT movie_id, comment_id, comment FROM user_movie_comments WHERE user_id=$1",
		user_id)
	if err != nil {
		return nil, err
	}
	defer user_movie_comment_rows.Close()

	user_movie_comments_list := []user_movie_comments{}

	for user_movie_comment_rows.Next() {
		var r user_movie_comments
		if err := user_movie_comment_rows.Scan(&r.MovieID, &r.CommentID, &r.Comment); err != nil {
			return nil, err
		}
		user_movie_comments_list = append(user_movie_comments_list, r)
	}

	for i := range user_activies {
		user_activies[i].Comment = []string{}
		for j := range user_movie_comments_list {
			if user_movie_comments_list[j].MovieID == user_activies[i].MovieID {
				user_activies[i].Comment = append(user_activies[i].Comment, user_movie_comments_list[j].Comment)
			}
		}
	}
	return user_activies, nil
}
