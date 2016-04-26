CREATE TABLE urls (id SERIAL PRIMARY KEY,
                   url_hash TEXT UNIQUE NOT NULL,
                   url TEXT UNIQUE NOT NULL);

CREATE TABLE short_url(id serial primary key, settings text);
insert into short_url values(1, '');
