package database

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func NewConfig() *Config {
    return &Config{
        Host:     getEnv("DB_HOST", "postgres"),
        Port:     getEnv("DB_PORT", "5432"),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", "postgres"),
        DBName:   getEnv("DB_NAME", "medical_bot"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
    }
}

func Connect() error {
    config := NewConfig()
    
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
        config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
    )
    
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time { return time.Now().UTC() },
        SkipDefaultTransaction: true,
        PrepareStmt: true,
    })
    
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    log.Println("✅ Database connected successfully")
    return nil
}

func Close() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
