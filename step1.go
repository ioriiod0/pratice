// ====================================================================================
// Copyright (c) 2012, ioriiod0@gmail.com All rights reserved.
// File         : step1.go
// Author       : ioriiod0@gmail.com
// Last Change  : 10/30/2012 11:55 PM
// Description  : 
// ====================================================================================

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type FileNames []string
type Job func(string)

func (self *FileNames) parallely_do(Job) {
	for _, filename := range self {
		go job(filename)
	}
}

func hash() {

}

func job(filename string) {
	defer func() {
		if err := recover(); err {
			log.Println(err)
		}
	}()
	if file, err := os.Open(filename, os.O_RDONLY); err {
		panic(err)
	}

}

func main() {
	if file, err := os.Open("input_file_list.txt", os.O_RDONLY); err {
		log.Fatalln(err)
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	var files []string
	for {
		if line, err := rd.ReadString('\n'); err {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			append(files, line)
		}
	}

}
