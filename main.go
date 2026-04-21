package main

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/creack/pty"
	"github.com/tarm/serial"
)

func main() {
	var args Args
	args.Read()

	physicalPort, err := serial.OpenPort(&serial.Config{
		Name: args.Device,
		Baud: args.BaudRate,
	})
	if err != nil {
		log.Fatalf("failed to open physical port %s: %v\n", args.Device, err)
	}
	defer physicalPort.Close()
	log.Printf("successfully opened physical port: %s\n", args.Device)

	var masters []*os.File
	var mu sync.Mutex

	for i := 0; i < args.Number; i++ {
		master, slave, err := pty.Open()
		if err != nil {
			log.Fatalf("failed to create virtual port %d: %v", i, err)
		}
		masters = append(masters, master)
		log.Printf("virtual port %d created: %s\n", i, slave.Name())
	}

	go func() {
		buf := make([]byte, 32)
		for {
			n, err := physicalPort.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("read physical port error: %v\n", err)
				}
				break
			}
			data := buf[:n]
			for _, m := range masters {
				_, _ = m.Write(data)
			}
		}
	}()

	for i, m := range masters {
		go func(index int, master *os.File) {
			buf := make([]byte, 32)
			for {
				n, err := master.Read(buf)
				if err != nil {
					log.Printf("virtual port %d disconnected\n", index)
					break
				}
				mu.Lock()
				_, err = physicalPort.Write(buf[:n])
				mu.Unlock()
				if err != nil {
					log.Printf("write to physical port failed: %v\n", err)
				}
			}
		}(i, m)
	}

	log.Println("multiplexing service running... Press Ctrl+C to exit")
	select {}
}
