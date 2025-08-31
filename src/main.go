package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ReqBody struct {
	Url    string `json:"url"`
	Length string `json:"length"`
	Name   string `json:"name"`
}

func check(db UrlDB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

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
		found_url, _ := db.CheckUrl(rb.Url)
		if found_url != "" {
			w.Write([]byte(found_url + " is your Url"))
			return
		}
	}
}
func create(db UrlDB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

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
    if !(strings.HasPrefix(rb.Url, "https://") || strings.HasPrefix(rb.Url, "http://")) {
        rb.Url = "https://" + rb.Url
    }

		found_url, _ := db.CheckUrl(rb.Url)
		log.Println("Found Url: ", found_url)
		if found_url != "" {
			w.Write([]byte(found_url))
			return
		}

		length, err := strconv.Atoi(rb.Length)
		if err != nil {
			log.Println("Could not Atoi the length assigning default")
			length = 30
		}
		if rb.Name != "" {
			result,_ := db.CheckSurl(rb.Name)
			if result == "" {
				db.WriteUrl(rb.Url, rb.Name)
				w.Write([]byte(req.Host + "/" + rb.Name))
				return
			}else {
				w.Write([]byte("already taken"))
			}
		} else {
			var shortUrl string
			for {

				id := Uuid{}
				shortUrl, err = id.Generate(int32(length))
				if err != nil {
					continue
				}
				if db.CheckForCollision(&shortUrl) == nil {
					break
				}
			}

			log.Println("short: ", shortUrl)
			db.WriteUrl(rb.Url, shortUrl)
			w.Write([]byte(shortUrl))
		}

	}
}

func home(db UrlDB) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {

		switch req.URL.Path {
		case "/":
			http.ServeFile(writer, req, "site/index.html")
		default:
			url, err := db.CheckSurl(req.URL.Path[1:])
			if err != nil {
				log.Fatal("Could not get surl from db")
			}
			if url != "" {
				log.Println("Redirecting to: " + url)
				serveSUrl(writer, req, &url)
			} else {
				http.ServeFile(writer, req, "site/index.html")
			}
		}
	}
}

func serveSUrl(writer http.ResponseWriter, req *http.Request, url *string) {
	http.Redirect(writer, req, *url, http.StatusMovedPermanently)
}

func main() {
	useText := flag.Bool("textDb",false,"Use the naive textFileDatabase")
	flag.Parse()

    var db UrlDB
    if *useText {
        db = &textFileDb{
            pathToTxt: "testdb.txt",
        }
    } else {
        sdb, err := startSqlDb()
		defer sdb.closeSqlDb()
        if err != nil {
            log.Fatal(err)
        }
        db = &sdb
    }

    http.HandleFunc("/", home(db))
    http.HandleFunc("POST /create", create(db))
    http.HandleFunc("POST /check", check(db))
	log.Println("Running on localhost:8000")
	fs := http.FileServer(http.Dir("site"))
	http.Handle("/site/", http.StripPrefix("/site/", fs))
	http.ListenAndServe(":8000", nil)
}
