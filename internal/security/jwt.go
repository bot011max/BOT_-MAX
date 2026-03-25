package security

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("medical-bot-secret-key-2026-military-grade-security")

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// GenerateJWT генерирует JWT токен
func GenerateJWT(userID, email, role string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "medical-bot",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ValidateJWT проверяет JWT токен и возвращает claims
func ValidateJWT(tokenString string) (map[string]interface{}, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}

// GetUserIDFromToken извлекает user_id из токена
func GetUserIDFromToken(tokenString string) (string, error) {
    claims, err := ValidateJWT(tokenString)
    if err != nil {
        return "", err
    }
    
    userID, ok := claims["user_id"].(string)
    if !ok {
        return "", fmt.Errorf("user_id not found in token")
    }
    
    return userID, nil
}
