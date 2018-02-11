# gopushbkp

## Overview
This program is a Go based backup application. It creates encrypted archives of
data files and it uploads them to S3 in the AWS Cloud. It is meant to be very
easy to use and to provide secure and reliable offsite backups. Backup archives
are uploaded to an S3 bucket as this storage service offers very reliable
protection against disk failures and corruptions. The restoration has to be done
manually using applications such as GnuPG to decrypt the archive and 7zip or
unzip to extract it.

## Security
The encryption is based on GnuPG/OpenPGP public and private keys. The backup
only requires the public key to be present. Hence the data would still be safe
if the public key could be read by an unauthorized person. The private key is
only required to decrypt backups and this is supposed to be rare. It is
recommended to keep it on a YubiKey. This offers very good protection for the
private key as it cannot be copied even if the computer where the decryption
happen had been compromised.

## Requirements
This program is built as a static binary and hence it does not require any
library to be present at run time. Also all the logic to compress, encrypt and
upload data are embedded in the program and unlike many scripts it does not rely
on any external program to perform these tasks. That is why it is very easy to
get this backup solution working.

Here are the requirements:
   * A public/private encryption keypair is required but only the public key
     needs to be accessible during the backup. This keypair can be generated
     using GnuPG. The public key can be exported to a text file using the
     following GnuPG command: "gpg --armor --export <keyid> > public-key.txt"
   * Enough disk space to store the archive on the computer where the backup
     is created
   * An S3 bucket where to upload backup archives and an access keypair to
     authenticate with AWS
   * A configuration file to provide the parameters required to perform the
     backup
   * A YubiKey is recommended for keeping the decryption key safe but it is
     optional.

## Implementation
This program first creates a PGP encrypted ZIP archive of the data. Both the
archiving and encryption are done on the fly hence no additional processing is
required at the end of the archiving. This saves both temporary disk space
and it is much faster to execute this way. Once the archive is ready it is
uploaded to the S3 bucket.

## How to use it
   * Create a bucket in S3 and create AWS Access Keys with access to this bucket
   * Copy the binary of this application to your system
   * Creates a configuration file named gopushbkp.toml and follow examples
     provided on github to provide values for all parameters
   * Run the binary. It will try to locate the configuration in multiple
     locations in the following order: first it will try to find it in the
     current directory, and next it will try to find it in the directory where
     the binary is stored.
