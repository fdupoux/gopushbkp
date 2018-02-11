package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func archiveDirectory(rootdir string, archive *zip.Writer, excludes []string) error {

	files, err := ioutil.ReadDir(rootdir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := file.Name()
		fullpath := filepath.Join(rootdir, file.Name())

		// Check exclusions
		skipfile := false
		for _, myexclude := range excludes {
			if matched, _ := filepath.Match(myexclude, filename); matched == true {
				fmt.Printf("Excluding %s\n", fullpath)
				skipfile = true
			}
			if matched, _ := filepath.Match(myexclude, fullpath); matched == true {
				fmt.Printf("Excluding %s\n", fullpath)
				skipfile = true
			}
		}

		// Skip excluded files or directories
		if skipfile == true {
			continue
		}

		// Populate header
		header, err := zip.FileInfoHeader(file)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(fullpath, "/")
		if file.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		fmt.Printf("Archiving %s\n", fullpath)

		// Write header to archive
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		// Write contents to archive
		if file.IsDir() == true {
			err := archiveDirectory(fullpath, archive, excludes)
			if err != nil {
				return err
			}
		} else {
			file, err := os.Open(fullpath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
