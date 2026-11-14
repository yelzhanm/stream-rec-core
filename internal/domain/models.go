package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id" swaggertype:"string" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
	Email          string             `bson:"email" json:"email" example:"yelzhan@aitu.edu.kz"`
	PasswordHash   string             `bson:"password_hash" json:"-"`
	Password       string             `bson:"-" json:"password,omitempty" example:"secret123"`
	Name           string             `bson:"name" json:"name" example:"Yelzhan"`
	FavoriteGenres []string           `bson:"favorite_genres" json:"favorite_genres" example:"Sci-Fi,Drama"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}

type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" swaggertype:"string" example:"65f1a2b3c4d5e6f7a8b9c0d2"`
	Name        string             `bson:"name" json:"name" example:"Inception"`
	Description string             `bson:"description" json:"description" example:"A thief steals corporate secrets through dream-sharing technology."`
	Category    string             `bson:"category" json:"category" example:"Sci-Fi"`
	ReleaseYear int                `bson:"release_year" json:"release_year" example:"2010"`
	DurationMin int                `bson:"duration_min" json:"duration_min" example:"148"`
	Rating      float64            `bson:"rating" json:"rating" example:"8.8"`
}

type Interaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id" swaggertype:"string" example:"65f1a2b3c4d5e6f7a8b9c0d3"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id" swaggertype:"string" example:"65f1a2b3c4d5e6f7a8b9c0d1"`
	ItemID       primitive.ObjectID `bson:"item_id" json:"item_id" swaggertype:"string" example:"65f1a2b3c4d5e6f7a8b9c0d2"`
	Type         string             `bson:"type" json:"type" example:"rating"` // "view", "rating", "selection"
	RatingValue  int                `bson:"rating_value,omitempty" json:"rating_value,omitempty" example:"5"`
	WatchTimeSec int                `bson:"watch_time_sec,omitempty" json:"watch_time_sec,omitempty" example:"8800"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type ItemSearchFilter struct {
	Query    string
	Category string
	Year     int
}
