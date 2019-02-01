package main

import (
	"fmt"
	"encoding/binary"
	"io"
	"bufio"
	"os"
	"bytes"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	inputfile, err := os.Open("/var/tmp/ut")
	check(err)
	hashfile, err := os.Create("hashfile")
	check(err)

	defer inputfile.Close()
	defer hashfile.Close()

	bufferedReader := bufio.NewReader(inputfile)
	hashWriter := bufio.NewWriter(hashfile)
	defer hashWriter.Flush()

	var count uint32
	var currentword []byte
	var currentHash []byte
	hashArray := make([]uint32, 30*30*30)

	for {
		line, err := bufferedReader.ReadBytes(0x0a) // 0x0a /n i LATIN-1
		if err == io.EOF {
			break
		}
		
		tmp := bytes.Split(line, []byte(" ")) // Split on space
		newWord := tmp[0]

		if !bytes.Equal(currentword, newWord) {
			currentword = newWord
			newHash := FirstThree(currentword)
			
			if !bytes.Equal(currentHash, newHash) {
				currentHash = newHash
				index := Hash(currentHash)
				hashArray[index] = count
			}
		}

		count += uint32(len(line))
	}

	for i := range(hashArray) {
		b := make([]byte, 4)
		n := hashArray[i]
		binary.LittleEndian.PutUint32(b, n)
		hashWriter.Write(b)
	}
}
