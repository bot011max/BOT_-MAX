"""
ABSOLUTE ARMOR - Военный уровень защиты
Главный файл безопасности, объединяющий все механизмы
"""

import os
import sys
import time
import json
import hashlib
import hmac
import base64
import secrets
import logging
import threading
from typing import Dict, Any, Tuple, Optional
from datetime import datetime, timedelta
from collections import defaultdict
import redis
import jwt
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2
import bcrypt
import pyotp
import qrcode
import io
import ipaddress
import re
import requests

class AbsoluteArmor:
    """
    Абсолютная броня - объединяет все механизмы защиты
    """
    
    def __init__(self, logger=None):
        self.logger = logger or logging.getLogger(__name__)
        self.redis = None
        self.waf_rules = self._load_waf_rules()
        self.blocked_ips = set()
        self.attack_stats = defaultdict(int)
        self.master_key = None
        self.jwt_secret = None
        self.cipher = None
        
    def Init(self) -> bool:
        """Инициализация всех компонентов безопасности"""
        try:
            # Подключение к Redis
            self.redis = redis.Redis(
                host=os.getenv('REDIS_HOST', 'localhost'),
                port=int(os.getenv('REDIS_PORT', 6379)),
                password=os.getenv('REDIS_PASSWORD', ''),
                decode_responses=True,
                socket_connect_timeout=5
            )
            self.redis.ping()
            self.logger.info("✅ Redis connected")
            
            # Загрузка мастер-ключа
            self._load_master_key()
            
            # Запуск фоновых задач
            self._start_background_tasks()
            
            self.logger.info("✅ AbsoluteArmor initialized")
            return True
            
        except Exception as e:
            self.logger.error(f"❌ Failed to initialize AbsoluteArmor: {e}")
            return False
    
    def _load_master_key(self):
        """Загрузка мастер-ключа из файла или переменной окружения"""
        key_path = '/run/secrets/master_key'
        if os.path.exists(key_path):
            with open(key_path, 'rb') as f:
                self.master_key = f.read().strip()
        else:
            self.master_key = os.getenv('MASTER_KEY', Fernet.generate_key()).encode()
        
        self.cipher = Fernet(base64.urlsafe_b64encode(self.master_key[:32]))
        self.jwt_secret = os.getenv('JWT_SECRET', secrets.token_urlsafe(64))
    
    def _load_waf_rules(self) -> Dict:
        """Загрузка правил WAF"""
        return {
            'sql_injection': [
                r"(\bSELECT\b.*\bFROM\b)",
                r"(\bINSERT\b.*\bINTO\b)",
                r"(\bUPDATE\b.*\bSET\b)",
                r"(\bDELETE\b.*\bFROM\b)",
                r"(\bDROP\b.*\bTABLE\b)",
                r"(\bUNION\b.*\bSELECT\b)",
                r"--",
                r";",
                r"(')\s*(OR|AND)\s*(')",
            ],
            'xss': [
                r"<script.*?>.*?</script>",
                r"javascript:",
                r"onerror\s*=",
                r"onload\s*=",
                r"alert\s*\(",
                r"prompt\s*\(",
                r"confirm\s*\(",
                r"document\.cookie",
            ],
            'path_traversal': [
                r"\.\.[\/\\]",
                r"\/etc\/passwd",
                r"\/windows\/win\.ini",
            ],
            'command_injection': [
                r";\s*(cat|rm|wget|curl|nc|bash|sh)",
                r"\|\s*(cat|rm|wget|curl|nc|bash|sh)",
                r"`.*`",
                r"\$\(.*\)",
            ]
        }
    
    def _start_background_tasks(self):
        """Запуск фоновых задач безопасности"""
        threading.Thread(target=self._cleanup_blocks, daemon=True).start()
        threading.Thread(target=self._sync_blacklist, daemon=True).start()
    
    def _cleanup_blocks(self):
        """Очистка устаревших блокировок"""
        while True:
            time.sleep(3600)  # Каждый час
            # Очистка через Redis с TTL
    
    def _sync_blacklist(self):
        """Синхронизация черного списка с Redis"""
        while True:
            time.sleep(300)  # Каждые 5 минут
            try:
                blacklisted = self.redis.smembers('global:blacklist')
                self.blocked_ips = set(blacklisted)
            except:
                pass
    
    def ProtectRequest(self):
        """Middleware для защиты запросов"""
        def middleware(c):
            # Получение IP
            ip = c.ClientIP()
            
            # Проверка в черном списке
            if ip in self.blocked_ips:
                self.logger.warning(f"Blocked request from blacklisted IP: {ip}")
                c.AbortWithStatusJSON(403, {"error": "Access denied", "code": "IP_BLOCKED"})
                return
            
            # Rate limiting
            if not self._check_rate_limit(ip):
                self.logger.warning(f"Rate limit exceeded for IP: {ip}")
                c.AbortWithStatusJSON(429, {"error": "Too many requests", "code": "RATE_LIMITED"})
                return
            
            # WAF проверка
            if not self._waf_check(c.Request):
                self.logger.warning(f"WAF blocked request from IP: {ip}")
                c.AbortWithStatusJSON(403, {"error": "Malicious request detected", "code": "WAF_BLOCK"})
                return
            
            c.Next()
        return middleware
    
    def _check_rate_limit(self, ip: str) -> bool:
        """Проверка лимитов запросов"""
        key = f"ratelimit:{ip}"
        current = self.redis.get(key)
        
        if current and int(current) > 100:  # 100 запросов в минуту
            return False
        
        pipe = self.redis.pipeline()
        pipe.incr(key)
        pipe.expire(key, 60)
        pipe.execute()
        return True
    
    def _waf_check(self, request) -> bool:
        """Web Application Firewall проверка"""
        # Проверка URL
        url = request.URL.String()
        for attack_type, patterns in self.waf_rules.items():
            for pattern in patterns:
                if re.search(pattern, url, re.IGNORECASE):
                    self.attack_stats[attack_type] += 1
                    return False
        
        # Проверка тела запроса для POST/PUT
        if request.Method in ['POST', 'PUT', 'PATCH']:
            body = request.GetBody()
            if body:
                for attack_type, patterns in self.waf_rules.items():
                    for pattern in patterns:
                        if re.search(pattern, body, re.IGNORECASE):
                            self.attack_stats[attack_type] += 1
                            return False
        
        return True
    
    def CreateJWT(self, user_id: str, role: str) -> str:
        """Создание JWT токена"""
        payload = {
            'user_id': user_id,
            'role': role,
            'iat': datetime.utcnow(),
            'exp': datetime.utcnow() + timedelta(hours=1),
            'jti': secrets.token_urlsafe(16)
        }
        return jwt.encode(payload, self.jwt_secret, algorithm='HS256')
    
    def VerifyJWT(self, token: str) -> Dict:
        """Проверка JWT токена"""
        try:
            return jwt.decode(token, self.jwt_secret, algorithms=['HS256'])
        except jwt.ExpiredSignatureError:
            raise Exception("Token expired")
        except jwt.InvalidTokenError:
            raise Exception("Invalid token")
    
    def ValidatePassword(self, password: str) -> bool:
        """Проверка сложности пароля"""
        if len(password) < 8:
            return False
        if not re.search(r'[A-Z]', password):
            return False
        if not re.search(r'[a-z]', password):
            return False
        if not re.search(r'[0-9]', password):
            return False
        if not re.search(r'[!@#$%^&*(),.?":{}|<>]', password):
            return False
        return True
    
    def GenerateTOTPSecret(self) -> str:
        """Генерация секрета для TOTP (Google Authenticator)"""
        return pyotp.random_base32()
    
    def VerifyTOTP(self, secret: str, code: str) -> bool:
        """Проверка TOTP кода"""
        totp = pyotp.TOTP(secret)
        return totp.verify(code)
    
    def Encrypt(self, data: bytes) -> bytes:
        """Шифрование данных"""
        return self.cipher.encrypt(data)
    
    def Decrypt(self, encrypted: bytes) -> bytes:
        """Дешифрование данных"""
        return self.cipher.decrypt(encrypted)
    
    def HashPassword(self, password: str) -> str:
        """Хеширование пароля"""
        salt = bcrypt.gensalt(rounds=12)
        return bcrypt.hashpw(password.encode(), salt).decode()
    
    def AuditLog(self, action: str, user_id: str, ip: str):
        """Запись в аудит лог"""
        event = {
            'timestamp': datetime.utcnow().isoformat(),
            'action': action,
            'user_id': user_id,
            'ip': ip,
            'event_id': secrets.token_urlsafe(16)
        }
        # Сохранение в Redis для последующей отправки в SIEM
        self.redis.lpush('audit:queue', json.dumps(event))
        self.logger.info(f"AUDIT: {action} - {user_id}")
    
    def MetricsHandler(self):
        """Handler для Prometheus метрик"""
        def handler(c):
            metrics = f"""
# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total {self.redis.get('metrics:requests') or 0}

# HELP waf_blocks_total Total WAF blocks
# TYPE waf_blocks_total counter
waf_blocks_total {self.attack_stats.get('sql_injection', 0) + self.attack_stats.get('xss', 0)}

# HELP rate_limiter_blocks_total Total rate limiter blocks
# TYPE rate_limiter_blocks_total counter
rate_limiter_blocks_total {self.redis.get('metrics:rate_blocks') or 0}

# HELP active_users Current active users
# TYPE active_users gauge
active_users {self.redis.scard('active:sessions')}
"""
            c.String(200, metrics)
        return handler
