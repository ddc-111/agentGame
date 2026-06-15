package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver   string `yaml:"driver"`   // sqlite 或 mysql
	DSN      string `yaml:"dsn"`      // 数据库连接字符串
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Database struct {
	DB *gorm.DB
}

func New(cfg Config) (*Database, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "sqlite":
		dsn := cfg.DSN
		if dsn == "" {
			dsn = "game.db"
		}
		dialector = sqlite.Open(dsn)

	case "mysql":
		dsn := cfg.DSN
		if dsn == "" {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		}
		dialector = mysql.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Printf("Connected to %s database", cfg.Driver)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}
