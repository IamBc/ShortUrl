sudo su postgres -c 'psql -c "CREATE DATABASE short_url_def;"'
sudo su postgres -c 'psql -d short_url_def  -f ./db/schema.sql'
sudo su postgres -c 'psql -d short_url_def  -f ./db/permissions.sql'
