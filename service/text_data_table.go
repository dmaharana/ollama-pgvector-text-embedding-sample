package service

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"text-embedding-sample/dao"

	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

const EmbeddingTokenSize = 4096
const MaxInsertParams = 65535

// create test data table

func (db *DBConn) CreateTextDataTable() error {
	log.Println("Create text data table")
	_, err := db.Conn.Exec(db.Ctx, dao.CreateTextDataStmt)
	if err != nil {
		panic(err)
	}

	return nil
}

func (db *DBConn) CreateVectorDataTable() error {
	log.Println("Create vector data table")
	sqlStmt := fmt.Sprintf(dao.CreateTextEmbbedingStmt, EmbeddingTokenSize)
	_, err := db.Conn.Exec(db.Ctx, sqlStmt)
	if err != nil {
		panic(err)
	}

	return nil
}

func (db *DBConn) SaveAllTextData(textData []dao.TextData) []int {
	log.Println("Save all text data")
	insertedRows := []int{}

	// if no records then return
	if len(textData) == 0 {
		return insertedRows
	}

	log.Printf("Will be saving %d records", len(textData))

	m := dao.TextData{}
	t := reflect.TypeOf(m)
	colCount := t.NumField()
	// remove id and created_at
	colCount -= 2

	batchSize := MaxInsertParams / colCount
	log.Printf("Col count: %d, Batch size: %d", colCount, batchSize)

	for i := 0; i < len(textData); i += batchSize {
		end := i + batchSize
		if end > len(textData) {
			end = len(textData)
		}
		batch := textData[i:end]

		// log.Println("Inserting batch", batch)

		valStr := make([]string, 0, len(batch)*colCount)
		valArgs := make([]interface{}, 0, len(batch)*colCount)

		for j, data := range batch {
			ph := make([]string, colCount)
			for k := 0; k < colCount; k++ {
				ph[k] = fmt.Sprintf("$%d", j*colCount+k+1)
			}

			valStr = append(valStr, fmt.Sprintf("(%s)", strings.Join(ph, ",")))
			valArgs = append(valArgs, data.ChunkId, data.FilePath, data.Text)
		}

		sqlStmt := fmt.Sprintf(dao.InsertTextDataStmt, strings.Join(valStr, ","))

		rows, err := db.Conn.Query(db.Ctx, sqlStmt, valArgs...)
		if err != nil {
			log.Fatalln(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				log.Fatalln(err)
			}
			insertedRows = append(insertedRows, id)
		}
	}

	log.Printf("Inserted rows: %v", insertedRows)
	log.Printf("Number of rows of text data saved: %d", len(textData))

	return insertedRows
}

func (d *DBConn) SelectNewText(data []dao.TextData) []dao.TextData {
	log.Println("Select new text")
	newTextData := []dao.TextData{}

	if len(data) == 0 {
		return newTextData
	}

	filePaths := make([]string, len(data))
	textContents := make([]string, len(data))
	chunkIds := make([]int, len(data))

	for _, c := range data {
		filePaths = append(filePaths, c.FilePath)
		textContents = append(textContents, c.Text)
		chunkIds = append(chunkIds, c.ChunkId)
	}

	rows, err := d.Conn.Query(d.Ctx, dao.SelectNewContentStmt, pq.Array(filePaths), pq.Array(textContents), pq.Array(chunkIds))
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunkId int
		var filePath string
		var text string
		err := rows.Scan(&filePath, &text, &chunkId)
		if err != nil {
			log.Fatalln(err)
		}

		if text == "" || filePath == "" {
			continue
		}
		newTextData = append(newTextData, dao.TextData{ChunkId: chunkId, FilePath: filePath, Text: text})
	}

	// log.Println("New text data:", PpObj(newTextData))
	log.Printf("Number of new text data: %d", len(newTextData))

	return newTextData
}

func (d *DBConn) SelectOldText(data []dao.TextData) []dao.TextData {
	log.Println("Select old text")
	oldTextData := []dao.TextData{}

	if len(data) == 0 {
		return oldTextData
	}

	filePaths := make([]string, len(data))
	textContents := make([]string, len(data))
	chunkIds := make([]int, len(data))

	for _, c := range data {
		filePaths = append(filePaths, c.FilePath)
		textContents = append(textContents, c.Text)
		chunkIds = append(chunkIds, c.ChunkId)
	}

	rows, err := d.Conn.Query(d.Ctx, dao.SelectOldContentStmt, pq.Array(filePaths), pq.Array(textContents), pq.Array(chunkIds))
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var chunkId int
		var filePath string
		var text string
		err := rows.Scan(&id, &filePath, &text, &chunkId)
		if err != nil {
			log.Fatalln(err)
		}
		if text == "" || filePath == "" {
			continue
		}
		oldTextData = append(oldTextData, dao.TextData{Id: id, ChunkId: chunkId, FilePath: filePath, Text: text})
	}

	log.Println("Old text data:", PpObj(oldTextData))
	log.Printf("Number of old text data: %d", len(oldTextData))
	return oldTextData
}

func (d *DBConn) InsertTextEmbbeddingRecord(data dao.TextEmbbeding) int {
	log.Println("Insert text embbeding record")
	var rowId int
	err := d.Conn.QueryRow(d.Ctx, dao.InsertTextEmbbedingStmt, data.TextDataId, data.Embbeding).Scan(&rowId)
	if err != nil {
		log.Fatalln(err)
	}

	return rowId
}

func (d *DBConn) UpdateTextEmbbedingRecord(data dao.TextEmbbeding) error {
	log.Println("Update text embbeding record")
	_, err := d.Conn.Exec(d.Ctx, dao.UpdateTextEmbbedingStmt, data.Embbeding, data.TextDataId)
	if err != nil {
		// log.Fatalln(err)
		log.Println("Update text embbeding record failed")
		return err
	}
	return nil
}

func (d *DBConn) SelectNearestTextEmbbedingRecord(qEmbed []float32, limit int) []dao.MatchingTextRecord {
	log.Println("Select text embbeding record")
	matchingRecs := []dao.MatchingTextRecord{}

	rows, err := d.Conn.Query(d.Ctx, dao.SelectNearestTextStmt, pgvector.NewVector(qEmbed), limit)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rec dao.MatchingTextRecord
		err := rows.Scan(&rec.TextDataId, &rec.ChunkId, &rec.FilePath, &rec.Text, &rec.Embbeding)
		if err != nil {
			log.Fatalln(err)
		}
		matchingRecs = append(matchingRecs, rec)
	}
	log.Println("Matching records:", PpObj(matchingRecs))

	return matchingRecs
}

func (d *DBConn) GetCosineSimilarity(qEmbed []float32, txtIds []int) []dao.SimilarityScore {

	rows, err := d.Conn.Query(d.Ctx, dao.SelectCosineSimilarityStmt, pgvector.NewVector(qEmbed), pq.Array(txtIds))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer rows.Close()

	simScores := []dao.SimilarityScore{}

	for rows.Next() {
		var txtId int
		var sim float64
		err := rows.Scan(&txtId, &sim)
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		simScores = append(simScores, dao.SimilarityScore{
			TextDataId: txtId,
			Similarity: sim,
		})
	}

	log.Println(PpObj(simScores))

	return simScores
}
