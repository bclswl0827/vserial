package main

import (
	"io"
	"log"
	"os"

	"github.com/creack/pty"
	"github.com/tarm/serial"
)

func createDevice(port *serial.Port, master *os.File, slave *os.File) error {
	_, err := io.Copy(master, port)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var args Args
	args.Read()

	serialPort, err := serial.OpenPort(&serial.Config{
		Name: args.Device,
		Baud: args.BaudRate,
	})
	if err != nil {
		log.Fatalln(err)
	}

	defer serialPort.Close()
	log.Printf("source port %s opened\n", args.Device)

	for i := 0; i < args.Number; i++ {
		master, slave, err := pty.Open()
		if err != nil {
			log.Fatalln(err)
		}
		defer master.Close()
		defer slave.Close()

		log.Printf("virtual port %s created\n", slave.Name())
		if i == args.Number-1 {
			err = createDevice(serialPort, master, slave)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			go createDevice(serialPort, master, slave)
		}
	}
}
