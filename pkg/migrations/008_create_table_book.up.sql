ALTER TABLE books
    ADD COLUMN avg_rating FLOAT DEFAULT 0.0,
    ADD COLUMN rating_count INT DEFAULT 0;