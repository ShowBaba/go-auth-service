package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/showbaba/go-auth-service/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToPgDB(host, user, password, dbname string, port int) (*gorm.DB, *sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, nil, err
	}

	log.Println("Database connection established!")
	return db, sqlDB, nil
}

func ConnectToRedis(address string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: utils.GetConfig().RedisAddr,
	})
}
