package connection

import (
	"encoding/binary"
	"fmt"
	"net"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

    handleLegacyPing(conn)

	state, err := handleHandshake(conn)
	if err != nil {
		fmt.Println("Handshake error:", err)
		return
	}

	if state == 1 {
		handleStatus(conn)
	}
}

func handleLegacyPing(conn net.Conn) {
	legacyPingID := make([]byte, 1)
	discard := make([]byte, 1)

    _, err := conn.Read(legacyPingID)
    if err != nil {
        fmt.Println("Error checking legacy ping")
        return
    }
	fmt.Println("ping:", legacyPingID)

    _, err = conn.Read(discard)
    if err != nil {
        fmt.Println("Error checking legacy ping")
        return
    }
	fmt.Println("ping:", discard)
}

func handleHandshake(conn net.Conn) (int, error) {
	protocolVersion, err := readVarInt(conn)
	if err != nil {
		return 0, err
	}
	fmt.Println("Client Protocol Version:", protocolVersion)

	serverAddress, err := readVarText(conn)
	if err != nil {
		return 0, err
	}
	fmt.Println("Server address:", serverAddress)

    port := make([]byte, 2)
    _, err = conn.Read(port)
    if err != nil {
        return 0, err
    }
    fmt.Println("Server port:", binary.BigEndian.Uint16(port))

	nextState, err := readVarInt(conn)
	if err != nil {
		return 0, err
	}
	fmt.Println("Next state:", nextState)

	return int(nextState), nil
}

func handleStatus(conn net.Conn) {
	packetID, err := readVarInt(conn)
	if err != nil || packetID != 0x00 {
        //TODO
        fmt.Println("---------------------Debug Log!!!-----------------------")
		return
	}

	sendServerStatus(conn)

	packetID, err = readVarInt(conn)
	if err != nil || packetID != 0x01 {
		return
	}

	payload := make([]byte, 8)
	_, err = conn.Read(payload)
	if err != nil {
		return
	}

	sendPingRespons(conn, payload)
}

func sendPingRespons(conn net.Conn, payload []byte) {
	writePacket(conn, 0x01, payload)
}

func sendServerStatus(conn net.Conn) {
	response := `{"version":{"name":"1.21.4","protocol":769},"players":{"max":20,"online":0},"description":{"text":"Japiernicze Dziala"}}`
	data := append([]byte{byte(len(response))}, []byte(response)...)
	conn.Write(data)
}
