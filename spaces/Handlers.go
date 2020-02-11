package main

import (
	"fmt"
	"net/http"
	"sync"
	"os"
	"time"
	"encoding/json"
	"github.com/gorilla/mux"
)

// Services
func Files(w http.ResponseWriter, r *http.Request) {

	fmt.Println("[%s] - Request from %s ", time.Now().Format(time.RFC3339), r.RemoteAddr)

	// Check unauthorized. Replace this Authorization token by a valid one
	// by automatic generation and / or a new and dedicated web service
	/*
		if r.Header.Get("Token") != "Jkd855c6x9Aqcf" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusForbidden)
			panic("Non authorized access detected")
		}
	*/

	//vars := mux.Vars(r)

	
	vars := mux.Vars(r)
	id := vars["id"]

	concretePath :="E:\\Projects\\ProjectFIles\\private\\" + id + "\\myspace\\"

	qChannel := make(chan FileModel)
	var wg sync.WaitGroup
	wg.Add(1)
	go grDiscover(concretePath, "", qChannel, &wg)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(<-qChannel); err != nil {
		panic(err)
	}


	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func CreateSpace(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var post map[string]interface{}
	err := decoder.Decode(&post)

	if err != nil {
		panic(err)
	}

	id := post["id"].(string);
	os.MkdirAll("E:\\Projects\\ProjectFiles\\private\\" + id + "\\myspace\\", 0755)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}


func Download(w http.ResponseWriter, r *http.Request) {
/*
	vars := mux.Vars(r)
	Module := vars["file"]

	log.Printf("Downloading module %s ", Module)
	h := md5.New()
	io.WriteString(h, Module)

	out, err := os.Create("E:\\Projects\\ProjectFIles\\private\\" + fmt.Sprintf("%x", h.Sum(nil)) + "\\" + Module + ".dll")
	if err != nil {
		panic(err)
	}

	defer out.Close()

	w.WriteHeader(http.StatusOK)

	FileHeader := make([]byte, 512)
	out.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := out.Stat()                          //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+Module+".dll")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	out.Seek(0, 0)
	io.Copy(w, out) //'Copy' the file to the client
	//io.Copy(out, resp.Body)
	*/
}