package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"os"

	"github.com/google/uuid"
)




func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}

type ReqBody struct {
	Url string `json:"url"`
	Length int32 `json:"length"`
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
	// Check if url is already in db

	var shortUrl string
	for {
		shortUrl, err := getUuid(rb.Length)
    if err != nil {
        continue
    }
    if dbCheckIfExists(shortUrl) == nil {
        break
    }
	}
	dbWriteNewUrl(rb.Url,shortUrl)
	w.Write([]byte(shortUrl))

	// return short url to sender
}
func dbWriteNewUrl(url string, shortUrl string) error {
	f,err := os.Open("testdb.txt")
	defer f.Close()
	if err != nil {
		fmt.Println("Testdb could not be opened")
		return err
	}

	input := url + " " + shortUrl + "\n"
	_,err = f.Write([]byte(input))
	if err != nil {
		fmt.Println("could not write to Testdb")
		return err
	}

	return nil
}

func dbCheckIfExists(shortUrl string) error {
	    dat, err := os.ReadFile("testdb.txt")
    if err != nil {
			fmt.Println("Error: could not find db file")
		}
		db_string := string(dat)
		for line := range strings.SplitSeq(db_string, "\n")	{
			content := strings.Split(line," ")
			if(len(content) > 2) {
				fmt.Println("Error: Something is wrong with line:",strings.Join(content,""))
			}
			sUrl:= content[1]
			if sUrl == shortUrl	{
				return err
			}
		}
    fmt.Print(string(dat))
		return nil
}

func getUuid(length int32) (string, error) {
    uuidLong := strings.ReplaceAll(uuid.New().String(),"-","")

    if length > 32 {
        return "", errors.New("length too long")
    }
    if length != 0 {
        return uuidLong[:length+1], nil
    }
    return uuidLong, nil
}

func main() {

    http.HandleFunc("/hello", hello)
		http.HandleFunc("/create",create)
		fmt.Println("Running on localhost:8000")
    http.ListenAndServe(":8000", nil)
}
