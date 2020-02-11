package main

/*
	--------------- MODELS
*/

type FileModel struct {
	id   string `json:"id"`
	name string `json:"name"`
	path string `json:"path"`
	size int    `json:"size"`
}
