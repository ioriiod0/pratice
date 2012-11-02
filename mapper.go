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
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
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

var (
	config_file_path = flag.String("config", "config.json", "config file path")
	num_of_cpus      = flag.Uint("cpus", runtime.NumCPU(), "num of cpus")
)

type Config struct {
	inputs, workers []string
}

func load_config(file_name string) (*Config, error) {
	file, err := os.Open(file_name)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(file)
	config = new(Config)
	if err := dec.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

type Package int8

type Worker struct {
	addr string
	enc  *gob.Encoder
}

type Job struct {
	signal  chan bool
	input   string
	workers string
}

func main() {

	flag.Parse()
	runtime.GOMAXPROCS(num_of_cpus)
	log.SetOutput(os.Stderr)

	config, err := load_config(config_file_path)
	if err != nil {
		log.Fataln(err)
	}

	inputs, err := open_inputs(config.inputs)
	if err != nil {
		log.Fatalln(err)
	}

	outputs, err := open_outputs(config.outputs)
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
