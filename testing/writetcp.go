package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func writeTCP() {
	const name = "writetcp"
	log.SetPrefix(name + "\t")

	port := flag.Int("p", 8080, "server port")
	flag.Parse()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: *port})

	if err != nil {
		log.Fatal("Error connecting to localhost:%d: %v", *port, err)
	}
	log.Printf("Connected to %s: will forward stdin", conn.RemoteAddr())

	defer conn.Close()

	go func() {

		for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {

			fmt.Printf("%s\n", connScanner.Text())

			if err := connScanner.Err(); err != nil {
				log.Fatalf("Error reading from %s: %v", conn.RemoteAddr())
			}
			if connScanner.Err() != nil {
				log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
			}
		}
	}()

	for stdinScanner := bufio.NewScanner(os.Stdin); stdinScanner.Scan(); {
		log.Printf("Sent: %s\n", stdinScanner.Text())
		if _, err := conn.Write(stdinScanner.Bytes()); err != nil {
			log.Fatalf("Error writing to %s: %v", conn.RemoteAddr(), err)
		}
		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Fatalf("Error writing to %s: %v", conn.RemoteAddr(), err)
		}
		if stdinScanner.Err() != nil {
			log.Fatalf("Error reading from %s: %v", conn.RemoteAddr(), err)
		}

	}

}
