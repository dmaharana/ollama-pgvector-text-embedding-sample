package service

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"text-embedding-sample/dao"

	"github.com/pgvector/pgvector-go"
)

// read text files from a directory
// parse them and create chunks
// create embeddings and  store in the database

// get files from a directory
func GetFilesFromDirectory(directory string) ([]string, error) {
	log.Printf("Entering GetFilesFromDirectory")
	files := make([]string, 0)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("Failed to retrieve files: " + err.Error())
	}
	return files, nil
}

// create chunks of a file and save into  the database, return new & existing  chunk ids
func CreateChunksFromFileAndSaveToDb(filename string, db *DBConn) ([]dao.TextData, []dao.TextData, error) {
	log.Printf("Entering CreateChunksFromFile")

	fileType, err := FileType(filename)
	if err != nil {
		return nil, nil, err
	}

	var chunks []string

	log.Println("file  type is ", fileType)
	switch fileType {
	case "text/plain; charset=utf-8":
		chunks = breakToChunks(filename)
	default:
		return nil, nil, errors.New("File type not supported: " + fileType)
	}

	textChunks := []dao.TextData{}

	for i, c := range chunks {
		textChunks = append(textChunks, dao.TextData{
			FilePath: filename,
			ChunkId:  i,
			Text:     c,
		})
	}

	nc := db.SelectNewText(textChunks)
	oc := db.SelectOldText(textChunks)

	ids := db.SaveAllTextData(nc)
	if len(ids) != len(nc) {
		log.Fatal("Error  saving new chunks from ", filename)
	} else {
		log.Println("Save data ids:", ids)
		for i, id := range ids {
			nc[i].Id = id
		}
	}

	return nc, oc, nil
}

func SaveEmbedding(txtDocs []dao.TextData, db *DBConn) {
	// insert embedding records for new text records
	for _, v := range txtDocs {
		emb := GenerateTextEmbbeding([]string{v.Text})
		log.Printf("Text: %s, Embedding: %v\n", v.Text, emb[0][:5])
		vRec := dao.TextEmbbeding{
			TextDataId: v.Id,
			Embbeding:  pgvector.NewVector(emb[0]),
		}

		db.InsertTextEmbbeddingRecord(vRec)
		log.Println("Inserted record:", vRec.TextDataId)
	}
}

func ParseFilesSaveEmbedding(dataDir string, db *DBConn) {
	fileList, _ := GetFilesFromDirectory(dataDir)
	if len(fileList) == 0 {
		log.Println(dataDir + " is empty.")
		return
	}
	log.Println("found files: ", fileList)

	for _, fname := range fileList {
		nc, _, _ := CreateChunksFromFileAndSaveToDb(fname, db)
		if len(nc) == 0 {
			continue
		}
		SaveEmbedding(nc, db)
	}
}
