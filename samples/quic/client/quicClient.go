package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/quic-go/quic-go"
)

const (
	serverIP          = "3.83.100.0"
	serverPort        = "5282"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "jarkom-quic-sample-evelyn"
	sslKeyLogFileName = "ssl-key.log"
)

func main() {

	sslKeyLogFile, err := os.Create(sslKeyLogFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer sslKeyLogFile.Close()

	fmt.Printf("QUIC Client Socket Program Example in Go\n")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{appLayerProto},
		KeyLogWriter:       sslKeyLogFile,
	}
	connection, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	defer connection.CloseWithError(0x0, "No Error")

	fmt.Printf("[quic] Dialling from %s to %s\n", connection.LocalAddr(), connection.RemoteAddr())

	fmt.Printf("[quic] Input message to be sent to server: ")
	message, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	stream1, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] Opened bidirectional stream %d to %s\n", stream1.StreamID(), connection.RemoteAddr())

	stream2, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] Opened bidirectional stream %d to %s\n", stream2.StreamID(), connection.RemoteAddr())

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Printf("[quic] Creating receive buffer of size %d\n", bufferSize)
		receiveBuffer := make([]byte, bufferSize)

		fmt.Printf("[quic] [Stream ID: %d] Sending message '%s'\n", stream1.StreamID(), message)
		_, err := stream1.Write([]byte(message))
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("[quic] [Stream ID: %d] Message sent\n", stream1.StreamID())

		receiveLength, err := stream1.Read(receiveBuffer)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("[quic] [Stream ID: %d] Received %d bytes of message from server\n", stream1.StreamID(), receiveLength)

		response := receiveBuffer[:receiveLength]
		fmt.Printf("[quic] [Stream ID: %d] Received message: '%s'\n", stream1.StreamID(), response)
	}()

	go func() {
		defer wg.Done()

		fmt.Printf("[quic] Creating receive buffer of size %d\n", bufferSize)
		receiveBuffer := make([]byte, bufferSize)

		fmt.Printf("[quic] [Stream ID: %d] Sending message '%s'\n", stream2.StreamID(), message)
		_, err := stream2.Write([]byte(message))
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("[quic] [Stream ID: %d] Message sent\n", stream2.StreamID())

		receiveLength, err := stream2.Read(receiveBuffer)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("[quic] [Stream ID: %d] Received %d bytes of message from server\n", stream2.StreamID(), receiveLength)

		response := receiveBuffer[:receiveLength]
		fmt.Printf("[quic] [Stream ID: %d] Received message: '%s'\n", stream2.StreamID(), response)
	}()

	wg.Wait()
}