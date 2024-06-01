package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileHandler struct{}

func (fh *FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	dirname := r.URL.Query().Get("dirname")
	if dirname == "" {
		dirname = `C:\xampp821\htdocs\dashcode\assets\css` // default directory
	}

	var files []FileDetails
	infos, err := os.ReadDir(dirname)
	if err != nil {
		fmt.Printf("error reading directory %v: %v\n", dirname, err)
		return
	}

	// Get the parent directory
	parentDir := filepath.Dir(dirname) + string(filepath.Separator)

	// Add the parent directory to the beginning of the files slice
	files = append([]FileDetails{{
		TID:          0,
		FileName:     filepath.Base(parentDir),
		Path:         parentDir,
		Size:         0,
		LastModified: "",
		IsDirectory:  true,
	}}, files...)

	tid := 1
	for _, info := range infos {
		// Skip hidden files
		if strings.HasPrefix(info.Name(), ".") {
			continue
		}

		fileInfo, err := info.Info()
		if err != nil {
			fmt.Printf("error getting info for file %v: %v\n", info.Name(), err)
			continue
		}

		files = append(files, FileDetails{
			TID:          tid,
			FileName:     info.Name(),
			Path:         filepath.Join(dirname, info.Name()),
			Size:         fileInfo.Size(),
			LastModified: fileInfo.ModTime().Format("02-01-2006 15:04:05"),
			IsDirectory:  info.IsDir(),
		})

		tid++
	}

	// Sort the files slice so that directories are always at the top
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDirectory != files[j].IsDirectory {
			return files[i].IsDirectory
		}
		return files[i].TID < files[j].TID
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (fh *FileHandler) CopyFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body
	var req CopyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Could not parse JSON body", http.StatusBadRequest)
		return
	}

	// Check if the source file exists
	if _, err := os.Stat(req.Src); os.IsNotExist(err) {
		http.Error(w, "Source file does not exist", http.StatusBadRequest)
		return
	}

	// Open the source file
	srcFile, err := os.Open(req.Src)
	if err != nil {
		http.Error(w, "Could not open source file", http.StatusInternalServerError)
		return
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(req.Dst)
	if err != nil {
		http.Error(w, "Could not create destination file", http.StatusInternalServerError)
		return
	}
	defer dstFile.Close()

	// Copy the file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		http.Error(w, "Could not copy file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File copied successfully"))
}

func (fh *FileHandler) ReadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body
	var req ReadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Could not parse JSON body", http.StatusBadRequest)
		return
	}

	// Check if the source file exists
	if _, err := os.Stat(req.Src); os.IsNotExist(err) {
		http.Error(w, "Source file does not exist", http.StatusBadRequest)
		return
	}

	// Read the file
	data, err := ioutil.ReadFile(req.Src)
	if err != nil {
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
