package main

import (
	"bufio"
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
	log.Printf("successfully opened physical port: %s (LineMode: %v)\n", args.Device, args.LineMode)

	var masters []*os.File
	var mu sync.Mutex

	for i := 0; i < args.Number; i++ {
		master, slave, err := pty.Open()
		if err != nil {
			log.Fatalf("failed to create virtual port %d: %v", i, err)
		}
		if err = os.Chmod(slave.Name(), 0660); err != nil {
			log.Printf("chmod failed: %v\n", err)
		}
		masters = append(masters, master)
		log.Printf("virtual port %d created: %s\n", i, slave.Name())
	}

	go func() {
		if args.LineMode {
			scanner := bufio.NewScanner(physicalPort)
			for scanner.Scan() {
				data := scanner.Bytes()
				line := append(data, '\n')

				log.Printf("RX (Line) %d bytes: %q\n", len(line), line)

				for _, m := range masters {
					_, _ = m.Write(line)
				}
			}
			if err := scanner.Err(); err != nil {
				log.Printf("physical port scanner error: %v\n", err)
			}
		} else {
			buf := make([]byte, 32)
			for {
				n, err := physicalPort.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("read physical port error: %v\n", err)
					}
					break
				}
				if n == 0 {
					continue
				}

				log.Printf("RX %d bytes: %q\n", n, buf[:n])
				for _, m := range masters {
					_, _ = m.Write(buf[:n])
				}
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
