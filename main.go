package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Proceso struct {
	id          int
	i           int
	logEnabled  bool
	log         chan bool
	hasFinished bool
}

func (self *Proceso) start() {
	go func() {
		for {
			var ok bool
			self.logEnabled, ok = <-self.log
			if !ok {
				self.hasFinished = true
			}
		}
	}()

	for !self.hasFinished {
		self.i++
		if self.logEnabled {
			fmt.Printf("id %d: %d\n", self.id, self.i)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (self *Proceso) close() {
	close(self.log)
}

func printMenu() {
	fmt.Println("1) Agregar proceso")
	fmt.Println("2) Mostrar procesos")
	fmt.Println("3) Terminar proceso")
	fmt.Println("4) Salir")
}

type Menu struct {
	scanner   *bufio.Scanner
	opc       string
	processes []Proceso
	idCount   int
}

func (self *Menu) processOpc() {
	switch self.opc {
	case "1":
		self.addProcess()
	case "2":
		self.showProcesses()
	case "3":
		self.endProcess()
	}
}

func (self *Menu) addProcess() {
	p := Proceso{
		id:  self.idCount,
		log: make(chan bool),
	}
	self.idCount++
	go p.start()

	self.processes = append(self.processes, p)
}

func (self *Menu) showProcesses() {
	for _, process := range self.processes {
		process.log <- true
	}
	self.scanner.Scan()
	for _, process := range self.processes {
		process.log <- false
	}
}

func (self *Menu) endProcess() {
	var id int
	fmt.Scan(&id)

	for i := 0; i < len(self.processes); i++ {
		if self.processes[i].id == id {
			self.processes[i].close()
			self.processes = append(self.processes[:i], self.processes[i+1:]...)
		}
	}
}

func (self *Menu) start() {
	self.scanner = bufio.NewScanner(os.Stdin)
	self.processes = make([]Proceso, 0)

	self.opc = ""
	for self.opc != "4" {
		printMenu()
		self.scanner.Scan()
		self.opc = self.scanner.Text()
		self.processOpc()
	}
}

func main() {
	menu := Menu{}
	menu.start()
}
