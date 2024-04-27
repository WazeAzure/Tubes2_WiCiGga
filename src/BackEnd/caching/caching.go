package caching

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"

	"git.mills.io/prologic/bitcask"
)

var CachedWebpage *bitcask.Bitcask
var CachedRedirect *bitcask.Bitcask

var CachedRootFolder string = ".cache/"

func CheckCacheFolder() {
	if _, err := os.Stat(".cache"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(".cache", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func CheckCacheFile(url string) bool {
	key := GetKeyHash(url)

	if _, err := os.Stat(CachedRootFolder + key); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CheckCacheRedirect(url string) bool {
	key := GetKeyHash(url)

	return CachedRedirect.Has([]byte(key))
}

func GetKeyHash(url string) string {
	keyHash := fnv.New64()
	keyHash.Write([]byte(url))
	return fmt.Sprintf("%v", keyHash.Sum64())
}

func InitCache() {
	var err error
	CachedWebpage, err = bitcask.Open(CachedRootFolder + "#cached-webpage")
	if err != nil {
		log.Fatalln(err)
	}

	CachedRedirect, err = bitcask.Open(CachedRootFolder + "#cached-redirect")
	if err != nil {
		log.Fatalln(err)
	}
	_ = CachedWebpage
	_ = CachedRedirect
}

func SetCacheVisited(currenct_url string) {
	key := GetKeyHash(currenct_url)
	if CachedWebpage.Has([]byte(key)) {
		return
	}

	CachedWebpage.Put([]byte(key), []byte(currenct_url))
}

func SetCacheUrl(current_url string, list_url []string) {

	// pastikan list url sudah bersih dan sudah di handle redirect sebelumnya!
	key := GetKeyHash(current_url)
	var err error
	Db, err := bitcask.Open(CachedRootFolder + key)
	if err != nil {
		log.Fatalln(err)
	}

	for val := range list_url {
		if !Db.Has([]byte(list_url[val])) {
			Db.Put([]byte(list_url[val]), []byte{1})
		}
	}

	SetCacheVisited(current_url)

	defer Db.Close()
}

func SetCacheRedirect(current_url string, redirect_url string) {
	key := GetKeyHash(current_url)

	if !CachedRedirect.Has([]byte(key)) {
		CachedRedirect.Put([]byte(key), []byte(redirect_url))
	}
}

func GetCacheUrl(current_url string) []string {
	// sudah pernah dikunjungi
	var list_url []string

	key := GetKeyHash(current_url)

	var err error
	Db, err := bitcask.Open(CachedRootFolder + key)
	if err != nil {
		log.Fatalln(err)
	}

	x := Db.Keys()

	for val := range x {
		list_url = append(list_url, string(val))
	}

	defer Db.Close()

	return list_url
}

func GetCacheRedirect(current_url string) string {
	key := GetKeyHash(current_url)

	ans, err := CachedRedirect.Get([]byte(key))
	if err != nil {
		log.Fatalln(err)
	}
	return string(ans)
}

func DeleteCache(url string) {
	key := GetKeyHash(url)

	CachedWebpage.Delete([]byte(key))
}
