package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const uploadDir = "./uploads/" // Directory where uploaded files are stored

func uploadChunkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Request received")
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get file from the request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get the chunk number and total chunks from the request
	chunk, err := strconv.Atoi(r.FormValue("chunk"))
	if err != nil {
		http.Error(w, "Invalid chunk number", http.StatusBadRequest)
		return
	}
	totalChunks, err := strconv.Atoi(r.FormValue("totalChunks"))
	if err != nil {
		http.Error(w, "Invalid total chunks", http.StatusBadRequest)
		return
	}

	// Create a temporary file for storing the chunk
	chunkFileName := fmt.Sprintf("%sfile_chunk_%d", uploadDir, chunk)
	chunkFile, err := os.Create(chunkFileName)
	if err != nil {
		http.Error(w, "Failed to create chunk file", http.StatusInternalServerError)
		return
	}
	defer chunkFile.Close()

	// Write the uploaded chunk to the temporary file
	_, err = io.Copy(chunkFile, file)
	if err != nil {
		http.Error(w, "Failed to write chunk to file", http.StatusInternalServerError)
		return
	}

	// If this is the last chunk, reassemble the file
	if chunk == totalChunks-1 {
		finalFile, err := os.Create(uploadDir + "final_file")
		if err != nil {
			http.Error(w, "Failed to create final file", http.StatusInternalServerError)
			return
		}
		defer finalFile.Close()

		// Reassemble the chunks into the final file
		for i := 0; i < totalChunks; i++ {
			chunkFileName := fmt.Sprintf("%sfile_chunk_%d", uploadDir, i)
			chunkFile, err := os.Open(chunkFileName)
			if err != nil {
				http.Error(w, "Failed to open chunk file", http.StatusInternalServerError)
				return
			}
			_, _ = io.Copy(finalFile, chunkFile)
			chunkFile.Close()
			os.Remove(chunkFileName) // Delete the chunk file after reassembly
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Chunk %d uploaded successfully", chunk)
}

func main() {
	// Ensure the upload directory exists
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	http.HandleFunc("/upload-chunk", uploadChunkHandler)
	fmt.Println("Server is listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
