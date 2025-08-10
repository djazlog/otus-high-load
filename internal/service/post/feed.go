package post

import (
	"context"
	"fmt"
	"otus-project/internal/metric"
	"otus-project/internal/model"
)

// Feed Получить ленту постов
func (s *serv) Feed(ctx context.Context, id string, offset *float32, limit *float32) ([]*model.Post, error) {
	// Смотрим посты в редисе
	posts, err := s.postRRepository.Feed(ctx, id, offset, limit)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}

	if len(posts) > 0 {
		fmt.Println("posts from redis")
		metric.IncResponseCounter("GetPostFeedService", "redis")
		return posts, nil
	}

	// Если посты не найдены, то ищем в базе
	posts, err = s.postPgRepository.Feed(ctx, id, offset, limit)
	if err != nil {
		return nil, err
	}

	metric.IncResponseCounter("GetPostFeedService", "db")
	fmt.Println("posts from db")

	// СОхраняем посты в редис TODO: не сохраняет
	err = s.postRRepository.CacheFeed(ctx, id, posts)
	if err != nil {
		fmt.Println("posts not saved in redis", err)
		//return nil, err
	}

	return posts, nil
}
