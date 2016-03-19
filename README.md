# ShortUrl

Yet another URL shortening service. Coming with a simple (web) interface, JSON REST API and capable of serving thousands of requests per second per core.



## Installation

### Packages
In order to run ShortUrl you need to have the following packages installed:
* PostgreSQL
* Go(lang)

Although the names of the packages may vary from distribution to distribution for ubuntu/debian you can go with
```bash
 sudo apt-get install postgresql postgresql-contrib
 sudo apt-get install golang
```

### Golang libraries

Once you have golang installed you must install the golang specific libraries which are used by the project.
You can do this by typing: ```go install``` inside the ```server/``` directory in the project.
 

## Build and Run

To build and run the script you must trigger the script ```run.sh```, eg:
``` bash run.sh ```

Note that without configuring the project you will get some error.

## Configuration

The configuration of ShortUrl us done via environment variables which are set inside ```run.sh```

| Name                 | Description                                                                                                     | Example Value                                            |
|----------------------|-----------------------------------------------------------------------------------------------------------------|----------------------------------------------------------|
| SHORT_URL_FILES_DIR  | Absolute path to where the web-client files are stored.                                                         | /home/<user>/short_url/web-client                        |
| SHORT_URL_FILE_PORT  | Port on which the file server will be launched.                                                                 | 9002                                                     |
| SHORT_URL_API_PORT   | Port on which the Rest(ful) API will be hosted                                                                  | 9003                                                     |
| DB_CONNECTION_DRIVER | Driver name for the database that is being used.  See https://golang.org/s/sqldrivers for a list of drivers.    | postgres                                                 |
| DB_CONNECTION_STRING | Connection string for the database. Depends on the actual DB that is being used. The example is with PostgreSQL | user=user password=pass dbname=short_url sslmode=disable |

## Interface

ShortUrl comes with a simple web interface. From there you can create short urls. 

## REST API
| API Call | HTTP Method | Path              | Request Payload              | Description                                                                                                                                     |
|----------|-------------|-------------------|------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| Add      | POST        | /add/             | URL that must  be shortened. | Adds a link and returns it’s short code as a response.  Links with containing a hash (fragment identifier)  in their url path are not accepted. |
| Check    | GET         | /check/<urlHash>  | -                            | Returns the url to which this urlHash points to as a response.                                                                                  |
| Redirect | GET         | /g/<urlHash>      | -                            | Returns a 304 Moved permanently request, moved to the URL which was put .  The client should redirect to the original (long) url.               |
| Delete   | DELETE      | /delete/<urlHash> | -                            | Deletes the entry for urlHash. If the entry doesn’t exist it will still return a 200 OK status                                                  |
| Add User Selected Hash   | POST      | /add_user_hash/<urlHash> | -                            | Creates a short link with the urlHash which is selected by the API caller.                            |
## API Status Codes

| Code                      | Description                                              |
|---------------------------|----------------------------------------------------------|
| 200 OK                    | A request was completed successfully                     |
| 304 Moved permanently     | The correct response when visiting a short link.         |
| 404 Not found             | Method/File doesn’t exist                                |
| 400 Bad Request           | Wrong parameters were passed when making the API Call.   |
| 500 Internal Server Error | There was an internal error. The request can be retried. |

## API EXAMPLES
Let’s assume that we are testing the API from the host on which it’s deployed, using curl and it’s port is 9003.

### Invalid add request
```bash
-v --data '' http://localhost:9003/add/
* Connected to localhost (127.0.0.1) port 9003 (#0)
> POST /add/ HTTP/1.1
> User-Agent: curl/7.38.0
> Host: localhost:9003
> Accept: */*
> Content-Length: 0
> Content-Type: application/x-www-form-urlencoded
> 
< HTTP/1.1 400 Bad Request
< Access-Control-Allow-Origin: *
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Sat, 12 Mar 2016 18:49:28 GMT
< Content-Length: 13
* HTTP error before end of send, stop sending
< 
Invalid URL!
```

### Wrong method add request
```
curl -v http://localhost:9003/add/
* Connected to localhost (127.0.0.1) port 9003 (#0)
> GET /add/ HTTP/1.1
> User-Agent: curl/7.38.0
> Host: localhost:9003
> Accept: */*
> 
< HTTP/1.1 405 Method Not Allowed
< Access-Control-Allow-Origin: *
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Sat, 12 Mar 2016 18:50:09 GMT
< Content-Length: 14
< 
Wrong method!
* Connection #0 to host localhost left intact
```


### Successful add request

```
curl -v --data 'www.google.bg' http://localhost:9003/add/
* Connected to localhost (127.0.0.1) port 9003 (#0)
> POST /add/ HTTP/1.1
> User-Agent: curl/7.38.0
> Host: localhost:9003
> Accept: */*
> Content-Length: 13
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 13 out of 13 bytes
< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: *
< Date: Sat, 12 Mar 2016 18:50:40 GMT
< Content-Length: 8
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
4hxrNG9wz
```

### Successful redirect request

```
curl -v http://localhost:9003/g/4hxrNG9w
* Connected to localhost (127.0.0.1) port 9003 (#0)
> GET /g/4hxrNG9w HTTP/1.1
> User-Agent: curl/7.38.0
> Host: localhost:9003
> Accept: */*
> 
< HTTP/1.1 301 Moved Permanently
< Location: /g/www.google.bg
< Date: Sat, 12 Mar 2016 18:51:58 GMT
< Content-Length: 51
< Content-Type: text/html; charset=utf-8
< 
<a href="/g/www.google.bg">Moved Permanently</a>.

* Connection #0 to host localhost left intact
```

### Successful check request
```
curl -v http://localhost:9003/check/4hxrNG9w
* Connected to localhost (127.0.0.1) port 9003 (#0)
> GET /check/4hxrNG9w HTTP/1.1
> User-Agent: curl/7.38.0
> Host: localhost:9003
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Sat, 12 Mar 2016 18:53:08 GMT
< Content-Length: 14
< 
www.google.bg
```

### Successful delete request

```
curl -X 'DELETE' -v http://localhost:9003/delete/4hxrNG9w
* Connected to localhost (127.0.0.1) port 9003 (#0)
> DELETE /delete/4hxrNG9w HTTP/1.1
> User-Agent: curl/7.38.0
> Host: localhost:9003
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Sat, 12 Mar 2016 18:53:44 GMT
< Content-Length: 22
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
Deleted hash: 4hxrNG9w

```
