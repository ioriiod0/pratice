// ====================================================================================
// Copyright (c) 2012, ioriiod0@gmail.com All rights reserved.
// File         : mapper.go
// Author       : ioriiod0@gmail.com
// Last Change  : 10/30/2012 11:55 PM
// Description  : 
// ====================================================================================

package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func open_input_files(in io.Reader) ([]io.Reader, error) {

	var ret []io.Reader
	rd := bufio.NewReader(in)

	for {
		if file_name, err := rd.ReadString('\n'); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else {
			if file, err := os.Open(file_name); err != nil {
				return nil, err
			} else {
				ret = append(ret, file)
			}
		}
	}
	return ret, nil
}

func open_output_files(size uint32) ([]io.Writer, error) {

	ret := make([]io.Writer, size)
	for i := uint32(0); i < size; i += 1 {
		file, err := os.OpenFile(string(i)+".txt", os.O_TRUNC|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		} else {
			ret[i] = file
		}
	}
	return ret, nil
}

func mapper(id int, in io.Reader, outs []io.Writer, signals chan int) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		signals <- id
	}()

	rd := bufio.NewReader(in)
	size := len(outs)
	wds := make([]*bufio.Writer, size, size)
	for i, out := range outs {
		wds[i] = bufio.NewWriter(out)
	}

	for {
		if line, err := rd.ReadString('\n'); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Panicln(err)
			}
		} else {
			index := crc32.ChecksumIEEE([]byte(line)) % uint32(size) //hash
			if _, err := wds[index].WriteString(line); err != nil {
				log.Panicln(err)
			}
		}
	}

}

func usage() {
	fmt.Println("usage: mapper [N=?] where ? is a number >= 1")
	fmt.Println("examle: mapper N=10")
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutput(os.Stderr)

	var num_of_outputs uint32

	if len(os.Args) != 2 {
		usage()
		os.Exit(-1)
	}

	parts := strings.Split(os.Args[1], "=")
	if len(parts) != 2 {
		usage()
		os.Exit(-1)
	}

	if tmp, err := strconv.ParseUint(parts[1], 10, 32); err != nil {
		usage()
		os.Exit(-1)
	} else {
		num_of_outputs = uint32(tmp)
	}

	inputs, err := open_input_files(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	outputs, err := open_output_files(num_of_outputs)
	if err != nil {
		log.Fatalln(err)
	}

	num_of_workers := len(inputs)
	signals := make(chan int, num_of_workers)
	defer func() {
		close(signals)
	}()

	for id, input := range inputs {
		go mapper(id, input, outputs, signals)
	}

	for i := 0; i < num_of_workers; i += 1 {
		<-signals
	}

	for i := uint32(0); i < num_of_outputs; i += 1 {
		fmt.Println(string(i) + ".txt")
	}

}
