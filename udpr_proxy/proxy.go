package main

import (
    "bytes"
    "fmt"
    "io"
    "net"
    "time"
)

const (
    ProxyPort  = 2013 // Port where the proxy listens
    ServerPort = 2012 // Original server port
    ReadTimeout = 5 * time.Second // Read timeout for server response
)

// Proxy function to handle client requests
func runProxyLoop(proxyConn *net.UDPConn) {
    buffer := make([]byte, 4096) // Initial buffer size

    for {
        n, clientAddr, err := proxyConn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Error reading from proxy connection: %v\n", err)
            continue
        }

        // Assume the first 4 bytes are the msgID, the rest is the data
        msgID := string(buffer[:4])
        data := buffer[4:n]

        // Forward request to main server port 2012
        go handleRequest(proxyConn, clientAddr, msgID, data)
    }
}

// Handle request by communicating with the server and forwarding to the client
func handleRequest(proxyConn *net.UDPConn, clientAddr *net.UDPAddr, msgID string, data []byte) {
    serverAddr := net.UDPAddr{
        Port: ServerPort,
        IP:   net.ParseIP("127.0.0.1"), // Assuming the server is on localhost
    }

    // Send the request to the server on port 2012
    serverConn, err := net.DialUDP("udp", nil, &serverAddr)
    if err != nil {
        fmt.Printf("Failed to dial server connection: %v\n", err)
        return
    }
    defer serverConn.Close()

    _, err = serverConn.Write(append([]byte(msgID), data...))
    if err != nil {
        fmt.Printf("Failed to send request to server: %v\n", err)
        return
    }

    // Use a bytes.Buffer to dynamically handle the response
    var responseBuffer bytes.Buffer

    serverConn.SetReadDeadline(time.Now().Add(ReadTimeout))
    _, err = io.Copy(&responseBuffer, serverConn)
    if err != nil && err != io.EOF {
        fmt.Printf("Failed to receive response from server: %v\n", err)
        return
    }

    // Forward the server's response to the client from port 2013
    _, err = proxyConn.WriteToUDP(responseBuffer.Bytes(), clientAddr)
    if err != nil {
        fmt.Printf("Failed to forward response to client: %v\n", err)
        return
    }

    // Handle the client's ACK directly
    handleClientAck(proxyConn, clientAddr, msgID, serverConn)
}

// Handle the client's ACK and forward it to the server's ephemeral port
func handleClientAck(proxyConn *net.UDPConn, clientAddr *net.UDPAddr, msgID string, serverConn *net.UDPConn) {
    ackBuffer := make([]byte, 4) // Expecting only the ACK ID

    // Wait for the ACK
    _, _, err := proxyConn.ReadFromUDP(ackBuffer)
    if err != nil {
        fmt.Printf("Failed to read ACK from client: %v\n", err)
        return
    }

    // Forward the ACK to the server's ephemeral port
    _, err = serverConn.Write(ackBuffer)
    if err != nil {
        fmt.Printf("Failed to forward ACK to server: %v\n", err)
    }
}

func main() {
    // Set up the proxy to listen on port 2013
    proxyAddr := net.UDPAddr{
        Port: ProxyPort,
        IP:   net.ParseIP("0.0.0.0"),
    }
    proxyConn, err := net.ListenUDP("udp", &proxyAddr)
    if err != nil {
        fmt.Printf("Failed to listen on proxy port %d: %v\n", ProxyPort, err)
        return
    }
    defer proxyConn.Close()

    // Start the proxy loop
    runProxyLoop(proxyConn)
}
