package bot

import (
	"github.com/terag/kadok/internal/cache"
	"github.com/terag/kadok/internal/http"
)

type BotContext struct {
	Cache      cache.Cache
	HttpClient http.Client
}

func NewBotContext(cache cache.Cache, httpClient http.Client) BotContext {
	return BotContext{
		Cache:      cache,
		HttpClient: httpClient,
	}
}
