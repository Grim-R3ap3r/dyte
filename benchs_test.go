package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

var portIdx atomic.Int64
var readersCount = 100

const (
	defaultReadBufferSize  = 8 * 1024 * 1024
	defaultWriteBufferSize = 8 * 1024 * 1024
)

func newUDPSocket() (fd int, port int, err error) {
	// Create local udp socket on any random port
	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		panic(err)
	}

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		panic(err)
	}

	addr := [4]byte{127, 0, 0, 1}
	portBase := 5000
	for {
		port = int(portIdx.Add(1) + int64(portBase))
		err = syscall.Bind(fd, &syscall.SockaddrInet4{
			Port: port,
			Addr: addr,
		})
		if err == nil {
			break
		}
		port++
	}
	err = syscall.SetNonblock(fd, false)
	return
}

// DO NOT MODIFY THIS FUNCTION
func testInit(readersCount int, verbose bool) (ports []int, readChan chan []byte, closeChan chan struct{}, err error) {
	ports = make([]int, readersCount)
	portsChan := make(chan int, readersCount)

	readChan = make(chan []byte, readersCount)
	closeChan = make(chan struct{}, 1)

	// Create readersCount udp sockets to read
	wg := sync.WaitGroup{}
	for i := 0; i < readersCount; i++ {
		wg.Add(1)
		go func(threadId int) {
			fd, port, err := newUDPSocket()
			if err != nil {
				return
			}

			portsChan <- port

			buf := make([]byte, 1500)
			wg.Done()
			for {
				select {
				case <-closeChan:
					return
				default:
					n, _, err := syscall.Recvfrom(fd, buf, 0)
					if err != nil {
						// Close the socket
						syscall.Close(fd)
						return
					}
					if verbose {
						data := string(buf[:n])
						fmt.Println("threadId", threadId, "read: ", n, "bytes", "data: ", data)
					}
					readChan <- buf[:n]
				}
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < readersCount; i++ {
		ports[i] = <-portsChan
	}

	return
}

// DO NOT MODIFY THIS FUNCTION
func getTestMsg() []byte {
	// Generate a 1500 byte random message
	buf := make([]byte, 1500)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}

// DO NOT MODIFY THIS FUNCTION
func waitForReaders(readChan chan []byte, b *testing.B) {
	// Wait for reader to read
	for i := 0; i < readersCount; i++ {
		select {
		case <-readChan:
		case <-time.After(1 * time.Second):
			b.Fatal("timeout") // This should not happen
		}
	}
}

func BenchmarkConnections(b *testing.B) {
	b.Run("baseline", func(b *testing.B) {
		BenchmarkRawUDP(b)
	})

	b.Run("Sample", func(b *testing.B) {
		BenchmarkSample(b)
	})
}

func BenchmarkRawUDP(b *testing.B) {
	b.StopTimer()

	testPort := 40101
	// Create a udp network connection
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: testPort,
	})
	if err != nil {
		b.Fatal(err)
	}

	ports, readChan, closeChan, err := testInit(readersCount, false)
	if err != nil {
		b.Fatal(err)
	}
	_ = readChan

	writer := func() {
		for i := 0; i < readersCount; i++ {
			buf := getTestMsg()
			_, err := conn.WriteTo(buf, &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: ports[i],
			})
			if err != nil {
				b.Fatal(err)
			}
		}

		// End of code that you are permitted to modify
		waitForReaders(readChan, b)
	}

	// Sequential test
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		writer()
	}
	b.StopTimer()

	close(closeChan)
}

// Implement your benchmarks here -->
// Please read the comments carefully. You need to implement something atleast much faster than the baseline
func BenchmarkSample(b *testing.B) {
    b.StopTimer()

    // Initialize ports, readChan, and closeChan
    ports, readChan, closeChan, err := testInit(readersCount, false)
    if err != nil {
        b.Fatal(err)
    }
    defer close(closeChan)

    // Create a pool of reusable buffers
    bufPool := sync.Pool{
        New: func() interface{} {
            return make([]byte, 1500)
        },
    }

    // Channel to signal completion of write operations
    done := make(chan struct{}, readersCount)

    // Parallelized writing to each port
    var wg sync.WaitGroup
    wg.Add(readersCount)
    for i := 0; i < readersCount; i++ {
        go func(port int) {
            defer func() {
                done <- struct{}{}
                wg.Done()
            }()

            // Get buffer from the pool
            buf := bufPool.Get().([]byte)
            defer bufPool.Put(buf)

            // Create a UDP connection for writing
            conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
                IP:   net.IPv4(127, 0, 0, 1),
                Port: ports[port],
            })
            if err != nil {
                b.Fatalf("failed to dial UDP: %v", err)
            }
            defer conn.Close()

            // Perform write operations
            for i := 0; i < b.N; i++ {
                // Generate test message
                rand.Read(buf)

                // Write message to the port
                _, err := conn.Write(buf)
                if err != nil {
                    b.Fatalf("failed to write to port %d: %v", ports[port], err)
                }
            }
        }(i)
    }

    // Start the timer
    b.StartTimer()

    // Wait for all write operations to complete
    go func() {
        wg.Wait()
        close(done)
    }()

    // Wait for completion of all write operations
    for range done {
    }

    // Stop the timer
    b.StopTimer()

    // Wait for readers to read all messages
    for i := 0; i < readersCount; i++ {
        <-readChan
    }
}