package main

import (
	_ "bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	_ "log"
	_ "os"
	"strings"

	"github.com/chewxy/sexp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/niklasfasching/go-org/org"
	_ "github.com/niklasfasching/go-org/org"
)

type OrgFile struct {
	title string
	path  string
}

func connect() []OrgFile {
	dbLoc := "/Users/raheel/.emacs.d/org-roam.db"
	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var sexps sexp.Sexp
	done := make(chan struct{})
	var parser = sexp.NewParser(strings.NewReader("(hello goodbye)"), false)
	go func(ch chan sexp.Sexp, done chan struct{}) {
		for s := range ch {
			sexps = s
		}

		done <- struct{}{}
	}(parser.Output, done)
	parser.Run()
	<-done
	fmt.Printf("%d\n", sexps.LeafCount())

	rows, err := db.Query(`select * from titles`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	files := make([]OrgFile, 10)
	for rows.Next() {
		var file string
		var title string
		err = rows.Scan(&file, &title)
		if err != nil {
			log.Fatal(err)
		}
		file = strings.TrimPrefix(file, "\"")
		file = strings.TrimSuffix(file, "\"")

		files = append(files, OrgFile{title, file})
	}
	return files
}

func main() {
	files := connect()
	base_path := "/Users/raheel/Downloads/org-roam-export"
	for _, file := range files {
		path := file.path

		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Error with  %v\n", err)
			continue
		}
		fmt.Printf("read %d bytes\n", len(data))
		input := strings.NewReader(string(data)).read
		html, err := org.New().Parse(input, "./").Write(org.NewHTMLWriter())
		if err != nil {
			fmt.Printf("Error with %s:  %v\n", file.title, err)
			continue
		}

		err = ioutil.WriteFile(base_path+"/"+file.title+".html", []byte(html), 0666)
		fmt.Println(err)
	}

	// f1, err := ioutil.ReadFile("/Users/raheel/orgs/roam/20201116113425-shaders.org")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(f1))
}
