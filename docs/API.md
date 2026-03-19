# 📚 API Документация

## Базовый URL


## Аутентификация
Все защищенные эндпоинты требуют JWT токен в заголовке:


## Эндпоинты

### Регистрация
**POST** `/register`

**Тело запроса:**
```json
{
    "email": "user@example.com",
    "password": "Password123!",
    "first_name": "Иван",
    "last_name": "Петров",
    "role": "patient"
}

{
    "success": true,
    "message": "Registration successful",
    "twofa_secret": "JBSWY3DPEHPK3PXP"
}
