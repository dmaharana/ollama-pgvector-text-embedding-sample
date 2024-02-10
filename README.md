## This project implements RAG

Use Ollama and pgvector to create a Retrieval Augmented Generation (RAG) system. This will allow us to query about information present in our local documents, without fine-tuning the Large Language Model (LLM). When using RAG, on quering for a piece of information, we will first do a retrieval step to fetch any relevant document chunk from a vector database where the documents were indexed. This example uses "dolphin-mistral" LLM to create embeddings as well as act as a chat agent answering the query. [Ollama](https://ollama.com/) is used to serve the LLM and provides a REST interface to [ollama/ollama](https://github.com/ollama/ollama) golang module.
[pgvector/pgvector](https://github.com/pgvector/pgvector) is run as a container to serve as a vector database. SQLs are written as documented in the pgvector project to store, find closest match and also find similariy with the stored vector embeddedings.

## Build

`go build -o rag` - this will create a binary named `rag`. You can then run the program with `./rag`.

## Usage

- `./rag -l <path to documents directory>` - this option will read the documents in the directory, create chunks, create embeddings and store them to database.
- `./rag -q <"Question for the chat engine">` - this option will create embedding of the query string, find the closest match with the data and create a prompt for the LLM chat agent.
- We also get the reference to the document, chunks and closeness match to the chunks information.

## References

- [Retrieval Augmented Generation in Go](https://eli.thegreenplace.net/2023/retrieval-augmented-generation-in-go/)
- [go-rag-openai](https://github.com/eliben/code-for-blog/blob/master/2023/go-rag-openai/cmd/chunker/chunker.go)
- [VECTOR DATABASES - PGVECTOR AND LANGCHAIN](https://bugbytes.io/posts/vector-databases-pgvector-and-langchain/)
- [pgvector/pgvector](https://github.com/pgvector/pgvector)
- [pgvector/pgvector-go](https://github.com/pgvector/pgvector-go)
- [Sample Data](https://github.com/hwchase17/chat-your-data/blob/master/state_of_the_union.txt)
