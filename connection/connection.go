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

	if state == 2 {
		fmt.Println("Logging in")
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
	_, err = conn.Read(discard)
	if err != nil {
		fmt.Println("Error checking legacy ping")
		return
	}
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
	_, err := readVarInt(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sendServerStatus(conn)

	payload := make([]byte, 8)
	_, err = conn.Read(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sendPingResponse(conn, payload)
}

func sendPingResponse(conn net.Conn, payload []byte) {
	writePacket(conn, 0x01, payload)
}

func sendServerStatus(conn net.Conn) {
	response := `{"version":{"name":"1.21.4","protocol":769},"players":{"max":20,"online":0},"description":{"text":"Japiernicze Dziala"}}`
	data := wrapString(response)
	packet := append([]byte{byte(len(data) + 1), 0x00}, data...)
	_, err := conn.Write(packet)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
