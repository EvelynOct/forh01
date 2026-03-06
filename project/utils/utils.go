package utils

import (
	"bytes"
	"encoding/binary"
)

// Tipe data komponen boleh diubah, namun variabelnya jangan diubah
type LRTPIDSPacketFixed struct {
	TransactionId     uint16
	IsAck             bool
	IsNewTrain        bool
	IsUpdateTrain     bool
	IsDeleteTrain     bool
	IsTrainArriving   bool
	IsTrainDeparting  bool
	TrainNumber       uint16
	DestinationLength uint8
}

type LRTPIDSPacket struct {
	LRTPIDSPacketFixed
	Destination string
}

func Encoder(packet LRTPIDSPacket) []byte {
	buffer := new(bytes.Buffer)

	var flags uint8
	flags = 0

	if packet.IsAck {
		flags = flags + 32
	}
	if packet.IsNewTrain {
		flags = flags + 16
	}
	if packet.IsUpdateTrain {
		flags = flags + 8
	}
	if packet.IsDeleteTrain {
		flags = flags + 4
	}
	if packet.IsTrainArriving {
		flags = flags + 2
	}
	if packet.IsTrainDeparting {
		flags = flags + 1
	}

	binary.Write(buffer, binary.BigEndian, packet.TransactionId)
	binary.Write(buffer, binary.BigEndian, flags)
	binary.Write(buffer, binary.BigEndian, packet.TrainNumber)
	binary.Write(buffer, binary.BigEndian, packet.DestinationLength)
	buffer.Write([]byte(packet.Destination))

	return buffer.Bytes()
}

func Decoder(rawMessage []byte) LRTPIDSPacket {
	reader := bytes.NewReader(rawMessage)

	var transactionId uint16
	var flags uint8
	var trainNumber uint16
	var destinationLength uint8

	binary.Read(reader, binary.BigEndian, &transactionId)
	binary.Read(reader, binary.BigEndian, &flags)
	binary.Read(reader, binary.BigEndian, &trainNumber)
	binary.Read(reader, binary.BigEndian, &destinationLength)

	destinationBytes := make([]byte, destinationLength)
	reader.Read(destinationBytes)

	packet := LRTPIDSPacket{}

	packet.TransactionId = transactionId
	packet.TrainNumber = trainNumber
	packet.DestinationLength = destinationLength
	packet.Destination = string(destinationBytes)

	if flags >= 32 {
		packet.IsAck = true
		flags = flags - 32
	}
	if flags >= 16 {
		packet.IsNewTrain = true
		flags = flags - 16
	}
	if flags >= 8 {
		packet.IsUpdateTrain = true
		flags = flags - 8
	}
	if flags >= 4 {
		packet.IsDeleteTrain = true
		flags = flags - 4
	}
	if flags >= 2 {
		packet.IsTrainArriving = true
		flags = flags - 2
	}
	if flags >= 1 {
		packet.IsTrainDeparting = true
	}

	return packet
}