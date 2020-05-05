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

import (
	nacl "github.com/kevinburke/nacl"
	"io/ioutil"
	"log"
	"os"
)

// Read input file
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

// Transform nacl.Key to a byte array
func keyToByte(key nacl.Key) []byte {
	// Convert key to bytes slice
	b := []byte(key[:])
	return b
}

// Read key file
func readKeyFile(file string) nacl.Key {
	// Read key
	f := readFile(file)
	// Transform keys from type []byte to nacl.Key
	key := new([nacl.KeySize]byte)
	copy(key[:], f)
	return key
}

// Write output file
func writeToFile(b []byte, n string) {
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
		log.Printf("Can not write" + n + "\n")
		log.Fatal(err)
	}
	log.Printf("File writed: '%s'\n", n)
}
