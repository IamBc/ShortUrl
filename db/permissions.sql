CREATE USER my_user;
CREATE USER my_user WITH  PASSWORD 'my_pass';
GRANT SELECT, INSERT, DELETE, UPDATE on urls to my_user;

-- POSTGRES SPECIFIC
GRANT USAGE on urls_id_seq to my_user;
