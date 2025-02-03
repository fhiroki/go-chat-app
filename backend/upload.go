package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "Successful")
}
