
[![forthebadge](https://forthebadge.com/images/featured/featured-built-with-love.svg)](https://forthebadge.com)
# HTTP link shortener
Сокращатель ссылок

Тестовое задание 

Используемые технологии: 
- PostgreSQL (в качестве основного хранилища данных)
- Redis (в качестве хранилища данных для быстрого получения ссылки)
- Docker (для запуска сервиса)
- Fiber (веб фреймворк)
- golang-migrate/migrate (для миграций БД)

# Usage

**Скопируйте проект**
```bash
  git clone https://github.com/artemiyKew/http-link-shortener.git
```

**Перейдите в каталог проекта**
```bash
  cd http-link-shortener
```

**Запустите сервер**
```bash
  make compose
```

## Examples
- [Создание сокращенной ссылки](#создание-сокращенной-ссылки)
- [Редирект](#редирект)
- [Получение данных о сокращенной ссылке](#получение-данных-о-сокращенной-ссылке)

## Создание сокращенной ссылки
Создание сокращенной ссылки: 

```bash
curl -X POST \
    http://localhost:1234/ \
    -H "Content-Type: application/json" \
    -d '{
    "long_url": "https://vk.com"
    }'
```
Пример ответа: 
```json
{
    "full_url":"https://vk.com",
    "create_at":"2023-12-01T18:34:32.555217917Z", 
    "expired_at":"2023-12-02T18:34:32.555218084Z",
    "visit_counter":0,
    "token":"1395ec37e4",
}
```

## Редирект
Редирект:
```bash
curl http://localhost:1234/1395ec37e4 
```

## Получение данных о сокращенной ссылке
Получение данных о сокращенной ссылке:

```bash
 curl -X GET \
    http://localhost:1234/ \
    -H "Content-Type: application/json" \
    -d '{
    "token": "1395ec37e4"
    }'
```
Пример ответа: 
```json
{
    "full_url":"https://vk.com",
    "create_at":"2023-12-01T18:34:32.555218Z",
    "expired_at":"2023-12-02T18:34:32.555218Z",
    "visit_counter":2,
    "token":"1395ec37e4"    
}
```