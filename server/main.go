package main

import (
	"crypto/tls"
	"flag"
	"log"
)

func main() {
	// Command-line flags
	addr := flag.String("addr", ":443", "listener's network address")
	flag.Parse()

	// Load public/private key pair
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	conf := tls.Config{Certificates: []tls.Certificate{cert}}

	// Listen on TCP
	ln, err := tls.Listen("tcp", *addr, &conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %s %s", ln.Addr().Network(), ln.Addr())

	defer ln.Close()
	for {
		// Wait for connection
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("New connection from %s", conn.RemoteAddr())

		go newClient(conn)
	}
}
