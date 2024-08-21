package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "sync/atomic"
)

// Ensure that the shutdown process is also communicated clearly
func shutdownHandler(conn *net.UDPConn, shutdownFlag *int32) {
    interruptCount := 0 // Scoped to this function
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

    for sig := range signals {
        interruptCount++
        log.Printf("Signal received: %v", sig)  // Logging the signal for better traceability

        if interruptCount == 1 {
            fmt.Println("Interrupt received, initiating graceful shutdown...")
            fmt.Println("Press Ctrl+C up to 9 times in total to force quit immediately.")
            atomic.StoreInt32(shutdownFlag, 1)  // Signal to shutdown the manager and other components
            conn.Close()    // Close the listener to stop accepting new connections
            continue           // Skip to the next iteration
        }

        if interruptCount == 9 {
            fmt.Println("Force quitting after 9 interrupts...")
            os.Exit(1) // Force exit after receiving 9 interrupts
        }

        fmt.Printf("Interrupt received (%d/9), press Ctrl+C again to force quit...\n", interruptCount)
    }
}
