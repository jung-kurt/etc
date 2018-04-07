package util

// helper for the archive/zip package (archive/unarchive to/from a file/reader/writer)
// Cleaned up from https://github.com/pierrre/archivefile (MIT license)

import (
	zip_impl "archive/zip"
	// "fmt"
	"io"
	"os"
	//	"path"
	"path/filepath"
	"strings"
)

// Archive compresses a file/directory to a writer
//
// If the path ends with a separator, then the contents of the folder at that path
// are at the root level of the archive, otherwise, the root of the archive contains
// the folder as its only item (with contents inside).
//
// If progress is not nil, it is called for each file added to the archive.
func Archive(inFilePath string, writer io.Writer, progress ProgressFunc) (err error) {
	var zipWriter *zip_impl.Writer
	var zipFileWriter io.Writer
	var zipHeader *zip_impl.FileHeader
	var archivePath, relativeFilePath, basePath string
	var file *os.File

	zipWriter = zip_impl.NewWriter(writer)

	basePath = filepath.Dir(inFilePath)

	err = filepath.Walk(inFilePath, func(filePath string, fileInfo os.FileInfo, walkErr error) (err error) {
		if walkErr == nil {
			if !fileInfo.IsDir() {

				relativeFilePath, err = filepath.Rel(basePath, filePath)
				if err == nil {

					archivePath = filepath.ToSlash(relativeFilePath)

					if progress != nil {
						progress(archivePath)
					}

					zipHeader, err = zip_impl.FileInfoHeader(fileInfo)
					if err == nil {
						zipHeader.Name = archivePath
						zipHeader.Method = zip_impl.Deflate
						// fmt.Printf("archive path [%s]\n", archivePath)
						zipFileWriter, err = zipWriter.CreateHeader(zipHeader)
						if err == nil {
							file, err = os.Open(filePath)
							if err == nil {
								_, err = io.Copy(zipFileWriter, file)
								file.Close()
							}
						}
					}
				}
			}
		} else {
			err = walkErr
		}
		return
	})
	if err != nil {
		return err
	}

	return zipWriter.Close()
}

// ArchiveFile compresses a file/directory to a file
//
// See Archive() doc
func ArchiveFile(inFilePath string, outFilePath string, progress ProgressFunc) (err error) {
	var outFile *os.File

	outFile, err = os.Create(outFilePath)
	if err == nil {
		err = Archive(inFilePath, outFile, progress)
		outFile.Sync()
		outFile.Close()
	}
	return
}

// Unarchive decompresses a reader to a directory
//
// The data's size is required because the zip reader needs it.
//
// The archive's content will be extracted directly to outFilePath.
//
// If progress is not nil, it is called for each file extracted from the archive.
func Unarchive(reader io.ReaderAt, readerSize int64, outFilePath string, progress ProgressFunc) (err error) {
	var zipReader *zip_impl.Reader
	var j int

	zipReader, err = zip_impl.NewReader(reader, readerSize)
	if err == nil {
		for j = 0; j < len(zipReader.File) && err == nil; j++ {
			err = unarchiveFile(zipReader.File[j], outFilePath, progress)
		}

	}
	return
}

// UnarchiveFile decompresses a file to a directory
//
// See Unarchive() doc
func UnarchiveFile(inFilePath string, outFilePath string, progress ProgressFunc) (err error) {
	var inFile *os.File
	var inFileInfo os.FileInfo

	inFile, err = os.Open(inFilePath)
	if err == nil {
		inFileInfo, err = inFile.Stat()
		if err == nil {
			err = Unarchive(inFile, inFileInfo.Size(), outFilePath, progress)
		}
		inFile.Close()
	}
	return
}

func unarchiveFile(zipFile *zip_impl.File, outFilePath string, progress ProgressFunc) (err error) {
	var zipFileReader io.ReadCloser
	var filePath string
	var file *os.File

	if !zipFile.FileInfo().IsDir() {

		if progress != nil {
			progress(zipFile.Name)
		}

		zipFileReader, err = zipFile.Open()
		if err == nil {

			filePath = filepath.Join(outFilePath, filepath.Join(strings.Split(zipFile.Name, "/")...))

			err = os.MkdirAll(filepath.Dir(filePath), os.FileMode(0755))
			if err == nil {

				file, err = os.Create(filePath)
				if err == nil {
					_, err = io.Copy(file, zipFileReader)
					file.Close()
				}
			}
			zipFileReader.Close()
		}
	}
	return
}

// ProgressFunc is the type of the function called for each archive file.
type ProgressFunc func(archivePath string)
