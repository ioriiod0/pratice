

package main 

import (
    "bufio"
    "fmt"
    "log"
    "io"
    "net"
    "encoding/gob"
    "flag"
    "os"
    "strings"
)


func HandleConsole() {
    
}


type Console struct {
    ch chan string
}

func (self *Console) run() {
    defer close(self.ch)
    rd := bufio.NewReader(os.Stdin)
    for {
        if line,err := rd.ReadString('\n'); err != nil {
            if err == io.EOF {
                break
            } else {
                log.Fatalln(err)
            }
        } else {
            parts := strings.SplitN(line, " ",2)
            if len(parts) != 2 {
                log.Println("err input!")
                continue
            } else {
                f := strings.Trim(s, cutset)
            }


        }

    }
}

var (
    port := flag.Uint("port", 111111, "port to bind")
)

func main() {
    flag.Parse()
    log.SetOutput(os.Stderr)

    ln,err := net.Listen("tcp", ":"+string(port))
    if err != nil {
        log.Fatalln(err)
    }

    go HandleConsole()

    for {
        conn,err := ln.Accept()
        if err != nil {

        }
    }
}