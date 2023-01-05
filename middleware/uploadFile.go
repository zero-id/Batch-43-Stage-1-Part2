package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


func UploadFile(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        file, handler, err := r.FormFile("input-image")
        if err != nil {
            fmt.Println(err)
            json.NewEncoder(w).Encode("Error Retrieving the File")
            return
        }
        defer file.Close()
        fmt.Printf("Uploaded File: %+v\n", handler.Filename)

        tempFile, err := ioutil.TempFile("uploads", "image-*"+handler.Filename)
        if err != nil {
            fmt.Println(err)
            fmt.Println("path upload error")
            json.NewEncoder(w).Encode(err)
            return
        }
        defer tempFile.Close()

        fileBytes, err := ioutil.ReadAll(file)
        if err != nil {
            fmt.Println(err)
        }

        tempFile.Write(fileBytes)

        data := tempFile.Name()
        filename := data[8:]

        ctx := context.WithValue(r.Context(), "dataFile", filename)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}