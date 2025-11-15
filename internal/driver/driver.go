package driver

import (
	"database/sql"
	"fmt"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// OpenDB opens a SQLite database with given DSN.
func OpenDB1(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Ping error:", err)
		return nil, err
	}

	return db, nil
}

// OpenDB opens a GORM SQLite database with the given DSN
func OpenDB(dsn string) (*gorm.DB, error) {
	// ✅ Disable GORM SQL logging to prevent log pollution
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDatabase(dsn string, logger logger.Logger) (*models.DBModel, error) {
	conn, err := OpenDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := models.AutoMigrate(conn); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate schema: %w", err)
	}

	logger.Info("✅ Local DB ready (auth, CID cache, sync log)")
	return &models.DBModel{DB: conn}, nil
}
