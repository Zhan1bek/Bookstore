
CREATE TABLE ratings (
                         id SERIAL PRIMARY KEY,
                         user_id INTEGER NOT NULL REFERENCES users(id),
                         book_id INTEGER NOT NULL REFERENCES books(id),
                         rating INTEGER CHECK (rating >= 1 AND rating <= 5),
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

