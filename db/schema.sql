CREATE TABLE urls (id SERIAL PRIMARY KEY,
                   url_hash TEXT UNIQUE NOT NULL,
                   url TEXT UNIQUE NOT NULL);
