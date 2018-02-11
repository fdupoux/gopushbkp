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
	"github.com/spf13/viper"
	"os"
	"path"
)

type configuration struct {
	datarootdir string
	bkparchdir  string
	bkpbasename string
	pubkeyfile  string
	excludes    []string
	aws_region  string
	aws_bucket  string
	aws_access  string
	aws_secret  string
	aws_prefix  string
}

func readConfig() (configuration, error) {

	myconfig := configuration{}

	dircwd, err := os.Getwd()
	if err != nil {
		return myconfig, fmt.Errorf("Failed to get current working archiveDirectory: %s\n", err)
	}

	pathexe, err := os.Executable()
	if err != nil {
		return myconfig, fmt.Errorf("Failed to find path to executable: %s\n", err)
	}

	direxe := path.Dir(pathexe)

	viper.SetConfigName("gopushbkp")
	viper.AddConfigPath("/etc/gopushbkp/")
	viper.AddConfigPath(dircwd)
	viper.AddConfigPath(direxe)

	err = viper.ReadInConfig()
	if err != nil {
		return myconfig, fmt.Errorf("Failed to read configuration file: %s\n", err)
	}

	if viper.IsSet("datarootdir") == true {
		myconfig.datarootdir = viper.GetString("datarootdir")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'datarootdir' in configuration file\n")
	}

	if viper.IsSet("bkparchdir") == true {
		myconfig.bkparchdir = viper.GetString("bkparchdir")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'bkparchdir' in configuration file\n")
	}

	if viper.IsSet("bkpbasename") == true {
		myconfig.bkpbasename = viper.GetString("bkpbasename")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'bkpbasename' in configuration file\n")
	}

	if viper.IsSet("pubkeyfile") == true {
		myconfig.pubkeyfile = viper.GetString("pubkeyfile")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'pubkeyfile' in configuration file\n")
	}

	if viper.IsSet("excludes") == true {
		myconfig.excludes = viper.GetStringSlice("excludes")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'excludes' in configuration file\n")
	}

	if viper.IsSet("aws_region") == true {
		myconfig.aws_region = viper.GetString("aws_region")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'aws_region' in configuration file\n")
	}

	if viper.IsSet("aws_bucket") == true {
		myconfig.aws_bucket = viper.GetString("aws_bucket")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'aws_bucket' in configuration file\n")
	}

	if viper.IsSet("aws_access") == true {
		myconfig.aws_access = viper.GetString("aws_access")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'aws_access' in configuration file\n")
	}

	if viper.IsSet("aws_secret") == true {
		myconfig.aws_secret = viper.GetString("aws_secret")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'aws_secret' in configuration file\n")
	}

	if viper.IsSet("aws_prefix") == true {
		myconfig.aws_prefix = viper.GetString("aws_prefix")
	} else {
		return myconfig, fmt.Errorf("Failed to find entry 'aws_prefix' in configuration file\n")
	}

	return myconfig, nil
}
