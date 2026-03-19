"""
ПЕНТЕСТ - Тестирование безопасности
Проверяет все механизмы защиты
"""

import requests
import threading
import time
import json
import random
import string
import hashlib
from concurrent.futures import ThreadPoolExecutor

BASE_URL = "http://localhost:8080"
API_URL = f"{BASE_URL}/api"

class SecurityTester:
    """Тестер безопасности"""
    
    def __init__(self):
        self.results = {
            'passed': [],
            'failed': [],
            'warnings': []
        }
    
    def run_all_tests(self):
        """Запуск всех тестов"""
        print("🔍 ЗАПУСК ПЕНТЕСТА")
        print("=" * 50)
        
        self.test_sql_injection()
        self.test_xss()
        self.test_rate_limiting()
        self.test_brute_force()
        self.test_jwt_security()
        self.test_2fa()
        self.test_ssrf()
        self.test_dos_protection()
        
        self.print_report()
    
    def test_sql_injection(self):
        """Тест SQL инъекций"""
        print("\n📊 Тест SQL Injection...")
        
        payloads = [
            "' OR '1'='1",
            "'; DROP TABLE users; --",
            "' UNION SELECT * FROM users--",
            "admin'--",
            "1' ORDER BY 1--",
            "' OR 1=1; --",
            "' UNION SELECT username,password FROM users--"
        ]
        
        for payload in payloads:
            try:
                response = requests.post(
                    f"{API_URL}/login",
                    json={"email": f"test{payload}@test.com", "password": "test"},
                    timeout=5
                )
                
                if response.status_code in [400, 403, 429]:
                    self.results['passed'].append(f"SQL Injection blocked: {payload}")
                else:
                    self.results['failed'].append(f"SQL Injection possible: {payload}")
            except:
                self.results['warnings'].append(f"Timeout for payload: {payload}")
    
    def test_xss(self):
        """Тест XSS атак"""
        print("\n📊 Тест XSS...")
        
        payloads = [
            "<script>alert(1)</script>",
            "<img src=x onerror=alert(1)>",
            "javascript:alert(1)",
            "<svg onload=alert(1)>",
            "';alert(1);//",
            "<body onload=alert(1)>"
        ]
        
        for payload in payloads:
            try:
                response = requests.post(
                    f"{API_URL}/register",
                    json={
                        "email": f"test@test.com",
                        "password": "Test123!",
                        "first_name": payload,
                        "last_name": "User",
                        "role": "patient"
                    },
                    timeout=5
                )
                
                if response.status_code in [400, 403]:
                    self.results['passed'].append(f"XSS blocked: {payload}")
                else:
                    self.results['failed'].append(f"XSS possible: {payload}")
            except:
                self.results['warnings'].append(f"Timeout for payload: {payload}")
    
    def test_rate_limiting(self):
        """Тест rate limiting"""
        print("\n📊 Тест Rate Limiting...")
        
        def make_request():
            try:
                return requests.get(f"{BASE_URL}/health", timeout=1)
            except:
                return None
        
        # 200 быстрых запросов
        with ThreadPoolExecutor(max_workers=20) as executor:
            responses = list(executor.map(lambda _: make_request(), range(200)))
        
        blocked = sum(1 for r in responses if r and r.status_code == 429)
        
        if blocked > 0:
            self.results['passed'].append(f"Rate limiting works: {blocked} requests blocked")
        else:
            self.results['failed'].append("Rate limiting not detected")
    
    def test_brute_force(self):
        """Тест защиты от брутфорса"""
        print("\n📊 Тест Brute Force Protection...")
        
        def try_login(password):
            try:
                return requests.post(
                    f"{API_URL}/login",
                    json={"email": "admin@test.com", "password": password},
                    timeout=1
                )
            except:
                return None
        
        # 20 попыток входа с разными паролями
        with ThreadPoolExecutor(max_workers=5) as executor:
            passwords = [f"pass{i}" for i in range(20)]
            responses = list(executor.map(try_login, passwords))
        
        blocked = sum(1 for r in responses if r and r.status_code == 429)
        
        if blocked > 0:
            self.results['passed'].append(f"Brute force protection works: {blocked} attempts blocked")
        else:
            self.results['failed'].append("Brute force protection not detected")
    
    def test_jwt_security(self):
        """Тест безопасности JWT"""
        print("\n📊 Тест JWT Security...")
        
        # Попытка использовать JWT с неправильным алгоритмом
        import jwt
        
        try:
            # Попытка создать JWT с none алгоритмом
            token = jwt.encode({"user": "admin"}, "", algorithm="none")
            
            response = requests.get(
                f"{API_URL}/profile",
                headers={"Authorization": f"Bearer {token}"},
                timeout=5
            )
            
            if response.status_code == 401:
                self.results['passed'].append("JWT none algorithm blocked")
            else:
                self.results['failed'].append("JWT none algorithm accepted")
        except:
            self.results['warnings'].append("JWT test error")
    
    def test_2fa(self):
        """Тест 2FA"""
        print("\n📊 Тест 2FA...")
        
        # Попытка входа без 2FA для пользователя с включенной 2FA
        response = requests.post(
            f"{API_URL}/login",
            json={"email": "admin@test.com", "password": "admin123"},
            timeout=5
        )
        
        if response.status_code == 200 and response.json().get('twofa_required'):
            self.results['passed'].append("2FA required as expected")
        else:
            self.results['failed'].append("2FA not working properly")
    
    def test_ssrf(self):
        """Тест SSRF защиты"""
        print("\n📊 Тест SSRF Protection...")
        
        internal_urls = [
            "http://localhost:5432",
            "http://127.0.0.1:6379",
            "http://169.254.169.254/latest/meta-data/",
            "http://metadata.google.internal/"
        ]
        
        for url in internal_urls:
            try:
                response = requests.get(
                    f"{API_URL}/proxy",
                    params={"url": url},
                    timeout=5
                )
                
                if response.status_code in [400, 403]:
                    self.results['passed'].append(f"SSRF blocked: {url}")
                else:
                    self.results['failed'].append(f"SSRF possible: {url}")
            except:
                pass
    
    def test_dos_protection(self):
        """Тест защиты от DoS"""
        print("\n📊 Тест DoS Protection...")
        
        def send_large_request():
            try:
                # Отправка большого JSON
                large_data = {"data": "A" * 1000000}
                return requests.post(
                    f"{API_URL}/login",
                    json=large_data,
                    timeout=1
                )
            except:
                return None
        
        # 10 больших запросов параллельно
        with ThreadPoolExecutor(max_workers=10) as executor:
            responses = list(executor.map(lambda _: send_large_request(), range(10)))
        
        blocked = sum(1 for r in responses if r and r.status_code == 413)
        
        if blocked > 0:
            self.results['passed'].append(f"DoS protection works: {blocked} large requests blocked")
        else:
            self.results['failed'].append("DoS protection not detected")
    
    def print_report(self):
        """Вывод отчета"""
        print("\n" + "=" * 50)
        print("📊 ОТЧЕТ ПЕНТЕСТА")
        print("=" * 50)
        
        print(f"\n✅ Пройдено: {len(self.results['passed'])}")
        for test in self.results['passed']:
            print(f"  ✓ {test}")
        
        print(f"\n❌ Провалено: {len(self.results['failed'])}")
        for test in self.results['failed']:
            print(f"  ✗ {test}")
        
        print(f"\n⚠️ Предупреждения: {len(self.results['warnings'])}")
        for test in self.results['warnings']:
            print(f"  ! {test}")
        
        print("\n" + "=" * 50)
        
        if len(self.results['failed']) == 0:
            print("✅ СИСТЕМА БЕЗОПАСНА!")
        else:
            print("⚠️ ОБНАРУЖЕНЫ УЯЗВИМОСТИ!")

if __name__ == "__main__":
    tester = SecurityTester()
    tester.run_all_tests()
