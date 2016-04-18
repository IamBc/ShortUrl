package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"errors"
	"io/ioutil"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"expvar"
	"fmt"
	"unicode/utf8"
)

var cache *Cache // In memory cache

var (
	reqCounters = expvar.NewMap("reqCounters")
)

func main() {
	//Initialize glog
	flag.Parse()
	glog.Info("default server mux: ", http.DefaultServeMux)

	//Initialize cache
	cache = NewCache(10)

	//Initialize the DB
	InitStorage()

	//Init the appstate flags
	reqCounters.Add(`addUserHashReqCount`, 0)
	reqCounters.Add(`addReqCount`, 0)
	reqCounters.Add(`checkReqCount`, 0)
	reqCounters.Add(`redirectReqCount`, 0)

	//Used in url hash generation
	rand.Seed(time.Now().UTC().UnixNano())

	//File server
	go http.ListenAndServe(":"+os.Getenv("SHORT_URL_FILE_PORT"), http.FileServer(http.Dir(os.Getenv("SHORT_URL_FILES_DIR"))))
	glog.Info("Started the file server on port:" + os.Getenv("SHORT_URL_FILE_PORT"))

	//Initialize the REST API routes
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/g/{urlHash}", Redirect)
	router.HandleFunc("/check/{urlHash}", Check)
	router.HandleFunc("/add/", Add)
	router.HandleFunc("/add_user_hash/{userSelectedHash}", AddUserSelectedHash)
	router.HandleFunc("/delete/{urlHash}", Remove)
	router.HandleFunc("/appstate/", expvarHandler)

	glog.Info("Starting the API server on port:" + os.Getenv("SHORT_URL_API_PORT"))
	glog.Info(http.ListenAndServe(":"+os.Getenv("SHORT_URL_API_PORT"), router))

}

/* Handlers  */
func Redirect(w http.ResponseWriter, r *http.Request) {
	reqCounters.Add(`redirectReqCount`, 1)
	if r.Method != "GET" {
		reqCounters.Add(`wrongMethodRequests`, 1)
		WriteResp(w, http.StatusMethodNotAllowed, `Wrong method!`)
		return
	}

	vars := mux.Vars(r)

	//Check if URL is in cache
	cachedUrl := cache.GetVal(vars[`urlHash`])
	if cachedUrl != `` {
		glog.Info(`Redirecting (from cache) with hash: `, vars[`urlHash`], ` to: `, cachedUrl)
		http.Redirect(w, r, cachedUrl, http.StatusMovedPermanently)
		return
	}

	//URL is not in cache, fetch it and cache it
	url, err := GetURLFromStorage(vars[`urlHash`])
	if err != nil {
		WriteResp(w, http.StatusInternalServerError, `Please try again later!`)
		return
	}
	glog.Info(`Redirecting with hash: `, vars[`urlHash`], ` to: `, url)
	cache.Add(vars[`urlHash`], url)

	if url == `` {
		WriteResp(w, http.StatusNotFound, `Not found!`)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	reqCounters.Add(`checkReqCount`, 1)
	if r.Method != "GET" {
		reqCounters.Add(`wrongMethodRequests`, 1)
		WriteResp(w, http.StatusMethodNotAllowed, `Wrong method!`)
		return
	}

	vars := mux.Vars(r)
	url, err := GetURLFromStorage(vars[`urlHash`])
	if err != nil {
		WriteResp(w, http.StatusInternalServerError, `Please try again later!`)
		return
	}
	glog.Info(`Redirecting with hash: `, vars[`urlHash`], ` to: `, url)

	if url == `` {
		WriteResp(w, http.StatusNotFound, `Not found!`)
		return
	}

	WriteResp(w, http.StatusOK, url)
}

func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	reqCounters.Add(`addReqCount`, 1)

	if r.Method != "POST" {
		reqCounters.Add(`wrongMethodRequests`, 1)
		WriteResp(w, http.StatusMethodNotAllowed, `Wrong method!`)
		glog.Error(r.Method)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteResp(w, http.StatusInternalServerError, `Please try again later!`)
		return
	}
	bodyStr := string(body) //bytes to string
	glog.Info(bodyStr)
	if bodyStr == `` {
		WriteResp(w, http.StatusBadRequest, `Invalid URL!`)
		return
	}

	maxURLLength, err := strconv.Atoi(os.Getenv(`MAX_URL_LENGTH`))
	if err != nil {
		WriteResp(w, http.StatusBadRequest, `Invalid URL!`)
		glog.Error(err, ` Invalid URL:  `, bodyStr)
		return
	}

	if utf8.RuneCountInString(bodyStr) > maxURLLength {
		WriteResp(w, http.StatusBadRequest, `URL exceeds max length!`)
		glog.Error(err, `URL is too long!`)
		return
	}

	err = checkUrl(bodyStr)
	if err != nil {
		WriteResp(w, http.StatusBadRequest, `Invalid URL!`)
		glog.Error(err, ` Invalid URL:  `, bodyStr)
		return
	}

	urlHash := GenerateHash(bodyStr, 8)
	glog.Error(`urlHash: `, urlHash)
	urlHash, err = AddURLToStorage(urlHash, bodyStr)

	if err != nil && urlHash == `ui` {
		urlHash = ``
		WriteResp(w, http.StatusInternalServerError, err.Error())
		return
	} else if err != nil {
		WriteResp(w, http.StatusInternalServerError, `Please try again later!`)
		return
	}
	w.Write([]byte(urlHash))
}

func AddUserSelectedHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	reqCounters.Add(`addUserHashReqCount`, 1)

	if r.Method != "POST" {
		reqCounters.Add(`wrongMethodRequests`, 1)
		WriteResp(w, http.StatusMethodNotAllowed, `Wrong method!`)
		glog.Error(r.Method)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteResp(w, http.StatusInternalServerError, `Please try again later!`)
		return
	}
	bodyStr := string(body) //bytes to string
	glog.Info(bodyStr)
	if bodyStr == `` {
		WriteResp(w, http.StatusBadRequest, `Invalid URL!`)
		return
	}

	err = checkUrl(bodyStr)
	if err != nil {
		WriteResp(w, http.StatusBadRequest, `Invalid URL!`)
		glog.Error(err, ` Invalid URL:  `, bodyStr)
		return
	}

	vars := mux.Vars(r)
	glog.Error(`userSelectedHash: `, vars[`userSelectedHash`])
	vars[`userSelectedHash`], err = AddURLToStorage(vars[`userSelectedHash`], bodyStr)

	if err != nil && vars[`userSelectedHash`] == `ui` {
		vars[`userSelectedHash`] = ``
		WriteResp(w, http.StatusInternalServerError, err.Error())
		return
	} else if err != nil {
		WriteResp(w, http.StatusInternalServerError, `Please try again later!`)
		return
	}
	w.Write([]byte(vars[`userSelectedHash`]))
}

func Remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	reqCounters.Add(`removeReqCount`, 1)
	if r.Method != "DELETE" {
		reqCounters.Add(`wrongMethodRequests`, 1)
		WriteResp(w, http.StatusMethodNotAllowed, `Wrong method!`)
		return
	}

	vars := mux.Vars(r)
	err := DeleteURL(vars[`urlHash`])
	if err != nil {
		WriteResp(w, http.StatusInternalServerError, `Please try again later!`)
	}
	glog.Info(`Deleted hash: `, vars[`urlHash`])

	w.Write([]byte(`Deleted hash: ` + vars[`urlHash`]))
}

// HELPER FUNCTIONS
func WriteResp(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

func GenerateHash(inp string, length int) string {
	reqCounters.Add(`generatedHashes`, 1)
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func checkUrl(input string) error {

	//u, err := url.UrlParse(input)
	_, err := url.ParseRequestURI(strings.Split(input, `#`)[0])
	if err != nil {
		glog.Error("  error:", err)
		return errors.New(`Not an URL!`)
	}

	return nil
}

func expvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
