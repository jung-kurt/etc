package util

import (
	"fmt"
	"io"
	"os"
)

// FileCopy copies the source file identified by srcFileStr to the destination
// file identified by dstFileStr. If the destination file exists it is
// overwritten, otherwise it is created. Permission bits are carried over from
// the source file.
func FileCopy(srcFileStr, dstFileStr string) (err error) {
	var info os.FileInfo
	info, err = os.Stat(srcFileStr)
	if err == nil {
		if info.Mode().IsRegular() {
			var srcFile *os.File
			srcFile, err = os.Open(srcFileStr)
			if err == nil {
				var dstFile *os.File
				dstFile, err = os.OpenFile(dstFileStr, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
				if err == nil {
					_, err = io.Copy(dstFile, srcFile)
					dstFile.Close()
					if err == nil {
						// touch file with source timestamp but ignore error if any
						tm := info.ModTime()
						os.Chtimes(dstFileStr, tm, tm)
					}
				}
				srcFile.Close()
			}
		} else {
			err = fmt.Errorf("%s is not a regular file", srcFileStr)
		}
	}
	return
}
