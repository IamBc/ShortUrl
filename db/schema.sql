
-- MYSQL
CREATE TABLE urls (
    id  int(10) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    url_hash TEXT NOT NULL,
    url TEXT NOT NULL 
);

create unique index duplicate_url on urls(url(500));

-- POSTGRES

CREATE TABLE urls (id SERIAL PRIMARY KEY,
                   url_hash TEXT NOT NULL,
                   url TEXT NOT NULL);
