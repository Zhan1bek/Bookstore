# Учетная система для небольшого магазина книг

## О проекте

Этот проект представляет собой веб-приложение, созданное для упрощения учета товаров, заказов и клиентов в небольшом магазине книг. Приложение позволяет владельцам магазина эффективно управлять каталогом книг, отслеживать заказы и взаимодействовать с клиентами через удобный интерфейс.

### Основные функции

- **Управление каталогом книг**: Добавление, обновление, удаление и просмотр информации о книгах.
- **Управление заказами**: Создание, обновление и отслеживание статуса заказов.
- **Управление клиентами**: Регистрация, обновление информации клиентов и просмотр истории их заказов.

## Rest API
- **POST /books - добавление новой книги в каталог.
- **GET /books/:id - получение информации о книге по ID.
- **PUT /books/:id - обновление информации о книге по ID.
- **DELETE /books/:id - удаление книги из каталога.
  
#### Для управления заказами:

- **POST /orders - создание нового заказа.
- **GET /orders/:id - получение информации о заказе по ID.
- **PUT /orders/:id - обновление статуса заказа по ID.
- **DELETE /orders/:id - отмена заказа.

#### Для работы с клиентами:

- **POST /customers - добавление нового клиента.
- **GET /customers/:id - получение информации о клиенте по ID.
- **PUT /customers/:id - обновление информации клиента по ID.
- **DELETE /customers/:id - удаление клиента из системы.

## Структура сущностей в базе данных и их связи

В проекте используется следующая структура базы данных для управления каталогом книг, заказами и клиентами.

### Сущности:

#### Книги (`books`)

- `id` (bigserial): Уникальный идентификатор книги (первичный ключ).
- `title` (text): Название книги.
- `author` (text): Автор книги.
- `price` (numeric): Цена книги.
- `stock_quantity` (integer): Количество книг на складе.
- `created_at` (timestamp): Дата и время добавления книги в каталог.
- `updated_at` (timestamp): Дата и время последнего обновления информации о книге.

#### Заказы (`orders`)

- `id` (bigserial): Уникальный идентификатор заказа (первичный ключ).
- `customer_id` (bigserial): Идентификатор клиента, сделавшего заказ (внешний ключ).
- `status` (text): Статус заказа (например, "новый", "в обработке", "завершен").
- `created_at` (timestamp): Дата и время создания заказа.
- `updated_at` (timestamp): Дата и время последнего обновления статуса заказа.

#### Клиенты (`customers`)

- `id` (bigserial): Уникальный идентификатор клиента (первичный ключ).
- `name` (text): Имя клиента.
- `email` (text): Электронная почта клиента.
- `phone` (text): Телефонный номер клиента.
- `created_at` (timestamp): Дата и время добавления клиента.
- `updated_at` (timestamp): Дата и время последнего обновления информации о клиенте.

#### Элементы заказа (`order_items`)

- `id` (bigserial): Уникальный идентификатор элемента заказа (первичный ключ).
- `order_id` (bigserial): Идентификатор заказа, к которому относится элемент (внешний ключ).
- `book_id` (bigserial): Идентификатор книги в элементе заказа (внешний ключ).
- `quantity` (integer): Количество книг в элементе заказа.

### Связи:

- Каждый заказ (`orders`) может содержать один или несколько элементов заказа (`order_items`), а каждый элемент заказа принадлежит одному заказу.
- Каждый элемент заказа (`order_items`) связан с одной книгой (`books`), а книги могут входить в состав нескольких элементов заказа.
- Заказы (`orders`) связаны с клиентами (`customers`), при этом один клиент может сделать несколько заказов.


### Zhanibek Myrzakhanov 22B030409