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
// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"fmt"
	"github.com/pierrec/lz4"
	"log"
)

func compress(data []byte) []byte {
	buf := make([]byte, len(data))
	ht := make([]int, 64<<10) // buffer for the compression table, 64KB
	n, err := lz4.CompressBlock(data, buf, ht)
	if err != nil {
		log.Println("Can not compress data")
		log.Fatal(err)
	}
	if n >= len(data) {
		fmt.Printf("`%s` is not compressible", string(data))
	}
	buf = buf[:n] // compressed data
	fmt.Printf("%x\n", buf)
	return buf
}

func decompress(data []byte) []byte {
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
