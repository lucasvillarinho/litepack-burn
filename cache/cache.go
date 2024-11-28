package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lucasvillarinho/litepack"
	"github.com/lucasvillarinho/litepack/cache"
)

type CacheHandler struct {
	lcache cache.Cache
}

func NewCacheHandler(ctx context.Context) (*CacheHandler, error) {
	lcache, err := litepack.NewCache(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	return &CacheHandler{
		lcache: lcache,
	}, nil
}

func (ch *CacheHandler) Set(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := ch.lcache.Set(r.Context(), payload.Key, payload.Value, time.Duration(payload.TTL)*time.Second)
	if err != nil {
		http.Error(w, "Failed to set cache", http.StatusInternalServerError)
		return
	}
}
