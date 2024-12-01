package handler

import (
	"context"
	"log/slog"
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
		slog.Error("failed to bind request payload", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	err := ch.lcache.Set(c.Request().Context(), payload.Key, payload.Value, time.Duration(payload.TTL)*time.Second)
	if err != nil {
		slog.Error("failed to set cache", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set cache"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cache set successfully"})
}

func (ch *CacheHandler) Get(c echo.Context) error {
	key := c.Param("key")

	value, err := ch.lcache.Get(c.Request().Context(), key)
	if err != nil {
		slog.Error("failed to get cache", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get cache"})
	}

	return c.JSON(http.StatusOK, map[string]string{"value": value})
}

func (ch *CacheHandler) Delete(c echo.Context) error {
	key := c.Param("key")

	err := ch.lcache.Del(c.Request().Context(), key)
	if err != nil {
		slog.Error("failed to delete cache", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete cache"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cache deleted successfully"})
}
