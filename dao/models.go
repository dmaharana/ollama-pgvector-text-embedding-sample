package dao

import "github.com/pgvector/pgvector-go"

type TextData struct {
	Id        int    `json:"id"`
	ChunkId   int    `json:"chunk_id"`
	FilePath  string `json:"file_path"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type TextEmbbeding struct {
	Id         int             `json:"id"`
	TextDataId int             `json:"text_data_id"`
	Embbeding  pgvector.Vector `json:"embededing"`
	CreatedAt  string          `json:"created_at"`
}

type MatchingTextRecord struct {
	TextDataId int             `json:"text_data_id"`
	ChunkId    int             `json:"chunk_id"`
	FilePath   string          `json:"file_path"`
	Text       string          `json:"text"`
	Embbeding  pgvector.Vector `json:"embededing"`
}

type SimilarityScore struct {
	TextDataId int     `json:"text_data_id"`
	Similarity float64 `json:"similarity"`
}
