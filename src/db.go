package main

import (
	"errors"
	"log"
	"os"
	"strings"
)

type UrlDB interface {
	CheckUrl(url string) (string, error)
	WriteUrl(url, shortUrl string) error
	CheckForCollision(shortUrl *string) error
	CheckSurl(surl string) (string,error)
}

type textFileDbBuffer map[string]string

type textFileDb struct {
	pathToTxt string
	content   textFileDbBuffer
}

func (db *textFileDb) textFileDbRefreshBuffer() {
	db.content = nil
	data, err := os.ReadFile(db.pathToTxt)
	if err != nil {
		log.Fatal(db.pathToTxt + " Could not be found")
	}
	if len(data) < 1 {
		log.Println("Database is empty")
	}

	db.content = make(map[string]string)
	db_string := string(data)
    for _, line := range strings.Split(db_string, "\n") {

		content := strings.Split(line, " ")

		if len(content) > 2 || len(content) < 2 {
			continue
		}
		db_url, sUrl := content[0], content[1]
		db.content[db_url] = sUrl
	}
	log.Println("INFO: DB Buffer refreshed")

}
func (db *textFileDb) textFileDbGetBuffer() {
	if db.content == nil {
		data, err := os.ReadFile(db.pathToTxt)
		if err != nil {
			log.Fatal(db.pathToTxt + " Could not be found")
		}
		if len(data) < 1 {
			log.Println("Database is empty")
		}

		db.content = make(map[string]string)
		db_string := string(data)
        for _, line := range strings.Split(db_string, "\n") {

			content := strings.Split(line, " ")

			if len(content) > 2 || len(content) < 2 {
				continue
			}
			db_surl, db_url := content[0], content[1]
			db.content[db_surl] = db_url
		}

	}
}

func (db *textFileDb) CheckUrl(url string) (string, error) {
	db.textFileDbGetBuffer()
	log.Println("------------------------------")
	log.Println(db.content)
	log.Println("------------------------------")
	for k,v := range db.content {
		if v == url {
			return k,nil
		}
	}
	return "", nil
}

func (db *textFileDb) CheckSurl(surl string) (string,error) {
	db.textFileDbGetBuffer()
	log.Println("------------------------------")
	log.Println(db.content)
	log.Println("------------------------------")
	return strings.TrimSpace(db.content[surl]), nil
}

func (db *textFileDb) WriteUrl(url, shortUrl string) error {
	f, err := os.OpenFile("testdb.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERROR: Could not open file for Text DB")
	}
	defer f.Close()
	if err != nil {
		log.Println("Testdb could not be opened")
		return err
	}

	input := shortUrl + " " + url + "\n"
	_, err = f.Write([]byte(input))
	if err != nil {
		log.Println(err)
		log.Println("could not write to Testdb")
		return err
	}

	db.textFileDbRefreshBuffer()
	return nil
}

func (db *textFileDb) CheckForCollision(shortUrl *string) error {
	for entry := range db.content{
		if(entry == *shortUrl) {
			return errors.New("") 
		}
	}
	return nil
}
