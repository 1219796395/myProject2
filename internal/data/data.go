package data

import (
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	jsoniter "github.com/json-iterator/go"
	"github.com/patrickmn/go-cache"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// json serializer
var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRemoteConfigRepo, NewEnvManageRepo, NewRemoteConfigLogRepo, NewNetworkConfigRepo, NewNetworkConfigLogRepo, NewAdminUserRepo, NewAuthLogRepo)

// Data .
type Data struct {
	db    *gorm.DB
	rdb   *redis.Client
	cache *cache.Cache
}

// NewData .
func NewData(bc *conf.Bootstrap, logger log.Logger) (*Data, func(), error) {
	var c = bc.Data
	var log = log.NewHelper(logger)

	// mysql
	// TODO: gorm logger
	//TODO: log level
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{PrepareStmt: true})
	if err != nil {
		log.Errorf("failed opening connection to mysql %s", err)
		return nil, nil, err
	}
	// mysql config
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("failed opening connection to mysql %s", err)
		return nil, nil, err
	}
	sqlDB.SetMaxIdleConns(int(c.Database.MaxIdleConn))
	sqlDB.SetMaxOpenConns(int(c.Database.MaxOpenConn))
	sqlDB.SetConnMaxLifetime(c.Database.ConnMaxLifetime.AsDuration())
	sqlDB.SetConnMaxIdleTime(c.Database.ConnMaxIdleTime.AsDuration())
	// mysql otel
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		log.Errorf("failed use otelgorm %s", err)
		return nil, nil, err
	}

	// redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		PoolSize:     int(c.Redis.PoolSize),
	})
	// redis otel
	rdb.AddHook(redisotel.NewTracingHook())

	cache := cache.New(cache.NoExpiration, cache.NoExpiration)

	d := &Data{
		db:    db,
		rdb:   rdb,
		cache: cache,
	}

	cleanup := func() {
		log.Info("closing the data resources")

		// mysql
		sqlDB, err := d.db.DB()
		if err != nil {
			log.Error(err)
		}
		if err := sqlDB.Close(); err != nil {
			log.Error(err)
		}

		// redis
		if err = d.rdb.Close(); err != nil {
			log.Error(err)
		}
	}

	return d, cleanup, nil
}
