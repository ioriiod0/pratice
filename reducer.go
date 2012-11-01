package main

import (
	"encoding/gob"
	"io"
	"net"
    "flag"
    "os"
    "log"
    "math"
)

type Query struct {
    Freq int32
    Str string
}

type ACK struct {
    Seq int32
}

type WorkerAgent struct {
    conn net.Conn
    ch chan Query
}

func NewWorkerAgent(conn net.Conn) *WorkerAgent {
    return &WorkerAgent{conn,make(chan Query,100)}
}

func (self *WorkerAgent) run() {
    defer func () {
        if err:=recover();err != nil {
            log.Println(err)
        }
        close(self.ch)
        self.conn.Close()
    }
    wd := gob.NewEncoder(self.conn)
    rd := gob.NewDecoder(self.conn)
    var query Query
    var ack ACK
    for {
        if err := rd.Decode(query);err != nil {
            panic(err)
        } else {
            self.ch <- query
        }
        if err := wd.Encode(ack);err != nil {
            panic(err)
        }
    }
}

func (self *Query) get() Query {
    return <-self.ch
}
///////////////////////////////////////

type LeafNode Query
type NonLeafNode int32

type LoserTree struct {
    nodes []interface{}
    index int32
}

func NewLoserTree(k uint) *LoserTree {
    tree := &LoserTree{make([]interface{},k*2),k}
    for i := 0;i < k;i += 1 {
        tree.nodes[i] = -1
    }
}


func (self *LoserTree) GetValue(id uint32) {
    var freq uint32
    switch v := self.nodes[id].(type)
    {
    case NonLeafNode:
        if v == -1 {
            freq = math.MaxUint32
        } else {
            freq = self.nodes[v].(LeafNode).Freq
        }
    case LeafNode:
        freq = v.Freq
    }
    return freq
}


func (self *LoserTree) Compete(defender,attacker uint32) {
    if defender == 0 {
        self.nodes[defender] = attacker
        return
    }

    new_defender := defender/2
    new_attacker := attacker

    if self.GetValue(attacker) < self.GetValue(defender) 
    {
        new_attacker,self.nodes[defender] = defender,attacker
    }
    
    self.Compete(new_defender,new_attacker)
}

func (self *LoserTree) Init(querys []Query) {
    length := len(querys) 
    if len(self.nodes) / 2 != length {
        log.Panicln("logic err in Init")
    }

    for i,query := range querys {
        self.index = length + i
        self.Push(query)
    }
}

func (self *LoserTree) Pop() LeafNode {
    index := self.nodes[0].(NonLeafNode)
    if index == -1 {
        log.Panicln("logic err in Pop")
    }
    self.index = index
    return self.nodes[index].(LeafNode)
}

func (self *LoserTree) Push(node LeafNode) {
    if self.index == -1 {
        log.Panicln("logic err in Push")
    }

    self.nodes[self.index] = node 
    self.Compete(self.index/2,self.index)
    self.index = -1
}

///////////////////////////////////////////
type Reducer struct {
    workers []*WorkerAgent
    tree *LoserTree
}


///////////////////////////////////////////

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

    ch := make(chan string,num_of_workers)


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
                go HandleConnection(conn,ch)
                n += 1
            }

            
        }
    }

}