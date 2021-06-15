CREATE TABLE IF NOT EXISTS user (
    user_id INT PRIMARY KEY,
    user_name VARCHAR( 50 ) UNIQUE NOT NULL,
    user_role INT 
)

CREATE TABLE MOVIE (
    movie_id INT PRIMARY KEY,
    movie_name VARCHAR( 100 ) UNIQUE NOT NULL
)

CREATE TABLE user_movie_comments (
    comment_id INT PRIMARY KEY,
    user_id INT,
    movie_id INT, 
    comment VARCHAR ( 200 )
    FOREIGN KEY (user_id) REFERENCES user (user_id)
    FOREIGN KEY (movie_id) REFERENCES movie (movie_id)
)

CREATE TABLE user_movie_rating (
    user_id INT, 
    movie_id INT,
    rating INT,
    PRIMARY KEY (user_id, movie_id)
    FOREIGN KEY (user_id) REFERENCES user (user_id)
    FOREIGN KEY (movie_id) REFERENCES movie (movie_id)
    
)
