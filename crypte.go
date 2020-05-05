// Copyright (c) 2020 José Fagoaga jose.fagoagasancho@telefonica.com
//
// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the Free Software
// Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
// PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// this program. If not, see <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: GPL-3.0-onl

package main

/* To-Do
- Return errors on functions
- Write key to a custom file
*/

import (
	crypto_rand "crypto/rand"
	//	compress "custom/compress"
	"encoding/base64"
	"flag"
	"fmt"
	nacl "github.com/kevinburke/nacl"
	box "github.com/kevinburke/nacl/box"
	"log"
	"os"
	//"bufio"
)

var (
	keys       = flag.Bool("k", false, "Pub/Priv Key pair generator")
	enc        = flag.Bool("e", false, "Encrypt, sign and compress with lz4 a message")
	dec        = flag.Bool("d", false, "Decrypt, verify sign adn decompress with lz4 a message")
	publicKey  = flag.String("p", "", "Public Key File")
	privateKey = flag.String("s", "", "Private Key File")
	msg        = flag.String("m", "", "Message to encrypt/decrypt")
)

func main() {
	// Print banner
	banner()
	// Flags config
	flag.Usage = usage
	flag.Parse()
	if !*keys && !*enc && !*dec {
		fmt.Printf("\nERROR - Must complete all input params\n")
		usage()
	}
	// Public/Private keys generation
	if *keys {
		pubKey, privKey, err := genKeys()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Created Public/Private Key")
			fmt.Printf("Public Key: %s\n", base64.StdEncoding.EncodeToString([]byte(pubKey[:])))
			fmt.Printf("Private Key: %s\n", base64.StdEncoding.EncodeToString([]byte(privKey[:])))
			// Write keys to file
			//fmt.Println("Enter file in which to save the key (Default: (public/private)Key)")
			//reader := bufio.NewReader(os.Stdin)
			//line, _ := reader.ReadString('\n')
			p := keyToByte(pubKey)
			writeToFile(p, "publicKey")
			s := keyToByte(privKey)
			writeToFile(s, "privateKey")
		}
	}
	// Encrypt message
	if *enc {
		log.Printf("Encryption started\n")
		message := readFile(*msg)
		// compress message
		//message_c := compress.Compress(message)
		enc := encrypt(*publicKey, *privateKey, message)
		writeToFile(enc, *msg+".enc")
	}
	// Decrypt message
	if *dec {
		log.Printf("Decryption started\n")
		message := readFile(*msg)
		dec := decrypt(*publicKey, *privateKey, message)
		// decompress message
		//message_d := compress.Decompress(dec)
		writeToFile(dec, *msg+".dec")
	}
}

func banner() {
	fmt.Println("## Crypte ##")
	fmt.Printf("Tool for (de)compress and (de)crypt message\n\n")
}

func usage() {
	fmt.Printf("\nUsage:\n")
	fmt.Println("- Generate Public/Private Keys:")
	fmt.Printf("%s -k\n", os.Args[0])
	fmt.Println("- Encrypt, sign and compress with lz4 a message")
	fmt.Printf("%s -e -p <PublicKeyFile> -s <PrivateKeyFile> -m <Message>\n", os.Args[0])
	fmt.Println("- Decrypt, verify sign and decompress with lz4  a message:")
	fmt.Printf("%s -d -p <PublicKeyFile> -s <PrivateKeyFile> -m <Message>\n", os.Args[0])
	os.Exit(1)
}

func genKeys() (publicKey, privateKey nacl.Key, err error) {
	// Generate Public/Private Keys
	publicKey, privateKey, err = box.GenerateKey(crypto_rand.Reader)
	return
}

func encrypt(publicKey, privateKey string, message []byte) []byte {
	// Read Public Key
	pk := readKeyFile(publicKey)
	// Read Private Key
	sk := readKeyFile(privateKey)
	// Encrypt message
	enc := box.EasySeal([]byte(message), pk, sk)
	// Print crypted message (b64 encoding)
	//fmt.Printf("Encrypted message: %s\n", base64.StdEncoding.EncodeToString(enc))
	return enc
}

func decrypt(publicKey, privateKey string, message []byte) []byte {
	// Read Public Key
	pk := readKeyFile(publicKey)
	// Read Private Key
	sk := readKeyFile(privateKey)
	// Decode message
	fmt.Printf("Encrypted message: %s\n", base64.StdEncoding.EncodeToString(message))
	// Decrypt message
	dec, err := box.EasyOpen(message, pk, sk)
	if err != nil {
		log.Fatal(err)
	}
	// Print decrypted message
	//fmt.Printf("Decrypted message: %s\n", string(dec))
	return dec
}
