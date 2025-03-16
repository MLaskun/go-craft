package connection

import (
	"net"
)

const (
	SEGMENT_BITS = 0x7F
	CONTINUE_BIT = 0x80
)

func readVarInt(conn net.Conn) (int32, error) {
	var value int32 = 0
	var position uint = 0
	var currentByte byte

	for {
		var b [1]byte
		_, err := conn.Read(b[:])
		if err != nil {
			return 0, err
		}
		currentByte = b[0]

		value |= int32(currentByte&SEGMENT_BITS) << position

		if (currentByte & CONTINUE_BIT) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			panic("VarInt is too big")
		}
	}

	return value, nil
}

func readVarText(conn net.Conn) (string, error) {
	size, err := readVarInt(conn)
	if err != nil {
		return "", err
	}

	buf := make([]byte, size)
    _, err = conn.Read(buf[:])
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func writePacket(conn net.Conn, packetID byte, data []byte) {
	length := len(data) + 1
	conn.Write(append([]byte{byte(length), packetID}, data...))
}
