package main

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lucasvillarinho/litepack"

	"github.com/lucasvillarinho/litepack-burn/handler"
)

func main() {
	ctx := context.Background()

	lcache, err := litepack.NewCache(ctx)
	if err != nil {
		slog.Error("failed to create cache", slog.Any("error", err))
		return
	}
	defer lcache.Destroy(ctx)

	handler, err := handler.NewCacheHandler(ctx, lcache)
	if err != nil {
		slog.Error("failed to create cache handler", slog.Any("error", err))
		return
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("cache/set", handler.Set)
	e.GET("cache/get/:key", handler.Get)
	e.DELETE("cache/delete/:key", handler.Delete)

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
	}
}
