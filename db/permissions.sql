CREATE USER shrt_url_default WITH  PASSWORD '123';
GRANT SELECT, INSERT, DELETE, UPDATE on urls to shrt_url_default;

GRANT USAGE on urls_id_seq to shrt_url_default;
