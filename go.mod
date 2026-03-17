# 1. Исправляем go.mod - меняем версию Go на 1.22
cat > go.mod << 'EOF'
module github.com/bot011max/BOT_MAX

go 1.22

require (
    github.com/gin-gonic/gin v1.12.0
    github.com/prometheus/client_golang v1.23.2
    github.com/golang-jwt/jwt/v5 v5.3.1
    github.com/joho/godotenv v1.5.1
    gorm.io/driver/postgres v1.6.0
    gorm.io/gorm v1.31.1
    github.com/gin-contrib/cors v1.7.6
    github.com/go-redis/redis/v8 v8.11.5
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    golang.org/x/time v0.15.0
    github.com/google/uuid v1.6.0
    github.com/aws/aws-sdk-go v1.55.8
    github.com/hashicorp/golang-lru v1.0.2
    github.com/mssola/user_agent v0.6.0
    github.com/segmentio/kafka-go v0.4.50
)

require (
    github.com/beorn7/perks v1.0.1 // indirect
    github.com/bytedance/sonic v1.15.0 // indirect
    github.com/cespare/xxhash/v2 v2.3.0 // indirect
    github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
    github.com/jackc/pgpassfile v1.0.0 // indirect
    github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
    github.com/jackc/pgx/v5 v5.6.0 // indirect
    github.com/jinzhu/inflection v1.0.0 // indirect
    github.com/jinzhu/now v1.1.5 // indirect
    github.com/jmespath/go-jmespath v0.4.0 // indirect
    github.com/klauspost/compress v1.18.0 // indirect
    github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
    github.com/pierrec/lz4/v4 v4.1.15 // indirect
    github.com/prometheus/client_model v0.6.2 // indirect
    github.com/prometheus/common v0.66.1 // indirect
    github.com/prometheus/procfs v0.16.1 // indirect
    github.com/segmentio/kafka-go v0.4.50 // indirect
    github.com/xdg-go/pbkdf2 v1.0.0 // indirect
    github.com/xdg-go/scram v1.2.0 // indirect
    github.com/xdg-go/stringprep v1.0.4 // indirect
    golang.org/x/crypto v0.48.0 // indirect
    golang.org/x/sync v0.19.0 // indirect
    golang.org/x/sys v0.41.0 // indirect
    golang.org/x/text v0.34.0 // indirect
    google.golang.org/protobuf v1.36.10 // indirect
)
EOF

# 2. Удаляем старый go.sum
rm go.sum

# 3. Обновляем зависимости
go mod tidy

# 4. Проверяем, что создался go.sum
ls -la go.sum
