package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v4"
	"log"
	"net/http"
	"os"
)

const PORT = ":8080"

func archive(filename string) error {
	files, mapErr := archiver.FilesFromDisk(nil, map[string]string{
		filename: filename,
	})
	if mapErr != nil {
		return mapErr
	}

	// create the output file we'll write to
	out, fileErr := os.Create("archive.tar.gz")
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

func main() {
	router := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	//router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.POST("/upload", func(ctx *gin.Context) {
		// Extract file from request
		file, formErr := ctx.FormFile("file")
		if formErr != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("File not provided."))
			return
		}
		log.Println("Received file: " + file.Filename)

		// Archive file
		archiveErr := archive(file.Filename)
		if archiveErr != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Something went wrong during archiving."))
			return
		}
		ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

		// Upload the file to specific dst.
		//c.SaveUploadedFile(file, dst)

	})

	if router.Run(PORT) != nil {
		log.Println("Error starting server.")
	}
}
