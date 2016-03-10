package main

type storage interface {
    GetURLFromStorage(urlHash string) (string, error)
    DeleteURL(urlHash string) error
    AddURLToStorage(urlHash string, bodyStr string) error
}
