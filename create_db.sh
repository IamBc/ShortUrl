sudo su postgres -c 'psql -c "CREATE DATABASE test1;"'
sudo su postgres -c 'psql -d test1  -f ./db/schema.sql'
sudo su postgres -c 'psql -d test1  -f ./db/permissions.sql'
