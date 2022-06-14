package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type Frame struct {
	Cmd    string   `json:"cmd"`
	Sender string   `json:"sender"`
	Data   []string `json:"data"`
}

type Info struct {
	myNum    int
	nextNode string
	nextNum  int
	imFirst  bool
	cont     int
}

var (
	host         string
	remotes      []string
	chInfo       chan Info
	readyToStart chan bool
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		log.Println("Hostname not given")
	} else {
		go prepareAgrawalla()
		host = os.Args[1]
		if len(os.Args) == 3 {
			secondNode()
		} else {
			go agrawallaStart()
		}
		server()
	}
}

func server() {
	if ln, err := net.Listen("tcp", host); err == nil {
		defer ln.Close()
		log.Printf("Listening on %s\n", host)
		for {
			if cn, err := ln.Accept(); err == nil {
				go fauxDispatcher(cn)
			} else {
				log.Printf("%s: can't accept connection.\n", host)
			}
		}
	} else {
		log.Printf("Can't listen on %s\n", host)
	}
}

func secondNode() {
	remotes = append(remotes, os.Args[2])
	if !send(remotes[0], Frame{"hola", host, []string{}}, func(cn net.Conn) {
		dec := json.NewDecoder(cn)
		var frame Frame
		dec.Decode(&frame)
		remotes = append(remotes, frame.Data...)
		log.Printf("%s: friends: %s\n", host, remotes)
	}) {
		log.Printf("Unable to connect to %s\n", remotes)
	}
}

func send(remote string, frame Frame, callback func(net.Conn)) bool {
	if cn, err := net.Dial("tcp", remote); err == nil {
		defer cn.Close()
		enc := json.NewEncoder(cn)
		enc.Encode(frame)
		if callback != nil {
			callback(cn)
		}
		return true
	} else {
		log.Printf("%s: can't connect to %s\n", host, remote)
		idx := -1
		for i, rem := range remotes {
			if rem == remote {
				idx = i
				break
			}
		}
		if idx >= 0 {
			remotes = append(remotes[:idx], remotes[idx+1:]...)
		}
		return false
	}
}

func fauxDispatcher(cn net.Conn) {
	defer cn.Close()
	dec := json.NewDecoder(cn)
	var frame Frame
	dec.Decode(&frame)
	log.Println(frame)
	switch frame.Cmd {
	case "hola":
		enc := json.NewEncoder(cn)
		enc.Encode(Frame{"oh", host, remotes})
		notification := Frame{"agrega", host, []string{frame.Sender}}
		for _, remote := range remotes {
			send(remote, notification, nil)
		}
		remotes = append(remotes, frame.Sender)
		log.Printf("%s: friends: %s\n", host, remotes)
	case "agrega":
		remotes = append(remotes, frame.Data...)
		log.Printf("%s: friends: %s\n", host, remotes)
	case "agrawalla":
		for _, remote := range remotes {
			info := <-chInfo
			msg := Frame{"myNumIs", host, []string{strconv.Itoa(info.myNum)}}
			send(remote, msg, nil)
			go func() {
				chInfo <- info
			}()
		}
	case "myNumIs":
		handleSomeNum(&frame)
	case "start":
		<-readyToStart
		criticalSection()

		/*case "potato":
		handlePotato(&frame)*/
	}
}

func handleSomeNum(frame *Frame) {
	if num, err := strconv.Atoi(frame.Data[0]); err == nil {
		info := <-chInfo
		if num > info.myNum {
			if num < info.nextNum {
				info.nextNode = frame.Sender
				info.nextNum = num
			}
		} else {
			info.imFirst = false
		}
		info.cont++
		go func() {
			chInfo <- info
		}()
		if info.cont == len(remotes) {
			if info.imFirst {
				log.Printf("%s I'm first", host)
				criticalSection()
			} else {
				readyToStart <- true
			}
		}
	} else {
		log.Printf("Can't convert %s sent from %s\n", frame.Data[0], frame.Sender)
	}
}

func criticalSection() {
	/*for len(remotes) > 0 {
		remote := remotes[rand.Intn(len(remotes))]
		data := []string{strconv.Itoa(rand.Intn(20) + 10)}
		if send(remote, Frame{"potato", host, data}, nil) {
			break
		}
	}*/
	log.Printf("%s my time has come to start, so I'll do my thing.\n", host)
	info := <-chInfo
	if info.nextNode == "" {
		log.Printf("%s I was the last of my kind... :(\n", host)
	} else {
		log.Printf("%s I shall inform %s that his time is nigh", host, info.nextNode)
		send(info.nextNode, Frame{"start", host, []string{}}, nil)
	}
}

func prepareAgrawalla() {
	chInfo = make(chan Info)
	readyToStart = make(chan bool)
	randNum := rand.Intn(int(1e9))
	log.Printf("%s, my number is %d\n", host, randNum)
	chInfo <- Info{
		myNum:    randNum,
		nextNode: "",
		nextNum:  int(1e9),
		imFirst:  true,
		cont:     0,
	}
}

func agrawallaStart() {
	time.Sleep(5 * time.Second)
	for _, remote := range remotes {
		send(remote, Frame{"agrawalla", host, []string{}}, nil)
	}
}

/*
func handlePotato(frame *Frame) {
	if num, err := strconv.Atoi(frame.Data[0]); err == nil {
		log.Printf("%s: recibí %d from %s\n", host, num, frame.Sender)
		if num == 0 {
			log.Printf("%s: perdí\n", host)
		} else {
			for len(remotes) > 0 {
				remote := remotes[rand.Intn(len(remotes))]
				data := []string{strconv.Itoa(num - 1)}
				time.Sleep(100 * time.Millisecond)
				if send(remote, Frame{"potato", host, data}, nil) {
					break
				}
			}
		}
	} else {
		log.Printf("%s: can't convert %s to number\n", host, frame.Data)
	}
}*/
