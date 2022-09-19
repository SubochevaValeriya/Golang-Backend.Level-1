package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//1. Добавить в пример с файловым сервером возможность получить список всех файлов
//на сервере (имя, расширение, размер в байтах)
//2. С помощью query-параметра, реализовать фильтрацию выводимого списка по
//расширению (то есть, выводить только .png файлы, или только .jpeg)
//3. *Текущая реализация сервера не позволяет хранить несколько файлов с одинаковым
//названием (т.к. они будут храниться в одной директории на диске). Подумайте, как
//можно обойти это ограничение?
//4. К коду, написанному в рамках заданий 1-3, добавьте тесты с использованием
//библиотеки httptest

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

const UploadDir = "upload"

func main() {
	uploadHandler := &UploadHandler{
		UploadDir: UploadDir,
	}
	//dirToServe := http.Dir(uploadHandler.UploadDir)
	fs := &http.Server{
		Addr: ":8000",
		//Handler:      http.FileServer(dirToServe),
		//Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.Handle("/upload", uploadHandler)
	http.HandleFunc("/files", files)

	fs.ListenAndServe()
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("kk")
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	filePath := h.UploadDir + "/" + header.Filename
	err = os.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fileLink := h.HostAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)
}

func files(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		requiredExt := r.FormValue("extension")
		var fileList []File

		// walk the file tree
		filepath.WalkDir(UploadDir, func(s string, d fs.DirEntry, e error) error {
			fileInfo, err := d.Info()
			if err != nil {
				http.Error(w, "Unable to receive fileInfo", http.StatusBadRequest)
			}
			currentFile := File{Name: d.Name(), Extension: filepath.Ext(d.Name()), Size: fileInfo.Size()}
			fileList = append(fileList, currentFile)
			return nil
		})

		fileList = fileList[1:] // not including parent directory

		// marshal to JSON
		jsonList, err := json.Marshal(fileList)

		if err != nil {
			http.Error(w, "Unable to marshal", http.StatusInternalServerError)
		}

		fmt.Fprintln(w, string(jsonList))
	}
}

type File struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}
