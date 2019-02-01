package main

import (
	"strconv"
	"bytes"
	"fmt"
	"encoding/binary"
	"bufio"
	"os"
	"io"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	index, err := os.Open("var/tmp/ut")
	check(err)
	korpus, err := os.Open("info/adk18/labb1/korpus")
	check(err)

	defer index.Close()
	defer korpus.Close()

	hashArray := slurpHash()

	word := os.Args[1]
	if !inputCheck([]byte(word)) {
		fmt.Println("Ogiltigt argument")
		return
	}

	pointer := search([]byte(word), hashArray, index)

	if pointer == -1 {
		fmt.Println("Det ordet finns inte.")
		return
	}

	index.Seek(pointer, 0)
	count := read([]byte(word), index, korpus, true)
	index.Seek(pointer, 0)

	if count > 25 {
		fmt.Printf("Det finns %d matchingar. Vill du ha dem utskrivna?  y/n\n", count)
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if answer == "y\n" {
			read([]byte(word), index, korpus, false)
		}
	} else {
		fmt.Printf("Antal: %d\n", count)
		read([]byte(word), index, korpus, false)
	}
}

func inputCheck(word []byte) bool {
	for i := range word {
		r := word[i]
		if !((r >= 97 && r < 122) || r == 0xe4 || r == 0xe5 || r  == 0xf6) {
			return false
		}
	}

	return true
}

func read(word []byte, index *os.File, korpus *os.File, onlyCounting bool) int {
	reader := bufio.NewReader(index)
	count := 0

	for {
		currentword, err := reader.ReadBytes(' ')
		if err == io.EOF {
			break
		}

		currentword = currentword[:len(currentword) - 1]

		if !bytes.Equal(word, currentword) {
			return count
		}

		pointer, _ := reader.ReadBytes('\n')

		if !onlyCounting {
			pointer = pointer[:len(pointer) - 1]
			pointer, _ := strconv.Atoi(string(pointer))
			korpus.Seek(int64(pointer) - 30, 0)
			b := make([]byte, 60)
			_, err := korpus.Read(b)
			check(err)
			fmt.Println(string(oneLine(b)))
		}

		count++
	}

	return count
}

func oneLine(word []byte) []byte {
	line := make([]byte, 60)

	for i := range word {
		if word[i] != 0x0a {
			line[i] = word[i]
		} else {
			line[i] = 0x20
		}
	}
	
	return line
}

func search(word []byte, hashArray []uint32, index *os.File) int64 {
	reader := bufio.NewReader(index)
	wprefix := Hash(FirstThree(word))
	i := int64(hashArray[wprefix])

	if i == 0 {	
		return -1
	}

	var j int64

	if wprefix != uint32(len(hashArray) - 1) {
		j = int64(hashArray[wprefix + 1])
	}

	for j-i > 1000 {
		m := (i + j) / 2 //automagic flooring for integers
		m = readback(m, index) //read back to beginning of line
		index.Seek(m, 0)
		b, e := reader.ReadBytes(byte(' '))
		check(e)

		if bytes.Compare(b[:len(b) - 1], word) > 0 {
			i = m
		} else {
			j = m
		}

		reader.Discard(reader.Buffered())
	}
	
	index.Seek(i, 0)
	offset := i
	for {
		s, e := reader.ReadBytes(byte('\n')) //next word
		if e == io.EOF {
			return -1
		}

		split := bytes.Split(s, []byte(" "))
		if bytes.Equal(split[0], word) {
			return offset
		} else if bytes.Compare(split[0], word) == 1 {
			return -1 //oh noes
		}

		offset += int64(len(s))
	}

}

func readback(offset int64, file *os.File) int64 {
	b := make([]byte, 1)
	file.ReadAt(b, offset)

	for ; b[0] != byte('\n'); offset-- {
		_, e := file.ReadAt(b, offset)
		check(e)

	}

	return offset
}

func slurpHash() []uint32 {
	hashfile, err := os.Open("hashfile")
	check(err)
	reader := bufio.NewReader(hashfile)
	hashArray := make([]uint32, 30*30*30)

	for i := 0; ; i++ {
		b := make([]byte, 4)
		_, err := reader.Read(b)
		if err == io.EOF {
			break
		}

		hashArray[i] = binary.LittleEndian.Uint32(b)
	}

	return hashArray
}
