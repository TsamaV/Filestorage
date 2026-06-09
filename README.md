# 🗂 File Storage API (Go)

Сервис для хранения файлов с JWT-аутентификацией.  
Поддерживает регистрацию пользователей, загрузку изображений и управление файлами через S3-совместимое хранилище MinIO.

---

## ⚡ Системные требования

- **Go:** версия 1.21 или выше
- **PostgreSQL:** версия 16 или выше
- **Docker & Docker Compose:** для запуска инфраструктуры
- **ОС:** Windows, macOS, Linux
- **Память:** минимум 64 MB RAM

---

## 📌 Возможности

- 🔐 Регистрация и авторизация пользователей через JWT
- 📤 Загрузка изображений (PNG, JPEG) в MinIO
- 📋 Просмотр списка файлов с пагинацией
- 🗑 Удаление файлов с проверкой владельца
- 🛡 Валидация входящих данных и обработка ошибок
- 📝 Структурированное логирование через zap (stdout + файл)

---

## 📂 Структура проекта

```
.
├── cmd/
│   └── main.go                  # Точка входа в приложение
├── config/
│   └── config.go                # Загрузка конфигурации из .env
├── internal/
│   ├── auth/                    # Регистрация и авторизация
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── payload.go
│   │   └── handler_test.go
│   ├── file/                    # Управление файлами
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── handler_test.go
│   ├── user/                    # Модель и репозиторий пользователя
│   │   ├── model.go
│   │   └── repository.go
│   └── db/                      # Подключение к PostgreSQL
│       └── db.go
├── migrations/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── pkg/
│   ├── middleware/              # JWT Auth middleware
│   ├── logger/                  # Настройка zap логгера
│   └── response/                # JSON-ответы
├── tests/
│   └── integration_test.go      # Интеграционные тесты
├── docker-compose.yaml
├── Dockerfile
├── go.mod
└── .env                         # Переменные окружения (создать вручную)
```

---

## 🚀 Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/your-username/file-storage.git
cd file-storage
```

### 2. Создать файл `.env`

```env
DATABASE_URL=postgres://postgres:password@localhost:5433/postgres?sslmode=disable

JWT_SECRET=your_jwt_secret_key
JWT_TTL=24h

S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=filestorage

PORT=8080
MAX_FILE_SIZE=10485760
```

### 3. Запустить все сервисы через Docker

```bash
docker-compose up --build
```

Поднимает PostgreSQL, MinIO, создаёт бакет `filestorage` и запускает приложение.

### 4. Или запустить локально

```bash
# Поднять только инфраструктуру
docker-compose up postgres minio minio-create-bucket -d

# Применить миграции
psql $DATABASE_URL -f migrations/000001_init.up.sql

# Запустить сервер
go run cmd/main.go
```

Сервер запустится на порту **:8080**

---

## 📜 API

### 🔐 Аутентификация

#### Регистрация
```
POST /sign-up
```
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```
Ответ `201 Created`:
```json
{
  "id": 1,
  "email": "user@example.com"
}
```

#### Вход
```
POST /sign-in
```
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```
Ответ `200 OK`:
```json
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### 📁 Файлы

> Все маршруты 🔒 требуют заголовка `Authorization: Bearer <token>`

| Метод | Маршрут | Описание | Доступ |
|-------|---------|----------|--------|
| `POST` | `/upload` | Загрузить файл | 🔒 |
| `GET` | `/files?page=1&limit=10` | Список файлов | 🔒 |
| `DELETE` | `/files/{id}` | Удалить файл | 🔒 |

#### Загрузка файла

Тип запроса: `multipart/form-data`. Поле файла: `file`.  
Разрешённые форматы: `image/png`, `image/jpeg`.

```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@photo.jpg"
```

Ответ `200 OK`:
```json
{ "message": "file uploaded" }
```

#### Получение списка файлов

| Параметр | Тип | Описание |
|----------|-----|----------|
| `page` | int | Номер страницы (по умолчанию 1) |
| `limit` | int | Элементов на странице (по умолчанию 1, максимум 100) |

Ответ `200 OK`:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "url": "http://localhost:9000/filestorage/photo.jpg",
    "file_name": "photo.jpg",
    "size": 204800,
    "created_at": "2025-06-09T12:00:00Z"
  }
]
```

#### Удаление файла

```
DELETE /files/{id}
```

Удаляет файл из MinIO и базы данных. Возвращает `403 Forbidden`, если файл принадлежит другому пользователю. Успешный ответ: `204 No Content`.

---

## 🌐 Используемые технологии

- **Роутинг:** стандартная библиотека `net/http`
- **База данных:** [pgx v5](https://github.com/jackc/pgx) + pgxpool
- **Объектное хранилище:** [MinIO Go SDK](https://github.com/minio/minio-go)
- **JWT:** [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)
- **Пароли:** `golang.org/x/crypto/bcrypt`
- **Валидация:** [go-playground/validator v10](https://github.com/go-playground/validator)
- **Логирование:** [uber-go/zap](https://github.com/uber-go/zap)
- **Конфигурация:** [godotenv](https://github.com/joho/godotenv)

---

## 🧪 Тестирование

```bash
# Юнит-тесты (без внешних зависимостей)
go test ./internal/...

# Интеграционные тесты (требуется запущенный PostgreSQL)
go test ./tests/...
```

Юнит-тесты покрывают хендлеры `auth` (SignUp, SignIn, валидация невалидных данных) и `file` (Upload, GetFiles, Delete) через моки сервисов.

---

## 🚨 Ограничения

- Поддерживаются только форматы `image/png` и `image/jpeg`
- Требуется активное подключение к PostgreSQL и MinIO
- JWT не поддерживает отзыв токенов (refresh token не реализован)
- URL файла генерируется с хардкодом `localhost:9000` — в продакшене нужно заменить на реальный домен

---

## ❓ Часто задаваемые вопросы (FAQ)

**Q: Как авторизоваться в защищённых маршрутах?**  
A: После входа через `/sign-in` вы получаете JWT-токен. Передавайте его в заголовке: `Authorization: Bearer <token>`.

**Q: Где хранятся загруженные файлы?**  
A: Файлы хранятся в MinIO в бакете `filestorage`. Метаданные (имя, размер, ссылка) сохраняются в PostgreSQL.

**Q: Что происходит при удалении файла?**  
A: Сервис сначала проверяет, что файл принадлежит текущему пользователю, затем удаляет его из MinIO и из базы данных.

**Q: Какие коды ошибок могут возникать?**  
A:
- `400 Bad Request` — невалидные данные или неподдерживаемый формат файла
- `401 Unauthorized` — отсутствует или невалидный JWT-токен
- `403 Forbidden` — попытка удалить чужой файл
- `409 Conflict` — email уже зарегистрирован
- `500 Internal Server Error` — внутренняя ошибка сервера

**Q: Можно ли запустить без Docker?**  
A: Да, если у вас уже установлены PostgreSQL и MinIO. Укажите корректные значения в `.env` и пропустите шаг с Docker.

---

## 🛠 Планы развития

### Ближайшие релизы (v1.1–1.2)
- 🔄 Refresh токены и отзыв JWT
- 📎 Поддержка любых типов файлов (PDF, видео, архивы)
- 🔍 Фильтрация файлов по имени и дате
- 📏 Настраиваемый лимит размера файла из конфига

### Средний план (v2.0)
- 🗂 Папки и организация файлов
- 🔗 Публичные ссылки на файлы с TTL
- 📈 Статистика загрузок и использования места
- 🖼 Автоматическая генерация превью для изображений

### Долгосрочные планы (v3.0+)
- 🐳 Kubernetes-манифесты для деплоя
- 🌍 Поддержка AWS S3 и GCS как альтернативы MinIO
- 📊 Веб-дашборд для управления файлами
- 👥 Командный доступ и шаринг файлов между пользователями

---

## 📝 Лицензия

Проект распространяется под [MIT License](LICENSE).

---

## 🤝 Вклад в проект

Приветствуются pull requests и issue reports. Перед внесением изменений:

1. Форкните репозиторий
2. Создайте feature branch (`git checkout -b feature/my-feature`)
3. Закоммитьте изменения (`git commit -m 'Add my feature'`)
4. Запушьте в branch (`git push origin feature/my-feature`)
5. Создайте Pull Request

---

## 📞 Поддержка

Для вопросов и предложений создавайте issue в репозитории проекта.
