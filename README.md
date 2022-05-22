<!-- implement 
build_status, (https://travis-ci.org) 
godoc, (https://godoc.org), (https://pkg.go.dev)
go_report_card here (goreportcard.com)
 !-->


# goEthereumWalletGen

goEthereumWalletGen is an ethereum wallet generator written entirely using Go. This application is Terminal-based and does not require an Internet connection, as it runs completely offline.

## Requirements

Requires Go 1.18 or newer

## Usage

There are two options to use the script (both require Go 1.18 or newer preinstalled)

Running the script manually using go
`go run *.go -pwd passwordToEncrypt -dir ./directoryToStoreKeystore -mnemonic "insert mnemonic phrase" -logging true`

Make an executable using the makefile
`make build`

And executing the program afterwards as follows 
`./ethWalletGen -pwd passwordToEncrypt -dir ./directoryToStoreKeystore -mnemonic "insert mnemonic phrase" -logging true`

### Flags

There are several optional flags one can use

* `-pwd passwordToEncrypt`
  * specifies a password to encrypt the keystore file with
* `-dir ./directoryToStoreKeystore`
  * Specifies the direcotry to store the keystore file in
  * default is `./wallets`
* `-mnemonic "my mnemonic key phrase"` 
  * One may use a mnemonic key phrase to generate a wallet off of
* `-logging` `true` OR `false`
  * logs privateKey and publicKey into the terminal
  * default is false

## How the application works

The application generates a random private key using the elliptic curve digital signature algorithm if no mnemonic key phrase is given.
In case a mnemonic key phrase is given using the flag `-mnemonic`, it uses that to get the corresponding private key from the key phrase.

After generating a private key, the application gets the public key from that private key and uses that public key to get the corresponding
ehtereum account address. After that it generates a keystore file in the directoy specified by the `-dir` flag and prints the account
address and the directory where it has stored the keystore file into the terminal.

### Example

If run correctly, the program output in the terminal should look as following

```
PrivateKey: a389e66e50e200d0bfc97a72b3b6c8802b171bb8e4a898c7023931f99237d1f4
PublicKey: 04d2cacfe6e0a19f4f25925356d296dd92a17bda27c6c42d58944af399930896b83c7a4a196d471e25d7a54670a3bad550396bd19ef43be88bfcd09f78383afdb2
Ethereum Wallet (0x4B7ABF923484e0A66c045019B1dCf92333FB4FAA) has been generated and stored in ./keystores
```

## Important

Never share your keystore file or privateKey with anyone!

### Known Issues

The application will throw an error if one tries to generate an ethereum wallet with the same mnemonic key `-mnemonic` and directory path `-dir` phrase in a row.
This is due to the program trying to create a keystore file which is already existant from running the program with the same properties beforehand.
