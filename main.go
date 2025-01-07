package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    // Listen on UDP port 50000
    addr := net.UDPAddr{
        Port: 50000,
        IP: net.ParseIP("0.0.0.0"),
    }
    
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Printf("Failed to bind to port 50000: %v\n", err)
        os.Exit(1)
    }
    defer conn.Close()
    
    fmt.Printf("Server listening on %s\n", conn.LocalAddr().String())
    
    buffer := make([]byte, 1024)
    for {
        n, remoteAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Error receiving data: %v\n", err)
            continue
        }
        
        if string(buffer[:n]) == "whoami" {
            response := fmt.Sprintf("addr:%s", remoteAddr.String())
            _, err := conn.WriteToUDP([]byte(response), remoteAddr)
            if err != nil {
                fmt.Printf("Error sending response: %v\n", err)
                continue
            }
            fmt.Printf("Responded to %s with their address\n", remoteAddr.String())
        }
    }
}
