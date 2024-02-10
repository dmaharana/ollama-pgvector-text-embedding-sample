package main

import (
	"flag"
	"log"
	"os"
	"text-embedding-sample/service"
)

func main() {
	log.SetFlags(log.Lshortfile)

	db := service.CreateDbConnection()
	defer db.Conn.Close()

	db.CreateAllTables()

	var dataDir, query string

	flag.StringVar(&dataDir, "l", "", "Directory with text data for loading into the database")

	flag.StringVar(&query, "q", "", "Question  to be processed by the model.")
	flag.Parse()

	if dataDir == "" && query == "" {
		flag.PrintDefaults()
		log.Println("Please provide either -l or -q parameter.")
		os.Exit(1)
	}

	if dataDir != "" {
		service.ParseFilesSaveEmbedding(dataDir, db)
	} else if query != "" {
		service.AnswerQuestion(query, db)
	}

}
