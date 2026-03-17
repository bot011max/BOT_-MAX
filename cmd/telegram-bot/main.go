// cmd/telegram/main.go - ПОЛНОСТЬЮ ПЕРЕРАБОТАННЫЙ
package main

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "sync"
    "syscall"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/bot011max/BOT_MAX/internal/security"
    "github.com/bot011max/BOT_MAX/internal/monitoring"
    "golang.org/x/time/rate"
)

type SecureTelegramBot struct {
    bot           *tgbotapi.BotAPI
    redis         *redis.Client
    secretManager *security.SecretManager
    auditLogger   *security.AuditLogger
    rateLimiter   *security.AdaptiveRateLimiter
    waf           *security.WAFMiddleware
    commands      map[string]CommandHandler
    limiter       *rate.Limiter
    mu            sync.RWMutex
    webhookSecret string
}

type CommandHandler struct {
    Handler     func(update tgbotapi.Update, args []string)
    Description string
    AuthRequired bool
    RateLimit    int
}

type WebhookRequest struct {
    UpdateID int `json:"update_id"`
    Message  *struct {
        MessageID int `json:"message_id"`
        From      struct {
            ID           int64  `json:"id"`
            IsBot        bool   `json:"is_bot"`
            FirstName    string `json:"first_name"`
            Username     string `json:"username"`
        } `json:"from"`
        Chat struct {
            ID        int64  `json:"id"`
            Type      string `json:"type"`
        } `json:"chat"`
        Date int    `json:"date"`
        Text string `json:"text"`
        Voice *struct {
            FileID   string `json:"file_id"`
            Duration int    `json:"duration"`
        } `json:"voice"`
        Photo []struct {
            FileID   string `json:"file_id"`
            FileSize int    `json:"file_size"`
        } `json:"photo"`
    } `json:"message"`
}

func NewSecureTelegramBot() (*SecureTelegramBot, error) {
    // Инициализация секретов
    secretManager, err := security.NewSecretManager(true)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания secret manager: %w", err)
    }

    // Получаем токен
    token, err := secretManager.GetSecret("telegram_token")
    if err != nil {
        return nil, fmt.Errorf("ошибка получения токена: %w", err)
    }

    // Создаем бота
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания бота: %w", err)
    }

    // Redis для хранения состояний
    redisPass, _ := secretManager.GetSecret("redis_password")
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":6379",
        Password: redisPass,
        DB:       0,
        PoolSize: 100,
    })

    // Rate limiter
    rateLimiter, err := security.NewAdaptiveRateLimiter("redis:6379", 
        &security.RateLimiterConfig{
            RequestsPerSecond: 1,
            Burst:             3,
            BlockDuration:     30 * time.Minute,
            CleanupInterval:   time.Minute,
            EnableAdaptive:    true,
        })
    if err != nil {
        return nil, err
    }

    // WAF
    wafConfig := security.WAFConfig{
        EnableSQLInjection: true,
        EnableXSS:         true,
    }
    waf, _ := security.NewWAFMiddleware(wafConfig)

    // Audit logger
    auditLogger, _ := security.NewAuditLogger()

    stb := &SecureTelegramBot{
        bot:           bot,
        secretManager: secretManager,
        auditLogger:   auditLogger,
        rateLimiter:   rateLimiter,
        waf:           waf,
        redis:         rdb,
        limiter:       rate.NewLimiter(rate.Limit(10), 20),
        commands:      make(map[string]CommandHandler),
        webhookSecret: security.GenerateRandomString(32),
    }

    stb.registerCommands()
    return stb, nil
}

// Установка вебхука с проверкой подписи
func (b *SecureTelegramBot) SetWebhook() error {
    webhookURL := os.Getenv("WEBHOOK_URL") + "/webhook/telegram"
    
    // Добавляем секрет в URL для валидации
    webhookURLWithSecret := fmt.Sprintf("%s?secret=%s", webhookURL, b.webhookSecret)
    
    webhookConfig := tgbotapi.NewWebhook(webhookURLWithSecret)
    
    // Устанавливаем максимальное количество соединений
    webhookConfig.MaxConnections = 40
    
    // Устанавливаем список разрешенных обновлений
    webhookConfig.AllowedUpdates = []string{
        "message",
        "callback_query",
        "inline_query",
    }

    _, err := b.bot.Request(webhookConfig)
    if err != nil {
        return fmt.Errorf("ошибка установки webhook: %w", err)
    }

    // Проверяем статус webhook
    webhookInfo, err := b.bot.GetWebhookInfo()
    if err != nil {
        return err
    }

    if webhookInfo.LastErrorDate != 0 {
        log.Printf("Ошибка webhook: %s", webhookInfo.LastErrorMessage)
    }

    log.Printf("✅ Webhook установлен: %s", webhookURL)
    log.Printf("📊 Статус: ожидает %d обновлений", webhookInfo.PendingUpdateCount)
    
    return nil
}

// Валидация входящего вебхука
func (b *SecureTelegramBot) validateWebhook(r *http.Request) bool {
    // Проверка secret из URL
    secret := r.URL.Query().Get("secret")
    if secret != b.webhookSecret {
        security.SecurityAlert("INVALID_WEBHOOK_SECRET", map[string]interface{}{
            "remote_addr": r.RemoteAddr,
            "user_agent":  r.UserAgent(),
        })
        return false
    }

    // Проверка подписи Telegram (X-Telegram-Bot-Api-Secret-Token)
    signature := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
    if signature != "" {
        expected := b.calculateSignature(r)
        if !hmac.Equal([]byte(signature), []byte(expected)) {
            security.SecurityAlert("INVALID_WEBHOOK_SIGNATURE", nil)
            return false
        }
    }

    return true
}

func (b *SecureTelegramBot) calculateSignature(r *http.Request) string {
    body, _ := io.ReadAll(r.Body)
    r.Body = io.NopCloser(bytes.NewBuffer(body))
    
    h := hmac.New(sha256.New, []byte(b.webhookSecret))
    h.Write(body)
    return hex.EncodeToString(h.Sum(nil))
}

// HTTP handler для вебхука
func (b *SecureTelegramBot) WebhookHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Валидация
    if !b.validateWebhook(r) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // 2. Rate limiting по IP
    if !b.rateLimiter.Allow(r.RemoteAddr) {
        http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
        return
    }

    // 3. Декодирование
    var update WebhookRequest
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    // 4. Асинхронная обработка
    go b.processUpdate(&update)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (b *SecureTelegramBot) processUpdate(update *WebhookRequest) {
    // Трассировка
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Логирование
    b.auditLogger.Log("TELEGRAM_UPDATE", "process", "update", "START", 
        fmt.Sprintf("%d", update.UpdateID), map[string]interface{}{
            "chat_id": update.Message.Chat.ID,
            "type":    b.getUpdateType(update),
        })

    // Валидация через WAF
    if err := b.waf.ValidateUpdate(update); err != nil {
        security.SecurityAlert("WAF_BLOCKED", map[string]interface{}{
            "update_id": update.UpdateID,
            "reason":    err.Error(),
        })
        return
    }

    // Обработка по типу
    if update.Message != nil {
        b.handleMessage(ctx, update)
    }
}

func (b *SecureTelegramBot) handleMessage(ctx context.Context, update *WebhookRequest) {
    chatID := update.Message.Chat.ID
    text := update.Message.Text

    // Санитизация
    text = b.sanitizeInput(text)

    // Rate limiting по пользователю
    if !b.rateLimiter.Allow(fmt.Sprintf("user:%d", chatID)) {
        b.sendMessage(chatID, "⏳ Слишком много запросов. Подождите немного.")
        return
    }

    // Проверка на команды
    if strings.HasPrefix(text, "/") {
        parts := strings.Fields(text)
        command := parts[0]
        args := parts[1:]

        b.mu.RLock()
        handler, exists := b.commands[command]
        b.mu.RUnlock()

        if exists {
            // Проверка авторизации
            if handler.AuthRequired && !b.isAuthorized(chatID) {
                b.sendMessage(chatID, "🔐 Для этой команды нужно привязать аккаунт.\nИспользуйте /bind")
                return
            }

            // Вызов хендлера
            handler.Handler(*update, args)
            
            // Метрики
            monitoring.TelegramMessagesTotal.WithLabelValues("command", command).Inc()
        } else {
            b.sendMessage(chatID, "❌ Неизвестная команда. Напишите /help")
        }
        return
    }

    // Обработка голосовых сообщений
    if update.Message.Voice != nil {
        b.handleVoice(ctx, update)
        return
    }

    // Обработка фото
    if len(update.Message.Photo) > 0 {
        b.handlePhoto(ctx, update)
        return
    }
}

// Команда /start
func (b *SecureTelegramBot) handleStart(update tgbotapi.Update, args []string) {
    chatID := update.Message.Chat.ID
    
    msg := "👋 <b>Добро пожаловать в Медицинского бота!</b>\n\n"
    msg += "🔐 <b>Безопасность:</b>\n"
    msg += "• Все сообщения шифруются\n"
    msg += "• Данные хранятся в зашифрованном виде\n"
    msg += "• Двухфакторная аутентификация\n\n"
    msg += "📱 <b>Возможности:</b>\n"
    msg += "• 💊 Напоминания о лекарствах\n"
    msg += "• 📝 Голосовой ввод симптомов\n"
    msg += "• 📅 Запись к врачу\n"
    msg += "• 📸 Анализ фото рецептов\n\n"
    msg += "🔑 <b>Начните с команды /bind</b> для привязки аккаунта"

    b.sendMessage(chatID, msg)
}

// Команда /bind с временным кодом
func (b *SecureTelegramBot) handleBind(update tgbotapi.Update, args []string) {
    chatID := update.Message.Chat.ID
    
    // Генерируем криптостойкий код
    code := security.GenerateNumericCode(6)
    
    // Сохраняем в Redis с TTL 10 минут
    ctx := context.Background()
    key := fmt.Sprintf("bind:%d", chatID)
    
    // Хешируем код перед сохранением
    hashedCode := b.hashCode(code)
    
    pipe := b.redis.Pipeline()
    pipe.Set(ctx, key, hashedCode, 10*time.Minute)
    pipe.Set(ctx, key+":attempts", 0, 10*time.Minute)
    _, err := pipe.Exec(ctx)
    
    if err != nil {
        b.sendMessage(chatID, "❌ Ошибка генерации кода. Попробуйте позже.")
        return
    }

    msg := "🔐 <b>Привязка Telegram к аккаунту</b>\n\n"
    msg += "Ваш одноразовый код:\n"
    msg += fmt.Sprintf("<code>%s</code>\n\n", code)
    msg += "Введите этот код в личном кабинете на сайте.\n"
    msg += "⏰ Код действителен 10 минут.\n\n"
    msg += "⚠️ Никому не сообщайте этот код!"

    // Inline клавиатура
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonURL("🔑 Перейти на сайт", os.Getenv("SITE_URL")),
        ),
    )

    b.sendMessageWithKeyboard(chatID, msg, keyboard)
}

// Команда /medications с пагинацией
func (b *SecureTelegramBot) handleMedications(update tgbotapi.Update, args []string) {
    chatID := update.Message.Chat.ID
    
    // Получаем user_id
    userID, err := b.getUserID(chatID)
    if err != nil {
        b.sendMessage(chatID, "❌ Ошибка получения данных")
        return
    }

    // Параметры пагинации
    page := 1
    if len(args) > 0 {
        fmt.Sscanf(args[0], "%d", &page)
    }
    limit := 5
    offset := (page - 1) * limit

    // Получаем из БД с пагинацией
    var prescriptions []models.Prescription
    var total int64
    
    b.db.Where("user_id = ? AND is_active = ?", userID, true).
        Count(&total)
    
    b.db.Where("user_id = ? AND is_active = ?", userID, true).
        Limit(limit).Offset(offset).
        Find(&prescriptions)

    if len(prescriptions) == 0 {
        b.sendMessage(chatID, "💊 У вас нет активных назначений")
        return
    }

    msg := fmt.Sprintf("💊 <b>Ваши лекарства (страница %d/%d):</b>\n\n", 
        page, (total+limit-1)/limit)
    
    for i, p := range prescriptions {
        msg += fmt.Sprintf("%d. <b>%s</b>\n", offset+i+1, p.Name)
        msg += fmt.Sprintf("   💊 Доза: %s\n", p.Dosage)
        msg += fmt.Sprintf("   ⏰ Режим: %s\n", p.Frequency)
        msg += fmt.Sprintf("   📅 До: %s\n", p.EndDate.Format("02.01.2006"))
        msg += "---\n"
    }

    // Клавиатура для навигации
    keyboard := b.createPaginationKeyboard(page, (total+limit-1)/limit)
    
    b.sendMessageWithKeyboard(chatID, msg, keyboard)
}

// Обработка голосовых сообщений
func (b *SecureTelegramBot) handleVoice(ctx context.Context, update *WebhookRequest) {
    chatID := update.Message.Chat.ID
    voice := update.Message.Voice

    // Проверка размера
    if voice.Duration > 120 { // 2 минуты макс
        b.sendMessage(chatID, "❌ Слишком длинное сообщение (макс 2 минуты)")
        return
    }

    b.sendMessage(chatID, "🎤 <i>Обрабатываю голосовое сообщение...</i>")

    // Получаем файл
    fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", 
        b.bot.Token, voice.FileID)

    // Скачиваем с проверкой
    audioData, err := b.downloadWithCheck(fileURL)
    if err != nil {
        b.sendMessage(chatID, "❌ Ошибка загрузки файла")
        return
    }

    // Отправляем в voice service через защищенный канал
    result, err := b.sendToVoiceService(audioData)
    if err != nil {
        b.sendMessage(chatID, "❌ Ошибка распознавания")
        return
    }

    b.sendMessage(chatID, fmt.Sprintf("📝 <b>Распознано:</b>\n%s", result))
}

// Graceful shutdown
func (b *SecureTelegramBot) Shutdown() {
    log.Println("🔄 Завершение работы бота...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Удаляем вебхук
    b.bot.RemoveWebhook()
    
    // Закрываем соединения
    b.redis.Close()
    
    log.Println("✅ Бот остановлен")
}

func main() {
    // Инициализация
    bot, err := NewSecureTelegramBot()
    if err != nil {
        log.Fatalf("❌ Ошибка создания бота: %v", err)
    }

    // Установка вебхука
    if err := bot.SetWebhook(); err != nil {
        log.Fatalf("❌ Ошибка установки webhook: %v", err)
    }

    // HTTP сервер для вебхуков
    http.HandleFunc("/webhook/telegram", bot.WebhookHandler)
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })

    srv := &http.Server{
        Addr:         ":8081",
        Handler:      nil,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // Graceful shutdown
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Ошибка сервера: %v", err)
        }
    }()

    log.Println("✅ Telegram bot слушает на :8081")

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    bot.Shutdown()
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
}
