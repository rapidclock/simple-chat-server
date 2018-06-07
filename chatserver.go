package main

import (
	"log"
	"net"
	"io"
	"fmt"
)

const listenAddr = "localhost:4000"
const findingMatchMsg = "Finding a match for you..."
const partnerFoundMsg = "Partner Found!"
const reQMsg = "Sorry Requeuing..."

func main() {
	l, err := net.Listen("tcp", listenAddr)
	logErr(err)
	defer l.Close()
	closeConChannel := make(chan io.ReadWriteCloser)
	conQ := make(chan io.ReadWriteCloser)
	go terminator(closeConChannel)
	go matcher(conQ, closeConChannel)
	for {
		c, err := l.Accept()
		logErr(err)
		conQ<- c
	}
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func matcher(conQ chan io.ReadWriteCloser, closeCh chan<- io.ReadWriteCloser) {
	fmt.Println("Matcher Started...")
	for {
		a := <-conQ
		fmt.Fprintln(a, findingMatchMsg)
		b := <- conQ
		if isConAlive(a) {
			if isConAlive(b) {
				fmt.Fprintln(io.MultiWriter(a, b), partnerFoundMsg)
				go chat(a, b, conQ, closeCh)
			} else {
				closeCh<- b
				checkAndRequeue(a, conQ, closeCh)
			}
		} else {
			closeCh<- a
			checkAndRequeue(b, conQ, closeCh)
		}
	}
}

func chat(a, b io.ReadWriteCloser, conQ chan<- io.ReadWriteCloser, closeCh chan<- io.ReadWriteCloser) {
	errChA := make(chan error)
	errChB := make(chan error)
	go link(a, b, errChB)
	go link(b, a, errChA)
	fmt.Println("Chat has commenced between two parties...")
	select {
	case <-errChA:
		closeCh<- a
		checkAndRequeue(b, conQ, closeCh)
	case <-errChB:
		closeCh<- b
		checkAndRequeue(a, conQ, closeCh)
	}
	fmt.Println("Chat Ended :(")
}

func link(w io.Writer, r io.Reader, errChR chan<- error) {
	_, err := io.Copy(w, r)
	errChR<-err
}

func terminator(terminate <-chan io.ReadWriteCloser) {
	fmt.Println("Terminator has Risen!")
	for {
		c := <- terminate
		c.Close()
	}
}

func checkAndRequeue(c io.ReadWriteCloser, conQ chan<- io.ReadWriteCloser, closeCh chan<- io.ReadWriteCloser) {
	if isConAlive(c) {
		fmt.Fprintln(c, reQMsg)
		conQ<- c
	} else {
		closeCh<- c
	}
}

func isConAlive(c io.Writer) bool {
	one := []byte("")
	_, err := c.Write(one)
	return err == nil
}
