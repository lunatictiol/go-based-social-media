package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follow struct {
	UserId     int64  `json:"user_id"`
	FollowerId int64  `json:"followee_id"`
	CreatedAt  string `json:"created_at"`
}
type FollowStore struct {
	db *sql.DB
}

func (f *FollowStore) Follow(ctx context.Context, userId, followerId int64) error {
	query := `
INSERT INTO followers (user_id,follower_id) values( $1,$2)
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := f.db.ExecContext(
		ctx,
		query,
		userId, followerId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err

}

func (f *FollowStore) UnFollow(ctx context.Context, userId, followerId int64) error {
	query := `
	DELETE FROM followers WHERE user_id = $1 AND follower_id=$2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := f.db.ExecContext(
		ctx,
		query,
		userId, followerId)

	return err
}
