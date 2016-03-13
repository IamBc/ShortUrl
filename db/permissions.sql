CREATE USER my_user;
SET PASSWORD FOR 'my_user' = PASSWORD('my_pass');
GRANT SELECT, INSERT, DELETE, UPDATE on urls to my_user;
