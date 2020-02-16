package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Services
func GetFile(w http.ResponseWriter, r *http.Request) {

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
	space := vars["space"]

	if strings.Contains(space, "..") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	concretePath := "E:\\Projects\\ProjectFIles\\private\\" + space + "\\myspace\\"

	qChannel := make(chan FileModel)
	var wg sync.WaitGroup
	wg.Add(1)
	go grDiscoverFiles(concretePath, "", qChannel, &wg)

	if err := json.NewEncoder(w).Encode(<-qChannel); err != nil {
		panic(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

// Services
func GetFiles(w http.ResponseWriter, r *http.Request) {

	log.Printf("[ %s ] - Request from %s ", time.Now().Format(time.RFC3339), r.RemoteAddr)

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
	space := vars["space"]

	if strings.Contains(space, "..") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	concretePath := "E:\\Projects\\ProjectFIles\\private\\" + space + "\\myspace\\"

	qChannel := make(chan FileModel)
	var wg sync.WaitGroup
	wg.Add(1)
	go grDiscoverFiles(concretePath, "", qChannel, &wg)

	var files []FileModel
	for file := range qChannel {
		files = append(files, file)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Printf(err.Error())
		panic(err)
	}

}

func CreateSpace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var post map[string]interface{}
	err := decoder.Decode(&post)

	if err != nil {
		panic(err)
	}

	id := post["id"].(string)
	userSpace := "E:\\Projects\\ProjectFiles\\private\\" + id + "\\myspace\\"
	sharedFolder := "E:\\Projects\\ProjectFiles\\shared\\"

	sendMessage("user-notification", false, constructNotification(id, "CreateSpace", STATUS_ONGOING, PRIORITY_STD, TYPE_INFO))

	os.MkdirAll(userSpace, 0755)
	err = CopyDir(sharedFolder, userSpace)

	if err != nil {
		failOnError(err, "Failed to copy a directory")
		sendMessage("user-notification", false, constructNotification(id, "CreateSpace", STATUS_ERROR, PRIORITY_CRITICAL, TYPE_INFO))

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
	} else {

		sendMessage("user-notification", false, constructNotification(id, "CreateSpace", STATUS_DONE, PRIORITY_STD, TYPE_INFO))

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
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
