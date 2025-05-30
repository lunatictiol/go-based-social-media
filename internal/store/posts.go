package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/lib/pq"
)

type Post struct {
	Id        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"userId"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
	User      User      `json:"user"`
}

type PostMetaData struct {
	Post         `json:"post"`
	CommentCount int `json:"comment_count"`
}
type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title,user_id,tags) 
	VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&post.Id,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	var post Post
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		SELECT id, user_id, title, content, created_at,  updated_at, tags, version
		FROM posts
		WHERE id = $1
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) DeletePostByID(ctx context.Context, id int64) error {

	query := ` DELETE FROM posts WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	res, err := s.db.ExecContext(
		ctx,
		query,
		id)

	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) UpdatePost(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING version
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Id,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
func (S *PostStore) GetUserFeed(ctx context.Context, id int64, fq PaginatedFeedQuery) ([]PostMetaData, error) {
	fmt.Println("Running with: id =", id, "limit =", fq.Limit, "offset =", fq.Offset, "search =", fq.Search, "tags =", fq.Tags)
	fmt.Println(reflect.TypeOf(fq.Tags))
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.user_id = p.user_id AND f.follower_id = $1
		WHERE 
			(f.follower_id = $1 OR p.user_id = $1) AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			($5::varchar[] IS NULL OR array_length($5::varchar[], 1) = 0 OR p.tags @> $5::varchar[])
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3
	`

	rows, err := S.db.QueryContext(ctx, query, id, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostMetaData
	for rows.Next() {
		var p PostMetaData
		err := rows.Scan(
			&p.Post.Id, &p.Post.UserId, &p.Post.Title, &p.Post.Content, &p.Post.CreatedAt, &p.Post.Version, pq.Array(&p.Post.Tags),
			&p.Post.User.Username, &p.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}
	return feed, nil

}
