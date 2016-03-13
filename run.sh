echo "Initializing environment..."

export SHORT_URL_FILES_DIR="" #Must be absolute path !!!
export SHORT_URL_FILE_PORT=""
export SHORT_URL_API_PORT=""
export DB_CONNECTION_STRING=""
export DB_CONNECTION_DRIVER=""
export URL_CHECK_STRICT_MODE=""


echo "Starting server..."
/usr/local/go/bin/go run server/main.go server/storage_rdb.go server/cache.go  -logtostderr=true

