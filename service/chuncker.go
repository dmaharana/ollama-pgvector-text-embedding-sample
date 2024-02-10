package service

// https://github.com/eliben/code-for-blog/blob/master/2023/go-rag-openai/cmd/chunker/chunker.go
// https://eli.thegreenplace.net/2023/retrieval-augmented-generation-in-go/

import (
	"bufio"
	"bytes"
	"os"

	"github.com/pkoukk/tiktoken-go"
)

const tokenEncoding = "cl100k_base"
const tableName = "chunks"

const chunkSize = 1000

// breakToChunks reads the file in `path` and breaks it into chunks of
// approximately chunkSize tokens each, returning the chunks.
func breakToChunks(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	tke, err := tiktoken.GetEncoding(tokenEncoding)
	if err != nil {
		panic(err)
	}

	chunks := []string{""}

	scanner := bufio.NewScanner(f)
	scanner.Split(splitByParagraph)

	for scanner.Scan() {
		chunks[len(chunks)-1] = chunks[len(chunks)-1] + scanner.Text() + "\n"
		toks := tke.Encode(chunks[len(chunks)-1], nil, nil)
		if len(toks) > chunkSize {
			chunks = append(chunks, "")
		}
	}

	// If we added a new empty chunk but there weren't any paragraphs to add to
	// it, make sure to remove it.
	if len(chunks[len(chunks)-1]) == 0 {
		chunks = chunks[:len(chunks)-1]
	}

	return chunks
}

// splitByParagraph is a custom split function for bufio.Scanner to split by
// paragraphs (text pieces separated by two newlines).
func splitByParagraph(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		return i + 2, bytes.TrimSpace(data[:i]), nil
	}

	if atEOF && len(data) != 0 {
		return len(data), bytes.TrimSpace(data), nil
	}

	return 0, nil, nil
}
