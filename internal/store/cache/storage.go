package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (user *store.User, err error)
		Set(context.Context, *store.User) error
		Delete(context.Context, int64)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
	}
}
