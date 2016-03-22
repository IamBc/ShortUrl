echo "Initializing environment..."


#NOTE: This is example configuration. If you run the install script this configuration will be enough to run ShortUrl.
# It is **STRONGLY** encouraged for the dbname, user and password to be changed.

export SHORT_URL_FILES_DIR="./web-client" 
export SHORT_URL_FILE_PORT="9002"
export SHORT_URL_API_PORT="9003"
export DB_CONNECTION_STRING="user=shrt_url password=123 dbname=short_url sslmode=disable"
export DB_CONNECTION_DRIVER="postgres"


echo "Starting server..."
./server/server -logtostderr=true

