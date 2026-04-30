package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"server/internal/model"
)

const (
	UserTTL = 30 * time.Minute
)

func GetUser(ctx context.Context, userID uint) (*CachedUser, error) {
	key := fmt.Sprintf("user:%d", userID)
	data, err := Client().Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var u CachedUser
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func SetUser(ctx context.Context, userID uint, u *CachedUser) error {
	key := fmt.Sprintf("user:%d", userID)
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return Client().Set(ctx, key, data, UserTTL).Err()
}

func DeleteUser(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:%d", userID)
	return Client().Del(ctx, key).Err()
}

type CachedUser struct {
	ID             uint   `json:"id"`
	Username       string `json:"nickname"`
	Avatar         string `json:"avatar"`
	Bio            string `json:"bio"`
	Gender         uint   `json:"gender"`
	Birthday       string `json:"birthday"`
	FollowingCount uint   `json:"following_count"`
	FollowerCount  uint   `json:"follower_count"`
}

func (c *CachedUser) ToModel() *model.User {
	return &model.User{
		ID:             c.ID,
		Username:       c.Username,
		Avatar:         c.Avatar,
		Bio:            c.Bio,
		Gender:         c.Gender,
		Birthday:       c.Birthday,
		FollowingCount: c.FollowingCount,
		FollowerCount:  c.FollowerCount,
	}
}

func NewCachedUser(u *model.User) *CachedUser {
	return &CachedUser{
		ID:             u.ID,
		Username:       u.Username,
		Avatar:         u.Avatar,
		Bio:            u.Bio,
		Gender:         u.Gender,
		Birthday:       u.Birthday,
		FollowingCount: u.FollowingCount,
		FollowerCount:  u.FollowerCount,
	}
}
