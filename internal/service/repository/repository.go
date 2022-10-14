package repository

import (
	"context"

	"github.com/hetonei/arcanery-go-backend/config"
	"github.com/hetonei/arcanery-go-backend/internal/service"
)

func GetRepoConnection(ctx context.Context, cfg config.Config) service.ConnectionRepo {
	m := map[string]string{
		"name":     cfg.Database.Db,
		"uri":      cfg.Database.Url,
		"username": cfg.Database.Username,
		"password": cfg.Database.Password,
	}

	switch cfg.Database.Db {
	case "mongo":
		c := GetConnectionMongo(ctx, m)
		return c
	}

	return nil
}
