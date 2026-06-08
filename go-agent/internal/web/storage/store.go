package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"

	"nano_cc/internal/config"
)

type Store struct {
	DB    *sql.DB
	Redis *redis.Client
	ES    *elasticsearch.Client
	Cfg   config.AppConfig
}

var (
	//go:embed migrations/*.sql
	migrationFiles embed.FS
)

func NewStore(cfg config.AppConfig) (*Store, error) {
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.ESAddresses,
		Username:  cfg.ESUsername,
		Password:  cfg.ESPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("new elasticsearch client: %w", err)
	}

	return &Store{DB: db, Redis: rdb, ES: es, Cfg: cfg}, nil
}

func (s *Store) HealthCheck(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("mysql ping: %w", err)
	}
	if err := s.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	if s.ES != nil {
		res, err := s.ES.Ping(s.ES.Ping.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("elasticsearch ping: %w", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			return fmt.Errorf("elasticsearch ping: %s", res.String())
		}
	}
	return nil
}

func (s *Store) Close() error {
	if s.Redis != nil {
		_ = s.Redis.Close()
	}
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
