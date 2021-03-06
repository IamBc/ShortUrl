echo "Initializing environment..."


#NOTE: This is example configuration. If you run the install script this configuration will be enough to run ShortUrl.
# It is **STRONGLY** encouraged for the dbname, user and password to be changed.


export SHORT_URL_FILES_DIR="./web-client"
export SHORT_URL_FILE_PORT="9002"
export SHORT_URL_API_PORT="9003"
export DB_CONNECTION_STRING="user=shrt_url_default password=123 dbname=short_url_def sslmode=disable"
export DB_CONNECTION_DRIVER="postgres"
export MAX_URL_COUNT="100"
export MAX_URL_LENGTH="300"

echo "Starting server..."
./server/server -logtostderr=true
