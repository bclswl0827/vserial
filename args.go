package main

import "flag"

func (a *Args) Read() {
	flag.StringVar(&a.Device, "device", "/dev/ttyS0", "source serial device name or path")
	flag.IntVar(&a.BaudRate, "baudrate", 9600, "source serial device baud rate")
	flag.IntVar(&a.Number, "number", 1, "number of virtual serial devices to create")
	flag.Parse()
}
