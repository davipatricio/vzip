package main

import (
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

var options *flag.FlagSet
var compressionLevel *int
var compressMethod *string
var dest *string

func main() {
	// Read the first argument as a file name or folder name
	// and create a zip file with the same name
	if len(os.Args) < 2 {
		fmt.Println("Usage: vzip <folder_name | file_name> [--dest] [--level=int] [--method=none,gzip,zlib]\n  --dest: Where to store the zip file.\nDefaults to ./<filename>.zip\n  --level: Whether to compress the file or not. Between -1 and 9. The higher the number, better the compression.\nDefaults to 0 (disables compression).\n  --method: The method to use for compression.\nAccepted values: (default) none, gzip, zlib")
		return
	}

	name := os.Args[1]

	if !fileExists(name) {
		fmt.Println("File or folder with name " + name + " does not exist.")
		os.Exit(1)
		return
	}

	options = flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	compressionLevel = options.Int("level", 0, "Whether to compress the file or not. Between -1 and 9. The higher the number, better the compression.\n(default) 0 disables compression.")
	compressMethod = options.String("method", "none", "The method to use for compression.\nAccepted values: (default) none, gzip, zlib")
	dest = options.String("dest", "", "Where to store the zip file.\n(default) ./<filename>.zip")

	options.Parse(os.Args[2:])

	if *dest == "" {
		*dest = "./" + name + ".zip"
	}

	if !strings.HasSuffix(*dest, ".zip") {
		*dest += ".zip"
	}

	// Create a zip file with the same name
	// and write the contents of the folder to it
	finalArchive, err := os.Create(*dest)
	if err != nil {
		panic(err)
	}

	zipWriter := zip.NewWriter(finalArchive)
	addCompressorToWriter(zipWriter)

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

func addCompressorToWriter(dst *zip.Writer) {
	if *compressionLevel < -1 || *compressionLevel > 9 {
		fmt.Println("Compression level must be between -1 and 9.")
		os.Exit(1)
		return
	}

	switch *compressMethod {
	case "gzip":
		dst.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriterLevel(out, *compressionLevel)
		})
	case "zlib":
		dst.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return zlib.NewWriterLevel(out, *compressionLevel)
		})
	case "none":
		dst.RegisterCompressor(zip.Store, func(out io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(out, flate.NoCompression)
		})
		dst.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(out, flate.NoCompression)
		})
	default:
		fmt.Println("Invalid compression method.")
		os.Exit(1)
	}
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

func fileExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}

func isFile(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}
