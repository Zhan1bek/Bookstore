CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          user_id INT NOT NULL,
                          book_id INT NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          FOREIGN KEY (user_id) REFERENCES users(id),
                          FOREIGN KEY (book_id) REFERENCES books(id)
);