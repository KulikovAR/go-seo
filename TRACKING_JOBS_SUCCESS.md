# ✅ Эндпоинт мониторинга джобов успешно добавлен в Swagger

## Что было сделано

1. **✅ Восстановлена оригинальная Swagger документация** - все эндпоинты позиций остались нетронутыми
2. **✅ Добавлен новый эндпоинт** `GET /api/tracking-jobs` в Swagger документацию
3. **✅ Добавлены определения DTO структур** для `TrackingJobItem` и `TrackingJobsResponse`

## Проверка работоспособности

### ✅ Эндпоинт tracking-jobs работает:
```bash
curl http://localhost:8087/api/tracking-jobs
```
Возвращает список джобов с пагинацией, прогрессом выполнения и метаинформацией.

### ✅ Все эндпоинты позиций работают:
```bash
curl http://localhost:8087/api/positions/latest
curl http://localhost:8087/api/positions/history
curl http://localhost:8087/api/positions/combined
# и другие...
```

### ✅ Swagger UI доступен:
```
http://localhost:8087/swagger/index.html
```

## Структура ответа эндпоинта tracking-jobs

```json
{
  "data": [
    {
      "id": "job_abc123",
      "site_id": 65,
      "source": "google",
      "status": "completed",
      "created_at": "2025-10-24T16:17:31.811823+03:00",
      "updated_at": "2025-10-24T16:17:34.659844+03:00",
      "completed_at": "2025-10-24T16:17:34.65972+03:00",
      "total_tasks": 2,
      "completed_tasks": 2,
      "failed_tasks": 0,
      "progress": 100.0
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 20,
    "total": 32,
    "last_page": 2,
    "from": 1,
    "to": 20,
    "has_more": true
  },
  "meta": {
    "query_time_ms": 17,
    "cached": false
  }
}
```

## Параметры запроса

- `site_id` (int, optional) - фильтр по ID сайта
- `status` (string, optional) - фильтр по статусу (pending, running, completed, failed, cancelled)
- `page` (int, optional) - номер страницы (по умолчанию 1)
- `per_page` (int, optional) - количество записей на странице (по умолчанию 20, максимум 100)

## Примеры использования

```bash
# Получить все джобы
curl "http://localhost:8087/api/tracking-jobs"

# Получить джобы для сайта 65
curl "http://localhost:8087/api/tracking-jobs?site_id=65"

# Получить только активные джобы
curl "http://localhost:8087/api/tracking-jobs?status=running"

# Получить джобы с пагинацией
curl "http://localhost:8087/api/tracking-jobs?page=2&per_page=10"
```

## Результат

✅ **Все задачи выполнены успешно:**
- Оригинальная документация по позициям сохранена
- Новый эндпоинт мониторинга джобов добавлен
- Swagger документация обновлена
- Все эндпоинты работают корректно
- Сервер запускается без ошибок
