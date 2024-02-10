package service

import (
	"context"
	"fmt"
	"log"
	"text-embedding-sample/dao"

	"github.com/tmc/langchaingo/llms/ollama"
)

const modelName = "dolphin-mistral"

// GenerateTexEmbebeding takes in array of strings and returns array of embeddings
func GenerateTextEmbbeding(texts []string) [][]float32 {

	llm, err := ollama.New(ollama.WithModel(modelName))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	embs, err := llm.CreateEmbedding(ctx, texts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got %d embedding(s)", len(embs))
	for i, emb := range embs {
		log.Printf("%d: len=%d; first few=%v\n", i, len(emb), emb[:4])
	}

	return embs
}

//	 AnswerQuestion takes in a question string along with DB connection and provides
//			the answer to the given question and references and similarity score with the reference chunk
func AnswerQuestion(question string, d *DBConn) string {
	log.Println("Answering question...")

	textIds := []int{}

	// get the embedding for the question
	questionEmb := GenerateTextEmbbeding([]string{question})

	// find the closest match from the database
	contextInfo, closestMatches := FindClosestMatch(questionEmb[0], d)

	for _, match := range closestMatches {
		log.Printf("MatchId: %d, Text: %v", match.TextDataId, match.Text)
		textIds = append(textIds, match.TextDataId)
	}

	log.Printf("matching  text ids are %+v\n", textIds)
	// find similarity score
	d.GetCosineSimilarity(questionEmb[0], textIds)

	// Build the prompt and execute the LLM API.
	query := fmt.Sprintf(`Use the below information to answer the subsequent question.
Information:
%v

Question: %v`, contextInfo, question)

	log.Println("Query:\n", query)

	// create OLLAMA client
	llm, err := ollama.New(ollama.WithModel(modelName))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	resp, err := llm.Call(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Got response, ID:", resp)
	return resp
}

// FindClosestMatch  finds the closest match from the database based on the cosine similarity.
func FindClosestMatch(questionEmb []float32, d *DBConn) (string, []dao.MatchingTextRecord) {
	var searchLimit = 2
	// find the closest match from the database
	closestMatches := d.SelectNearestTextEmbbedingRecord(questionEmb, searchLimit)

	// get the context info from the database
	contextInfo := GetContextInfo(closestMatches)

	return contextInfo, closestMatches
}

// GetContextInfo generates the text to be sent as query context
func GetContextInfo(closestMatches []dao.MatchingTextRecord) string {
	contextText := ""
	for _, match := range closestMatches {
		contextText = contextText + " " + match.Text
	}

	return contextText
}
