
### **Makefile**
```makefile
.PHONY: help init up down logs clean backup restore test

help:
	@echo "Доступные команды:"
	@echo "  make init    - инициализация безопасности (создание ключей)"
	@echo "  make up      - запуск всех сервисов"
	@echo "  make down    - остановка всех сервисов"
	@echo "  make logs    - просмотр логов"
	@echo "  make clean   - очистка (удаление контейнеров и томов)"
	@echo "  make backup  - создание резервной копии"
	@echo "  make restore - восстановление из бэкапа"
	@echo "  make test    - запуск тестов"

init:
	@chmod +x scripts/init-security.sh
	@./scripts/init-security.sh

up:
	@docker-compose -f deployments/docker-compose.yml --env-file .env.production up -d
	@echo "✅ Сервисы запущены"

down:
	@docker-compose -f deployments/docker-compose.yml down
	@echo "✅ Сервисы остановлены"

logs:
	@docker-compose -f deployments/docker-compose.yml logs -f

clean:
	@docker-compose -f deployments/docker-compose.yml down -v
	@docker system prune -f
	@echo "✅ Очистка завершена"

backup:
	@./scripts/backup.sh

restore:
	@./scripts/restore.sh

test:
	@go test ./... -v
	@cd tests/security && python3 penetration_test.py
