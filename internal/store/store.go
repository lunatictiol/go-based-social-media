package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
	ErrConflict          = errors.New("resource already exists")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostByID(ctx context.Context, id int64) (*Post, error)
		DeletePostByID(ctx context.Context, id int64) error
		UpdatePost(ctx context.Context, post *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetUserByID(ctx context.Context, id int64) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, userId, followerId int64) error
		UnFollow(ctx context.Context, userId, followerId int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowStore{db: db},
	}
}
