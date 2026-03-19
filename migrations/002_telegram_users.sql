-- Дополнительные таблицы для Telegram бота

-- Таблица для состояний диалогов
CREATE TABLE IF NOT EXISTS telegram_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telegram_id BIGINT NOT NULL REFERENCES telegram_users(telegram_id) ON DELETE CASCADE,
    state TEXT DEFAULT 'none',
    temp_data JSONB,
    last_message_id INTEGER,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_telegram_sessions_telegram_id ON telegram_sessions(telegram_id);
CREATE INDEX idx_telegram_sessions_state ON telegram_sessions(state);

-- Таблица для истории сообщений
CREATE TABLE IF NOT EXISTS telegram_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telegram_id BIGINT NOT NULL REFERENCES telegram_users(telegram_id) ON DELETE CASCADE,
    message_id INTEGER NOT NULL,
    text TEXT,
    is_bot BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_telegram_messages_telegram_id ON telegram_messages(telegram_id);
CREATE INDEX idx_telegram_messages_created ON telegram_messages(created_at);

-- Функция для очистки старых сессий
CREATE OR REPLACE FUNCTION cleanup_old_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM sessions WHERE expires_at < NOW();
    DELETE FROM telegram_sessions WHERE updated_at < NOW() - INTERVAL '1 day';
END;
$$ LANGUAGE plpgsql;
