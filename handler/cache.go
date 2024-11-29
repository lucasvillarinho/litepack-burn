package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	lcache "github.com/lucasvillarinho/litepack/cache"
)

type CacheHandler struct {
	lcache lcache.Cache
}

func NewCacheHandler(ctx context.Context, lcache lcache.Cache) (*CacheHandler, error) {
	return &CacheHandler{
		lcache: lcache,
	}, nil
}

func (ch *CacheHandler) Set(c echo.Context) error {
	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	err := ch.lcache.Set(c.Request().Context(), payload.Key, payload.Value, time.Duration(payload.TTL)*time.Second)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set cache"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cache set successfully"})
}
