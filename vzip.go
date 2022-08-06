package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func main() {
	// Read the first argument as a file name or folder name
	// and create a zip file with the same name
	if len(os.Args) != 2 {
		fmt.Println("Usage: vzip <file_name>")
		return
	}
	name := os.Args[1]

	// Create a zip file with the same name
	// and write the contents of the folder to it
	finalArchive, err := os.Create(name + ".zip")
	if err != nil {
		panic(err)
	}

	zipWriter := zip.NewWriter(finalArchive)

	defer finalArchive.Close()
	defer zipWriter.Close()

	// If the chosen argument is a file
	if isFile(name) {
		addSingleFileToZip(name, zipWriter)
		return
	}

	// If the chosen argument is a folder
	files := readDir(name)
	doFileLoop(name, files, zipWriter)
}

func addSingleFileToZip(dir string, zipWriter *zip.Writer) {
	fileWriter, err := os.Open(dir)
	if err != nil {
		panic(err)
	}

	defer fileWriter.Close()
	createStream, err := zipWriter.Create(dir)
	if err != nil {
		panic(err)
	}

	// Copy file content to the zip file
	copyToFinalZip(createStream, fileWriter)
}

func loopInsideFolder(dir string, zipWriter *zip.Writer) {
	fileWriter, err := os.Open(dir)
	if err != nil {
		panic(err)
	}

	defer fileWriter.Close()

	if isFile(dir) {
		createStream, err := zipWriter.Create(dir)
		if err != nil {
			panic(err)
		}
	
		// Copy file content to the zip file
		copyToFinalZip(createStream, fileWriter)
	} else {
		files := readDir(dir)
		if len(files) != 0 {
			doFileLoop(dir, files, zipWriter)
		}
	}
}

func doFileLoop(dir string, files []fs.DirEntry, zipWriter *zip.Writer) {
	for _, file := range files {
		if !file.IsDir() {
			addSingleFileToZip(dir+"/"+file.Name(), zipWriter)
			return
		}

		loopInsideFolder(dir+"/"+file.Name(), zipWriter)
	}
}

func copyToFinalZip(dest io.Writer, src *os.File) {
	_, err := io.Copy(dest, src)
	if err != nil {
		panic(err)
	}
}

func readDir(name string) []fs.DirEntry {
	files, err := os.ReadDir(name)
	if err != nil {
		return []fs.DirEntry{}
	}
	return files
}

func isFile(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}
