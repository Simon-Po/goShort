package main

import (
	"encoding/json"
	"errors"
	"log"
	"io"
	"net/http"
	"os"
	"strings"
	"strconv"
	"github.com/google/uuid"
)


type ReqBody struct {
	Url    string `json:"url"`
	Length string  `json:"length"`
}

func create(w http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var rb ReqBody
	err = json.Unmarshal(b, &rb)
	if err != nil {
		log.Println("Error: Could not Unmarshal Body")
		http.Error(w, err.Error(), 500)
		return
	}
	found_url,_ := dbCheckIfExistsUrl(&rb.Url)
	log.Println("Found Url: ",found_url)
	if found_url != "" {
		 w.Write([]byte(found_url)) 
		 return
	}

	length,err := strconv.Atoi(rb.Length)
	if err != nil {
		log.Println("Could not Atoi the length")
	}
	var shortUrl string
	for {
		shortUrl, err = getUuid(int32(length))
		if err != nil {
			continue
		}
		if dbCheckIfExistsCollision(&shortUrl) == nil {
			break
		}
	}

	log.Println("short: ", shortUrl)
	dbWriteNewUrl(rb.Url, shortUrl)
	w.Write([]byte(shortUrl))

	// return short url to sender
}
func dbWriteNewUrl(url string, shortUrl string) error {
	f, err := os.OpenFile("testdb.txt",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	return nil
}


func dbCheckIfExistsUrl(url *string) (string,error) {
	dat, err := os.ReadFile("testdb.txt")
	if err != nil {
		log.Println("Error: could not find db file")
	}
	if len(dat) < 1 {
		log.Println("Database is empty")
		return "",err
	}
	db_string := string(dat)
	for line := range strings.SplitSeq(db_string, "\n") {

		content := strings.Split(line, " ")

		if len(content) > 2 || len(content) < 2 {
			continue
		}  
		db_url,sUrl := content[0],content[1]


			if db_url == *url {
				return sUrl,nil
		}
	}
	log.Print(string(dat))
	return "",err  
}
func dbCheckIfExistsCollision(shortUrl *string) error {
	dat, err := os.ReadFile("testdb.txt")
	if err != nil {
		log.Println("Error: could not find db file")
	}
	if len(dat) < 1 {
		log.Println("Database is empty")
		return nil
	}
	db_string := string(dat)
	log.Println("--------------")
	log.Println(db_string)
	log.Println("--------------")
	log.Println("db_string: ", db_string)
	for line := range strings.SplitSeq(db_string, "\n") {

		content := strings.Split(line, " ")

		if len(content) > 2 || len(content) < 2 {
			continue
		}  
			log.Println("content: ", content)
			sUrl := content[1]
			log.Println("sUrl: ",sUrl)

			log.Println("shortUrl: ",*shortUrl)

			if sUrl == *shortUrl {
				return err
		}
	}
	log.Print(string(dat))
	return nil
}

func getUuid(length int32) (string, error) {
	uuidLong := strings.ReplaceAll(uuid.New().String(), "-", "")

	if length > 32 {
		return "", errors.New("length too long")
	}
	if length != 0 {
		return uuidLong[:length], nil
	}
	return uuidLong, nil
}

func main() {

	http.HandleFunc("/create", create)
	log.Println("Running on localhost:8000")
	http.ListenAndServe(":8000", nil)
}
