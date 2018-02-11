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
	"fmt"
	"os"
)

var config configuration

func main() {

	var err error

	// Read configuration file
	config, err = readConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file: %s\n", err)
		os.Exit(1)
	}

	// Create the encrypted archive locally
	archfile, err := archiveData()
	if err != nil {
		fmt.Printf("Failed to create backup archive: %s\n", err)
		os.Exit(2)
	}

	// Upload archive to the cloud
	err = uploadArchive(archfile)
	if err != nil {
		fmt.Printf("Failed to upload archive: %s\n", err)
		os.Exit(3)
	}

	// Show results
	fmt.Printf("Operation completed successfully\n")
	os.Exit(0)
}
