package main

import (
	"encoding/gob"
	"io"
	"net"
    "flag"
    "os"
    "log"
)



var (
    num_of_cpus := flag.Uint("cpus", 1, "number of cpu to use")
    port := flag.Uint("port", 100100, "port ...")
    out := flag.String("out", "output.txt", "output file name")
    num_of_workers := flag.Uint("workers", 0, "number of workers, MUST be set")
)


func main() {
    log.SetOutput(os.Stderr)
    flag.Parse()

    if num_of_workers == 0 {
        flag.PrintDefaults()
        log.Fatalln("number of workers must be set!")
    }
    runtime.GOMAXPROCS(num_of_cpus)
 
    ln,err := net.Listen("tcp", ":"+string(port))
    if err != nil {
        log.Fatalln(err)
    }

    // chans := make([]char string,num_of_workers)
    // for i,ch := range chans {
    //     chans[i] = make(chan string,10)
    // }

    

    n := 0
    for {
        if conn,err := ln.Accept();err != nil {
            log.Println(err)
            continue
        } else {
            if n == num_of_workers {
                log.Println("to many workers!")
                conn.Close()
            } else {
                n += 1
                go HandleConnection(conn)
            }

            
        }
    }

}