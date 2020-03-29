package compress

import (
	"fmt"
	"github.com/pierrec/lz4"
	"log"
)

func Compress(data []byte) []byte {
	buf := make([]byte, len(data))
	ht := make([]int, 64<<10) // buffer for the compression table, 64KB
	n, err := lz4.CompressBlock(data, buf, ht)
	if err != nil {
		log.Println("Can not compress data\n")
		log.Fatal(err)
	}
	if n >= len(data) {
		fmt.Printf("`%s` is not compressible", string(data))
	}
	buf = buf[:n] // compressed data
    fmt.Printf("%x\n", buf)
    return buf
}

func Decompress(data []byte) []byte {
    // Allocated a very large buffer for decompression.
    out := make([]byte, 10*len(data))
    n, err := lz4.UncompressBlock(data, out)
    if err != nil {
        fmt.Println(err)
    }
    out = out[:n] // uncompressed data
    fmt.Println(string(out[:len(data)]))
    return out
}
