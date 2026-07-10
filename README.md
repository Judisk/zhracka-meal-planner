# Food Planner

Desktop-приложение на Go для генерации простого сбалансированного рациона. Продукты хранятся в SQLite, управляются через GUI на Fyne, а генерация рациона использует взвешенный вероятностный выбор, чтобы одни и те же продукты не повторялись слишком часто.

Проект написан как практика Go: разделение логики по слоям, работа с SQLite и транзакциями, обработка ошибок через sentinel errors, table-driven тесты и GUI на Fyne.

## Возможности

- Таблица продуктов в GUI с добавлением, редактированием и удалением.
- Фильтрация по категории, статусу banned/allowed и предпочтению.
- Сортировка по ID и имени (по возрастанию/убыванию).
- Генерация одного блюда.
- Генерация дня питания на 1–6 приёмов пищи.
- Сохранение сгенерированного дня с обновлением `selection_score`.
- Учёт banned-продуктов при генерации.
- Взвешенный выбор продуктов через `selection_score` и `preference`.
- SQLite без CGO через `modernc.org/sqlite`.

## Как работает генерация

Продукты делятся на три категории: `Grain`, `Protein`, `Vegetable`.

Одно блюдо — это три продукта, по одному из каждой категории:

```text
Dish = Grain + Protein + Vegetable
Day  = []Dish
```

При выборе продукт берётся не равномерно, а с учётом `selection_score`:

- чем выше `selection_score`, тем выше шанс выбора;
- выбранные продукты получают cooldown (сброс score), чтобы не повторяться;
- после сохранения дня scores обновляются;
- `preference` (Liked / Neutral / Disliked) влияет на скорость роста score.

## Архитектура

Зависимости между слоями направлены в одну сторону, без циклов:

```text
gui  →  service  →  dayone  →  foodgenerator
                 →  storage
(все слои используют пакет products)
```

- `products` — доменные типы: `Product`, `Category`, `PreferenceStatus`, `BlockedProducts`.
- `foodgenerator` — генерация одного блюда, взвешенный выбор.
- `dayone` — сборка дня из блюд, проверка banned-продуктов.
- `storage` — SQLite: CRUD, выборки, обновление scores, транзакции.
- `service` — оркестрация между storage / dayone для GUI.
- `gui` — интерфейс на Fyne.

`dayone` намеренно не импортирует `storage` — вся работа с БД проходит через `service`.

## Структура проекта

```text
.
├── cmd/
│   └── main.go                  # точка входа
├── internal/
│   ├── products/                # доменные типы
│   ├── foodgenerator/           # генерация блюда
│   ├── dayone/                  # генерация дня, banned-валидатор
│   ├── storage/                 # SQLite: CRUD, выборки, scores
│   ├── service/                 # бизнес-логика, seed
│   └── gui/                     # Fyne-интерфейс
├── go.mod
└── go.sum
```

## Требования

- Go `1.26.4` или новее (см. `go.mod`).
- Графическое окружение для запуска Fyne.

Прямые зависимости:

- `fyne.io/fyne/v2`
- `modernc.org/sqlite`

## Запуск

```bash
go run ./cmd
```

При старте приложение создаёт (или открывает) базу `products.db` рядом с местом запуска и, если таблица пуста, заполняет её набором продуктов по умолчанию через `SeedDefaultProductsIfEmpty`.

## Тесты

```bash
go test ./...
```

Покрыты генерация блюда и дня, storage- и service-слои, а также ошибочные сценарии: пустая база, banned-продукты, невалидные значения `n`.

## База данных

Таблица создаётся автоматически при первом запуске. Ограничения на уровне БД гарантируют корректность данных даже при прямой вставке:

```sql
CREATE TABLE IF NOT EXISTS products (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL CHECK (trim(name) <> ''),
    category        TEXT NOT NULL CHECK (category IN ('Grain', 'Protein', 'Vegetable')),
    banned          INTEGER NOT NULL DEFAULT 0 CHECK (banned IN (0, 1)),
    preference      REAL NOT NULL DEFAULT 1.0 CHECK (preference IN (0.5, 1.0, 1.5)),
    selection_score REAL NOT NULL DEFAULT 1.0
);
```

Файл базы игнорируется через `.gitignore` (`*.db`).

## Roadmap

- Балансировка рациона по КБЖУ.
- Экран статистики по scores и истории генераций.
- GitHub Actions для прогонки `go test ./...` на каждый push.
