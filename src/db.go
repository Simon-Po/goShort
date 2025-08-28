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
	CheckForCollion(shortUrl *string) error
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
	for line := range strings.SplitSeq(db_string, "\n") {

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
		for line := range strings.SplitSeq(db_string, "\n") {

			content := strings.Split(line, " ")

			if len(content) > 2 || len(content) < 2 {
				continue
			}
			db_url, sUrl := content[0], content[1]
			db.content[db_url] = sUrl
		}

	}
}

func (db *textFileDb) CheckUrl(url string) (string, error) {
	db.textFileDbGetBuffer()
	log.Println("------------------------------")
	log.Println(db.content)
	log.Println("------------------------------")
	return strings.TrimSpace(db.content[url]), nil
}

func (db *textFileDb) WriteUrl(url, shortUrl string) error {
	f, err := os.OpenFile("testdb.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Println("Testdb could not be opened")
		return err
	}

	input := url + " " + shortUrl + "\n"
	_, err = f.Write([]byte(input))
	if err != nil {
		log.Println(err)
		log.Println("could not write to Testdb")
		return err
	}

	db.textFileDbRefreshBuffer()
	return nil
}

func (db *textFileDb) CheckForCollion(shortUrl *string) error {
	for entry := range db.content{
		if(entry == *shortUrl) {
			return errors.New("") 
		}
	}
	return nil
}

