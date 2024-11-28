package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/lucasvillarinho/litepack-burn/cache"
)

func main() {
	ctx := context.Background()
	handler, err := cache.NewCacheHandler(ctx)
	if err != nil {
		slog.Error("failed to create cache handler", slog.Any("error", err))
		return
	}

	http.HandleFunc("/set", handler.Set)
	http.ListenAndServe(":8080", nil)
}
