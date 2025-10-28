# Authentication Service

Микросервис аутентификации с поддержкой JWT токенов. Реализована регистрация пользователей, вход в систему и защита маршрутов через middleware. Пароли хешируются с использованием bcrypt. Приложение контейнеризировано с Docker и запускается через docker-compose вместе с PostgreSQL.

## 🚀 Особенности

- **JWT аутентификация** - безопасные токены с временем жизни
- **Хеширование паролей** - использование bcrypt для защиты паролей
- **Swagger документация** - полная API документация
- **Docker контейнеризация** - легкий деплоймент
- **PostgreSQL** - надежное хранение данных
- **Fiber фреймворк** - высокопроизводительный веб-фреймворк

## Swagger
<img width="2494" height="1100" alt="Screenshot 2025-10-28 at 21-04-09 Swagger UI" src="https://github.com/user-attachments/assets/eae22be6-d275-4055-bd16-3e5114440637" />

## Быстрый старт
1. **Клонирование репозитория**
```bash
git clone https://github.com/your-username/auth-service.git
cd auth-service
```
2. **Запуск докер файла**
```bash
docker-compose up -d
```
