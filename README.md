# cryptgo

Simplistic CLI tool for encrypting and decrypting files based on RSA keys.

[![cryptgo-build-test](https://github.com/bartossh/cryptgo/actions/workflows/go.yml/badge.svg)](https://github.com/bartossh/cryptgo/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/bartossh/cryptgo/branch/main/graph/badge.svg?token=1748BU8XY2)](https://codecov.io/gh/bartossh/cryptgo)

## Features

- Asymmetric encryption/decryption.
- Hash SHA512 based.
- Random source of entropy ensure that encrypting the same message twice doesn't result in the same ciphertex. 
- Encrypts single file using your system `~/.ssh/id_rsa` to generate new public key.
- Decrypts single file using your system `~/.ssh/id_rsa` (private) key.
- Can encrypt and decrypt files with RSA key generated with passphrase.
- Encrypt single file using auto generated rsa key, saves key in to given path
- Decrypts single file using file key from provided path

## Usage

- dependencies

Please be sure you have rsa key: `RSA PRIVATE KEY` created under `~/.ssh/` folder

To create one write in terminal: `ssh-keygen -t rsa -m PEM`

- help

`cryptgo --help`

- encrypt

`cryptgo --input <file path to enrypt> --output <new encrypted file path> encrypt`

- decrypt

`cryptgo --input <path to encrypted file> --output <new decrypted file path> decrypt`

- passwd encrypt

`cryptgo --input <file to enrypt path> --output <new encrypted file path> --passwd <passphrase> encrypt`

- passwd decrypt

`cryptgo --input <encrypted file path> --output <new decrypted file path> --passwd <passphrase> decrypt`

- auto generate key encrypt

`cryptgo --input <file path to enrypt> --output <new encrypted file path> -generate <path where to save rsa key> encrypt`

- use key from path decrypt

`cryptgo --input <path to encrypted file> --output <new decrypted file path> -use <path to rsa key to be used> decrypt`

## Building

This software runs on OSX and Unix based systems, it is not working on Windows yet.

Build with `go build .`

To compress binary build with `go build -ldflags="-s -w"` then run `upx -9 -v cryptgo` (~2MB binary)

## Development

- Recommended go version is `go1.17` or higher.
- Write an issue please or make a PR against main branch (there is no development branch yet)
- Please write test for core functionalities before making PR.

## Features to implement in the future

- Support for folders encryption/decryption.
- Support for other asymmetric keys.
- P2P encrypted message communication. 
- make it available for windows
