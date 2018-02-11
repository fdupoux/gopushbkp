/******************************************************************************\
* Copyright (C) 2018-2018 Francois Dupoux. All rights reserved.                *
*                                                                              *
* This program is free software; you can redistribute it and/or                *
* modify it under the terms of the GNU General Public                          *
* License v2 as published by the Free Software Foundation.                     *
*                                                                              *
* This program is distributed in the hope that it will be useful,              *
* but WITHOUT ANY WARRANTY; without even the implied warranty of               *
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU             *
* General Public License for more details.                                     *
*                                                                              *
* Homepage: https://github.com/fdupoux/gopushbkp                               *
\******************************************************************************/

package main

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func archiveData() (string, error) {

	// Read public key from text file
	fmt.Printf("Reading public key from %s\n", config.pubkeyfile)
	keyfile, err := os.Open(config.pubkeyfile)
	defer keyfile.Close()
	if err != nil {
		return "", fmt.Errorf("Failed to read pubkey from file: %s\n", err)
	}

	// Parse public key
	keyring, err := openpgp.ReadArmoredKeyRing(keyfile)
	if err != nil {
		return "", fmt.Errorf("openpgp.ReadArmoredKeyRing() failed: %s\n", err)
	}

	// Make sure the data exists and is a directory
	info, err := os.Stat(config.datarootdir)
	if err != nil {
		return "", fmt.Errorf("Unable to find data directory '%s': %s\n", config.datarootdir, err)
	}
	if info.IsDir() == false {
		return "", fmt.Errorf("The data root must be a directory %s\n", config.datarootdir)
	}

	// Determine full path to output archive file
	curtime := time.Now()
	filename := fmt.Sprintf("%s-%04d%02d%02d-%02d%02d%02d.zip.gpg", config.bkpbasename,
		curtime.Year(), curtime.Month(), curtime.Day(), curtime.Hour(), curtime.Minute(), curtime.Second())
	archfile := filepath.Join(config.bkparchdir, filename)
	fmt.Printf("Creating archive in %s\n", archfile)

	// Prepare checksum file
	csumwriter := sha256.New()

	// Create output encrypted file
	encfile, err := os.Create(archfile)
	if err != nil {
		return "", fmt.Errorf("Failed to create output archive file in %s: %s\n", archfile, err)
	}

	// Create multi-writer for both archive and checksum
	mwriter := io.MultiWriter(encfile, csumwriter)

	// Encrypt data using OpenPGP
	gpgio, err := openpgp.Encrypt(mwriter, keyring, nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return "", fmt.Errorf("openpgp.Encrypt() failed: %s\n", err)
	}

	// The archived data will be sent to the GPG Writer
	archive := zip.NewWriter(gpgio)

	// Call recursive function to perform the archiving
	err = archiveDirectory(config.datarootdir, archive, config.excludes)
	if err != nil {
		archive.Close()
		gpgio.Close()
		fmt.Printf("Removing archive file %s\n", archfile)
		os.Remove(archfile)
		return "", fmt.Errorf("archiveDirectory() failed: %s\n", err)
	}

	// Close archive writer
	err = archive.Close()
	if err != nil {
		return "", fmt.Errorf("Failed to close archive: %s\n", err)
	}

	// Close OpenPGP writer
	err = gpgio.Close()
	if err != nil {
		return "", fmt.Errorf("Failed to close gpgio: %s\n", err)
	}

	// Close file writer
	err = encfile.Close()
	if err != nil {
		return "", fmt.Errorf("Failed to close archive file: %s\n", err)
	}

	// Write checksum to file
	finalsum := csumwriter.Sum(nil)
	suminfo := []byte(fmt.Sprintf("%x  %s\n", finalsum, filename))
	err = ioutil.WriteFile(archfile+".sha256", suminfo, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write checksum to file: %s\n", err)
	}

	return archfile, nil
}
