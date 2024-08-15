package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v4"
	"io"
	"log"
	"net/http"
	"os"
)

const PORT = ":8080"
const ARCHIVE_NAME = "archive.tar.gz"

func archive(filename string) error {
	files, mapErr := archiver.FilesFromDisk(nil, map[string]string{
		filename: filename,
	})
	if mapErr != nil {
		return mapErr
	}

	// create the output file we'll write to
	out, fileErr := os.Create(ARCHIVE_NAME)
	if fileErr != nil {
		return fileErr
	}
	defer out.Close()

	// we can use the CompressedArchive type to gzip a tarball
	// (compression is not required; you could use Tar directly)
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}

	// create the archive
	archiveErr := format.Archive(context.Background(), out, files)
	if archiveErr != nil {
		return archiveErr
	}
	return nil
}

func removeFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("FILE_REMOVING_ERROR:", err)
	}
}

func processFile(ctx *gin.Context) {
	// Extract file from request
	file, formErr := ctx.FormFile("file")
	if formErr != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("File not provided."))
		return
	}
	log.Println("Received file: " + file.Filename)

	dst, fileErr := os.Create(file.Filename)
	if fileErr != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Something went wrong."))
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Something went wrong.")
		return
	}

	if _, err := io.Copy(dst, src); err != nil {
		ctx.String(http.StatusInternalServerError, "Something went wrong.")
		return
	}

	// Archive file
	archiveErr := archive(file.Filename)
	if archiveErr != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Something went wrong during archiving."))
		return
	}
	defer removeFile(ARCHIVE_NAME)
	defer removeFile(file.Filename)
	defer src.Close()
	defer dst.Close()

	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func main() {
	router := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	//router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.POST("/upload", processFile)

	if router.Run(PORT) != nil {
		log.Println("Error starting server.")
	}
}
