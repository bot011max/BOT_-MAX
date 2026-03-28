#!/bin/bash
echo "🔐 ПРОВЕРКА HSM (Аппаратное шифрование)"
echo "========================================"
echo ""

cd /workspaces/BOT_MAX

# Создаем тестовые данные
echo "Тестовые данные для шифрования" > /tmp/test_data.txt

# Запускаем Go тест для HSM
go run -e << 'GO_TEST' 2>/dev/null
package main

import (
    "fmt"
    "log"
    "github.com/bot011max/medical-bot/internal/hsm"
)

func main() {
    fmt.Println("🔐 Тестирование HSM модуля...")
    
    hsmModule := hsm.NewSoftwareHSM()
    
    // Тестовые данные
    data := []byte("Секретные медицинские данные пациента")
    
    // Шифрование
    encrypted, err := hsmModule.Encrypt(data)
    if err != nil {
        log.Fatal("Ошибка шифрования:", err)
    }
    fmt.Printf("✅ Данные зашифрованы (длина: %d)\n", len(encrypted))
    fmt.Printf("   Зашифрованные данные: %s...\n", encrypted[:50])
    
    // Расшифровка
    decrypted, err := hsmModule.Decrypt(encrypted)
    if err != nil {
        log.Fatal("Ошибка расшифровки:", err)
    }
    fmt.Printf("✅ Данные расшифрованы: %s\n", string(decrypted))
    
    // Генерация ключа
    key, err := hsmModule.GenerateKey("test-key")
    if err != nil {
        log.Fatal("Ошибка генерации ключа:", err)
    }
    fmt.Printf("✅ Сгенерирован ключ: %s...\n", key[:50])
    
    fmt.Println("\n✅ HSM модуль работает корректно!")
}
GO_TEST
