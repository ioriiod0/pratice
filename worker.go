package main

import (
	"bufio"
	"encoding/gob"
	"github.com/pmylund/sortutil"
	"io"
	"log"
	"net"
    "flag"
    "runtime"
    "os"
)


///////////////////////////////////////////////////////////
type TrieNode struct {
	index int
	nodes []*TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{index: -1, nodes: make([]*TrieNode, 26)}
}

func (self *TrieNode) Find(key byte) *TrieNode {
	var offset int
	if key <= 'z' && key >= 'a' {
		offset := key - 'a'
	} else {
		log.Panicf("unkown key:%v\n", key)
	}

	if self.nodes[offset] == nil {
		self.nodes[offset] = NewTrieNode(key)
	}
	return self.node[offset]
}

///////////////////////////////////////////////

type Query struct {
    Freq int
    Str  string
}

type Querys []Query

func (self *Querys) Sort() {
	sortutil.DescByField(self, "Freq")
}

func (self *Querys) Each(f func(Query)) {
    for _:query := range self {
        f(query)
    }
}


////////////////////////////////////////////////
type Trie struct {
	root   *TrieNode
	querys []Query
}

func NewTrie(size) {
	return &Trie{NewTrieNode(), make([]Query, size)}
}

func (self *Trie) Insert(str string) {
    p := self.root
    for _,b := range byte(str) {
        p := p.Find(b)
    }
    if p.index != -1 {
        querys[p.index].Freq += 1
    } else {
        p.index = len(querys)
        querys = append(querys,Query{1,str})
    }
}

////////////////////////////////////////

type Response struct {
    Req int
}

func HandleInput(trie *Trie,ch chan string,out_addr string) {

    defer func(){
        if err := recover();err != nil {
            log.Println(err)
        }
    }()

    for {
        if str,ok := <-ch;ok {
            trie.Insert(str)
        } else {
            break
        }
    }

    trie.querys.Sort()

    conn,err := net.Dial("tcp", out_addr)
    if err != nil {
        log.Fatalln(err)
    }

    wd := gob.NewEncoder(conn)
    rd := gob.NewDecoder(conn)
    var res Response

    trie.querys.Each(func(q Query) {
        if err := wd.Encode(q); err!=nil {
            log.Fatalln(err)
        }
        if err := rd.Decode(res); err!=nil {
            log.Fatalln(err)
        }
    })

    conn.Close()
}



type Request struct {
    Str string
}

func HandleConnection(conn *net.TCPConn,ch chan string) {
    defer func(){
        if err := recover();err != nil {
            if err != io.EOF {
                log.Println(err)
            }
        }
        close(ch)
    }()
    
    rd := gob.NewDecoder(conn)
    var req Request
    for {
        if err := rd.Decode(req);err != nil {
            panic(err)
        }
        ch <- req.Str
    }
}


////////////////////////////////////////////////////

var (
    num_of_cpus := flag.Uint("cpus", 1, "number of cpu to use")
    port := flag.Uint("port", "9090", "port...")
    out := flag.String("out", "0.0.0.0:100100", "output/reducer's ip adress")
)

func main() {
    flag.Parse()
    runtime.GOMAXPROCS(num_of_cpus)
    log.SetOutput(os.Stderr)

    ln,err := net.Listen("tcp", ":"+string(port))
    if err != nil {
        log.Fatalln(err)
    }

    for {
        if conn,err := ln.Accept(); err != nil {
            log.Println(err)
            continue
        } else {

            trie := NewTrie(1000);
            ch := make(chan String,10)
            signal := make(chan bool)

            go HandleConnection(conn, ch)
            go HandleInput(trie,ch,out)

        }

    }



}
