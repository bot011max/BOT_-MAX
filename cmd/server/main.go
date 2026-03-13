package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "BOT_-MAX/internal/api"
    "BOT_-MAX/internal/middleware"
    "BOT_-MAX/internal/models"
)

func main() {
    // Загружаем .env файл
    if err := godotenv.Load(); err != nil {
        log.Println("Файл .env не найден, используем переменные окружения")
    }

    // Подключение к БД
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_PORT", "5432"),
        getEnv("DB_USER", "postgres"),
        getEnv("DB_PASSWORD", "postgres"),
        getEnv("DB_NAME", "medical_bot"),
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Ошибка подключения к БД:", err)
    }
    
    // Автоматическая миграция
    if err := db.AutoMigrate(
        &models.User{},
        &models.Patient{},
        &models.Doctor{},
        &models.Prescription{},
        &models.Reminder{},
    ); err != nil {
        log.Fatal("Ошибка миграции:", err)
    }
    
    log.Println("База данных успешно подключена и мигрирована")

    // Создаем роутер
    r := gin.Default()
    
    // CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Статические файлы
    r.Static("/static", "./web/static")
    r.LoadHTMLGlob("web/templates/*")

    // Главная страница
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
            "title": "Медицинский бот",
        })
    })

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
            "time":   time.Now(),
        })
    })

    // Инициализация обработчиков
    authHandler := api.NewAuthHandler(db)

    // API маршруты
    api := r.Group("/api")
    {
        // Публичные маршруты
        api.POST("/register", authHandler.Register)
        api.POST("/login", authHandler.Login)
        
        // Защищенные маршруты
        protected := api.Group("/")
        protected.Use(middleware.AuthRequired())
        {
            protected.GET("/profile", authHandler.Profile)
        }
    }

    // Запуск сервера
    port := getEnv("PORT", "8080")
    srv := &http.Server{
        Addr:         ":" + port,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("Сервер запущен на порту %s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Ошибка сервера:", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Завершение работы...")

    // Закрываем соединение с БД
    sqlDB, _ := db.DB()
    sqlDB.Close()
    
    log.Println("Сервер остановлен")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
