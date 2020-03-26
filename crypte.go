package main

// incluir los errores en las devoluciones de las funciones

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	nacl "github.com/kevinburke/nacl"
	box "github.com/kevinburke/nacl/box"
	"io/ioutil"
	"log"
	"os"
)

var (
	keys       = flag.Bool("k", false, "Pub/Priv Key pair generator")
	enc        = flag.Bool("e", false, "Encrypt message")
	dec        = flag.Bool("d", false, "Decrypt message")
	publicKey  = flag.String("p", "", "Public Key File")
	privateKey = flag.String("s", "", "Private Key File")
	msg        = flag.String("m", "", "Message to encrypt/decrypt")
)

func main() {
	// Flags config
	flag.Usage = usage
	flag.Parse()
	//	usage()
	// fmt.Printf("\nERROR - Must complete all input params\n")
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
			p := keyToByte(pubKey)
			if !(writeToFile(p, "publicKey")) {
				fmt.Printf("Can not write Public Key\n")
			}
			s := keyToByte(privKey)
			if !(writeToFile(s, "privateKey")) {
				fmt.Printf("Can not write Private Key\n")
			}
		}
	}
	// Encrypt message
	//msg := "The quick brown fox jumps over the lazy dog"
	if *enc {
		message := readFile(*msg)
		enc := encrypt(*publicKey, *privateKey, message)
		if !(writeToFile(enc, *msg+".enc")) {
			fmt.Printf("Can not write encrypted message\n")
		}
	}
	// Decrypt message
	if *dec {
		message := readFile(*msg)
		dec := decrypt(*publicKey, *privateKey, message)
		if !(writeToFile(dec, "message.dec")) {
			fmt.Printf("Can not write decrypted message\n")
		}
	}
}

func usage() {
	fmt.Printf("\nUsage:\n")
	fmt.Println("- Generate Public/Private Keys:")
	fmt.Printf("%s -k\n", os.Args[0])
	fmt.Println("- Encrypt message:")
	fmt.Printf("%s -e -p <PublicKeyFile> -s <PrivateKeyFile> -m <Message>\n", os.Args[0])
	fmt.Println("- Decrypt message:")
	fmt.Printf("%s -d -p <PublicKeyFile> -s <PrivateKeyFile> -m <Message>\n", os.Args[0])
	os.Exit(1)
}

func readFile(f string) []byte {
	// Open file
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Read file
	b, err := ioutil.ReadAll(file)
	return b
}

func keyToByte(key nacl.Key) []byte {
	// Convert key to bytes slice
	b := []byte(key[:])
	return b
}

func writeToFile(b []byte, n string) (ok bool) {
	// Open a new file for writing only
	file, err := os.OpenFile(
		n,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0600,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Write bytes to file
	_, err = file.Write(b)
	if err != nil {
		ok = false
		log.Fatal(err)
	}
	ok = true
	log.Printf("File writed: '%s'\n", n)
	return
}

func genKeys() (publicKey, privateKey nacl.Key, err error) {
	// Generate Public/Private Keys
	publicKey, privateKey, err = box.GenerateKey(crypto_rand.Reader)
	return
}

func encrypt(publicKey, privateKey string, message []byte) []byte {
	/*
		var nonce [24]byte
		if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
			panic(err)
		}
	*/

	// Random nonce (var nonce [24]byte)
	nonce := nacl.NewNonce()
	//fmt.Printf("Nonce: %s\n", base64.StdEncoding.EncodeToString([]byte(nonce[:])))

	// Read Public Key
	pk := readKeyFile(publicKey)
	// Read Private Key
	sk := readKeyFile(privateKey)

	// Encrypt message
	var out []byte
	enc := box.Seal(out, []byte(message), nonce, pk, sk)
	// Print crypted message (b64 encoding)
	fmt.Printf("Decrypted message: %s\n", message)
	fmt.Printf("Encrypted message: %s\n", base64.StdEncoding.EncodeToString(enc))
	return enc
}

func decrypt(publicKey, privateKey string, message []byte) []byte {
	/* If box.Open
	var decNonce [24]byte
	copy(decNonce[:], e[:24])
	*/

	// Read Public Key
	pk := readKeyFile(publicKey)
	// Read Private Key
	sk := readKeyFile(privateKey)
	// Decode message
	//fmt.Printf("b64 message: %s\n", message)
	//msg, err := base64.StdEncoding.DecodeString(message)
	//fmt.Printf("non encoded message: %q\n", msg)
	// Decrypt message
	dec, err := box.EasyOpen(message, pk, sk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted message: %s\n", base64.StdEncoding.EncodeToString(dec))
	return dec
}

func readKeyFile(file string) nacl.Key {
	// Read key
	f := readFile(file)
	// Transform keys from type []byte to nacl.Key
	key := new([nacl.KeySize]byte)
	copy(key[:], f)
	return key
}
