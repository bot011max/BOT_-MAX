package validator

import (
    "regexp"
)

func IsValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return re.MatchString(email)
}

func IsValidPhone(phone string) bool {
    re := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
    return re.MatchString(phone)
}
