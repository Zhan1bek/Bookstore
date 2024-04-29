create table if not exists purchases(
    id serial primary key,
    user_id integer references users(id),
    book_id integer references books(id),
    quantity INTEGER,
    total_price DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)