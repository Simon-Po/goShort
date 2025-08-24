package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"strconv"
	"github.com/google/uuid"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

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
		fmt.Println("Error: Could not Unmarshal Body")
		http.Error(w, err.Error(), 500)
		return
	}
	found_url,_ := dbCheckIfExistsUrl(&rb.Url)
	fmt.Println("Found Url: ",found_url)
	if found_url != "" {
		 w.Write([]byte(found_url)) 
		 return
	}

	length,err := strconv.Atoi(rb.Length)
	if err != nil {
		fmt.Println("Could not Atoi the length")
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

	fmt.Println("short: ", shortUrl)
	dbWriteNewUrl(rb.Url, shortUrl)
	w.Write([]byte(shortUrl))

	// return short url to sender
}
func dbWriteNewUrl(url string, shortUrl string) error {
	f, err := os.OpenFile("testdb.txt",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("Testdb could not be opened")
		return err
	}

	input := url + " " + shortUrl + "\n"
	_, err = f.Write([]byte(input))
	if err != nil {
		fmt.Println(err)
		fmt.Println("could not write to Testdb")
		return err
	}

	return nil
}


func dbCheckIfExistsUrl(url *string) (string,error) {
	dat, err := os.ReadFile("testdb.txt")
	if err != nil {
		fmt.Println("Error: could not find db file")
	}
	if len(dat) < 1 {
		fmt.Println("Database is empty")
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
	fmt.Print(string(dat))
	return "",err  
}
func dbCheckIfExistsCollision(shortUrl *string) error {
	dat, err := os.ReadFile("testdb.txt")
	if err != nil {
		fmt.Println("Error: could not find db file")
	}
	if len(dat) < 1 {
		fmt.Println("Database is empty")
		return nil
	}
	db_string := string(dat)
	fmt.Println("--------------")
	fmt.Println(db_string)
	fmt.Println("--------------")
	fmt.Println("db_string: ", db_string)
	for line := range strings.SplitSeq(db_string, "\n") {

		content := strings.Split(line, " ")

		if len(content) > 2 || len(content) < 2 {
			continue
		}  
			fmt.Println("content: ", content)
			sUrl := content[1]
			fmt.Println("sUrl: ",sUrl)

			fmt.Println("shortUrl: ",*shortUrl)

			if sUrl == *shortUrl {
				return err
		}
	}
	fmt.Print(string(dat))
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

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/create", create)
	fmt.Println("Running on localhost:8000")
	http.ListenAndServe(":8000", nil)
}
