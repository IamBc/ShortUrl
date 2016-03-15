package main

/*
File containing the interface that must be implemented in order to store the ShortUrl data somewhere
In case a new storage type will be implemented (eg in a file) it must implement these functions
*/

type storage interface {
	GetURLFromStorage(urlHash string) (string, error)
	DeleteURL(urlHash string) error
	AddURLToStorage(urlHash string, bodyStr string) error
}
