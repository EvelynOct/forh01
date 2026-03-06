package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
	quicutils "jarkom.cs.ui.ac.id/h01/samples/quic/utils"
)

const (
	serverIP          = ""
	serverPort        = "5282"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "lrt-jabodebek-2406365282"
	sslKeyLogFileName = "ssl-key.log"
)

func Handler(packet utils.LRTPIDSPacket) string {
	if packet.IsTrainArriving {
		return "Attention, the train to " + packet.Destination + " will arrive at Platform 1"
	}

	if packet.IsTrainDeparting {
		return "Attention, the train to " + packet.Destination + " will depart from Platform 1"
	}

	return ""
}

func main() {
	localUDPAddress, err := net.ResolveUDPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalln(err)
	}

	socket, err := net.ListenUDP(serverType, localUDPAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer socket.Close()

	tlsConfig := &tls.Config{
		Certificates: quicutils.GenerateTLSSelfSignedCertificates(),
		NextProtos:   []string{appLayerProto},
	}

	listener, err := quic.Listen(socket, tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection quic.Connection) {
	for {
		stream, err := connection.AcceptStream(context.Background())
		if err != nil {
			return
		}

		go handleStream(stream)
	}
}

func handleStream(stream quic.Stream) {
	defer stream.Close()

	rawMessage, err := io.ReadAll(stream)
	if err != nil {
		log.Println(err)
		return
	}

	packet := utils.Decoder(rawMessage)

	message := Handler(packet)
	fmt.Println(message)

	packet.IsAck = true
	ackPacket := utils.Encoder(packet)

	_, err = stream.Write(ackPacket)
	if err != nil {
		log.Println(err)
	}
}