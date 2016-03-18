CREATE TABLE urls (id SERIAL PRIMARY KEY,
                   url_hash UNIQUE TEXT NOT NULL,
                   url TEXT UNIQUE NOT NULL);
