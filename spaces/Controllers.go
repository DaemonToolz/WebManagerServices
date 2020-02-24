package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

/*
	-------------------
*/

// Goroutines related

func ReadDir(DirPath string, parent string) []FileModel {
	files, err := ioutil.ReadDir(DirPath)
	myFiles := make([]FileModel, 0)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return myFiles
	}

	for _, f := range files {
		if f.IsDir() == false {
			currentFile := FileModel{
				Id:   "",
				Name: f.Name(),
				Path: (DirPath + "/" + f.Name()),
				Size: f.Size(),
				Type: 0,
			}

			myFiles = append(myFiles, currentFile)
		} else {
			myFiles = append(myFiles, ReadDir(DirPath+"/"+f.Name(), parent+"/"+f.Name())...)
		}
	}

	return myFiles
}

// Goroutine discovery
func grDiscover(DirPath string, parent string, result chan FileModel, wg *sync.WaitGroup) {
	defer wg.Done() // Done at the end, ofc

	log.Printf("Goroutine %s started", DirPath)

	files, err := ioutil.ReadDir(DirPath)
	if err != nil {
		log.Println("Error opening file:", err)
		result <- (FileModel{})

		return // Do not continue
	}

	// Preparing space for routines
	totalLocks := 0
	qChannel := make(chan FileModel, len(files)+20) // Margin
	var thisWaitingGroup sync.WaitGroup

	for _, f := range files {
		log.Printf("Element %s found", f.Name())
		if f.IsDir() == false {

			currentFile := FileModel{
				Id:   "",
				Name: f.Name(),
				Path: (DirPath + "/" + f.Name()),
				Size: f.Size(),
				Type: 0,
			}

			result <- currentFile
		} else {
			thisWaitingGroup.Add(1)
			totalLocks++
			go grDiscover(DirPath+"/"+f.Name(), parent+"/"+f.Name(), qChannel, &thisWaitingGroup)
		}
	}

	log.Printf("Waiting nested goroutines of %s", DirPath)
	thisWaitingGroup.Wait()

	log.Printf("Done waiting in %s", DirPath)

	close(qChannel)
	log.Printf("Routine Locks: %d", totalLocks)
	if len(qChannel) > 0 {
		log.Printf("Total responses: %d", len(qChannel))

		for response := range qChannel {

			log.Println("Response acquired: ", response)
			result <- response
		}

		log.Println("Channel draining done")

	}

	log.Printf("Goroutine %s terminated", DirPath)
	close(result)
	//wg.Done()
	//result <- query
}

// Goroutine discovery
func grDiscoverFiles(DirPath string, parent string, result chan FileModel, wg *sync.WaitGroup) {
	defer wg.Done() // Done at the end, ofc

	log.Printf("Goroutine %s started", DirPath)

	files, err := ioutil.ReadDir(DirPath)
	if err != nil {
		log.Println("Error opening file:", err)
		result <- (FileModel{})

		return // Do not continue
	}

	for _, f := range files {
		currentFile := FileModel{
			Id:   "",
			Name: f.Name(),
			Path: (DirPath + "/" + f.Name()),
			Size: f.Size(),
		}

		if f.IsDir() {
			currentFile.Type = 1
		} else {
			currentFile.Type = 0
		}
		result <- currentFile
	}

	log.Printf("Goroutine %s terminated", DirPath)
	close(result)

}

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func WriteToFile(path string, content interface{}) error {
	jsonContent, _ := json.MarshalIndent(content, "", " ")
	err := ioutil.WriteFile(path, jsonContent, 0644)
	failOnError(err, "An error occured while writing to a file")
	return err
}

func getUsers() []string {

	file, err := os.Open(getUsersFolder())
	if err != nil {
		failOnError(err, "Unable to read directory")
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	var users []string
	for _, name := range list {
		users = append(users, name)
	}

	return users
}
