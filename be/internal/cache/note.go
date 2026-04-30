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
	NoteTTL   = 10 * time.Minute
	FeedTTL   = 2 * time.Minute
)

func GetNote(ctx context.Context, noteID uint) (*model.Note, error) {
	key := fmt.Sprintf("note:%d", noteID)
	data, err := Client().Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var n model.Note
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func SetNote(ctx context.Context, noteID uint, n *model.Note) error {
	key := fmt.Sprintf("note:%d", noteID)
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return Client().Set(ctx, key, data, NoteTTL).Err()
}

func DeleteNote(ctx context.Context, noteID uint) error {
	key := fmt.Sprintf("note:%d", noteID)
	return Client().Del(ctx, key).Err()
}

type CachedNoteList struct {
	List       []model.Note `json:"list"`
	Pagination Pagination   `json:"pagination"`
}

type Pagination struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	HasMore  bool  `json:"has_more"`
}

func GetFeed(ctx context.Context, tab string, page int) (*CachedNoteList, error) {
	key := fmt.Sprintf("feed:%s:%d", tab, page)
	data, err := Client().Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result CachedNoteList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func SetFeed(ctx context.Context, tab string, page int, result *CachedNoteList) error {
	key := fmt.Sprintf("feed:%s:%d", tab, page)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	ttl := FeedTTL
	if tab == "follow" {
		ttl = 1 * time.Minute
	}
	return Client().Set(ctx, key, data, ttl).Err()
}

func InvalidateFeed(ctx context.Context) error {
	keys, err := Client().Keys(ctx, "feed:*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return Client().Del(ctx, keys...).Err()
	}
	return nil
}

func InvalidateFeedByTab(ctx context.Context, tab string) error {
	keys, err := Client().Keys(ctx, fmt.Sprintf("feed:%s:*", tab)).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return Client().Del(ctx, keys...).Err()
	}
	return nil
}
