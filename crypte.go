package main

// incluir los errores en las devoluciones de las funciones

// secretbox -- symmetric
// box -- asymmetric
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
	keys    = flag.Bool("k", false, "Pub/Priv Key pair generator")
	cryp    = flag.Bool("e", false, "Crypt message")
	decryp  = flag.Bool("d", false, "Decrypt message")
	pK      = flag.String("p", "", "Public Key File")
	sK      = flag.String("s", "", "Private Key File")
	m       = flag.String("m", "", "Ecrypted message")
	pubKey  nacl.Key
	privKey nacl.Key
	enc     []byte
)

func main() {
	// Flags config
	flag.Usage = usage
	flag.Parse()
    usage()
	// Public/Private keys generation
	if *keys {
		pubKey, privKey, err := genKeys()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Public Key: %s\n", base64.StdEncoding.EncodeToString([]byte(pubKey[:])))
		fmt.Printf("Private Key: %s\n", base64.StdEncoding.EncodeToString([]byte(privKey[:])))

		if !(writeKeyToFile(pubKey, "publicKey")) {
			fmt.Printf("Can not write Public Key")
		}

		if !(writeKeyToFile(privKey, "privateKey")) {
			fmt.Printf("Can not write Private Key")
		}
	}
	msg := "The quick brown fox jumps over the lazy dog"
	if *cryp {
		enc = crypt(*pK, *sK, msg)
	}

	if *decryp {
		decrypt(*pK, *sK, *m)
	}
}

func usage() {
	fmt.Printf("\nERROR - Must complete all input params\n")
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
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return b

}

func writeKeyToFile(k nacl.Key, n string) (ok bool) {
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
	// Convert key to bytes slice
	b := []byte(k[:])
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

func genKeys() (pubKey, privKey nacl.Key, err error) {
	pubKey, privKey, err = box.GenerateKey(crypto_rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pub/Priv key pair created!")
	return
}

func crypt(p, k string, m string) []byte {
	/*
		var nonce [24]byte
		if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
			panic(err)
		}
	*/
	nonce := nacl.NewNonce()
	fmt.Printf("Nonce: %s\n", base64.StdEncoding.EncodeToString([]byte(nonce[:])))
	fmt.Printf("Decrypted message: %s\n", m)
	pKey := readFile(p)
	pk := new([nacl.KeySize]byte)
	copy(pk[:], pKey)
	sKey := readFile(k)
	sk := new([nacl.KeySize]byte)
	copy(sk[:], sKey)
	//fmt.Printf("%T\n", pKey)
	/*
		    a, err := nacl.Load(string(pKey))
		    if err != nil {
			    fmt.Println("Can not load public key")
		    }
		    fmt.Println(a)
		    sKey := readFile(k)
		    b, _ := nacl.Load(string(sKey))
		    fmt.Println(base64.StdEncoding.EncodeToString(b[:]))
	*/
	var out []byte
	enc := box.Seal(out, []byte(m), nonce, pk, sk)
	fmt.Printf("Encrypted message: %s\n", base64.StdEncoding.EncodeToString(enc))
	return enc
}

func decrypt(p, k, m string) {
	// Read public key from file
	pKey := readFile(p)
	pk := new([nacl.KeySize]byte)
	copy(pk[:], pKey)
	// Read private key from file
	sKey := readFile(k)
	sk := new([nacl.KeySize]byte)
	copy(sk[:], sKey)
	/*
		var decNonce [24]byte
		copy(decNonce[:], e[:24])
	*/
	msg, err := base64.StdEncoding.DecodeString(m)
	dec, err := box.EasyOpen([]byte(msg), pk, sk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted message: %s\n", base64.StdEncoding.EncodeToString(dec))
}
