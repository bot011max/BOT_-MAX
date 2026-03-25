package database

import (
    "fmt"
    "log"
    "github.com/bot011max/medical-bot/internal/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewDB() (*gorm.DB, error) {
    // Используем SQLite для простоты
    dsn := "data/medical_bot.db"
    
    db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    log.Println("✅ Database connected successfully (SQLite)")
    return db, nil
}

func Migrate(db *gorm.DB) error {
    log.Println("🔄 Running database migrations...")
    
    err := db.AutoMigrate(
        &models.User{},
        &models.Medication{},
        &models.Reminder{},
    )
    if err != nil {
        return fmt.Errorf("failed to migrate database: %w", err)
    }
    
    log.Println("✅ Database migration completed")
    return nil
}
