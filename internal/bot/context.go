package bot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/terag/kadok/internal/cache"
	"github.com/terag/kadok/internal/http"
)

type BotContext struct {
	Cache      cache.Cache
	HttpClient http.Client
	Voice      VoiceContext
}

func NewBotContext(session *discordgo.Session, cache cache.Cache, httpClient http.Client) BotContext {
	return BotContext{
		Cache:      cache,
		HttpClient: httpClient,
		Voice: VoiceContext{
			stop:    make(chan error),
			s:       session,
			playing: sync.RWMutex{},
		},
	}
}
