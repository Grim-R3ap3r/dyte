## Screenshot of the results attached
![image](https://github.com/Grim-R3ap3r/dyte/assets/62543734/54cde585-e5c8-4623-a236-bc5ea24f0324)


### Description of the solution
- ✅ Preallocation of UDP Connections:  Instead of creating a new UDP connection for each write operation within the benchmark loop, the solution preallocates UDP connections before running the benchmark. This avoids the overhead of repeatedly dialing UDP connections for each iteration, leading to improved performance.
- ✅ Concurrent Writing: The benchmark utilizes goroutines to write data concurrently to multiple UDP connections. Each goroutine is responsible for writing data to a specific UDP connection, allowing for parallel processing of write operations. This concurrent approach maximizes CPU utilization and reduces overall execution time compared to sequential writing.
- ✅ Shared Buffer: To minimize memory allocations, a shared buffer is used for writing data to UDP connections. Instead of allocating a new buffer for each write operation, the same buffer is reused across all goroutines. This reduces memory allocation overhead and improves memory efficiency.

### Trade-offs
- ✅ Increased Memory Usage: Preallocating UDP connections and using a shared buffer for writing data may lead to increased memory usage compared to the baseline implementation. This is because resources are allocated upfront and retained throughout the benchmark execution. While this approach improves performance by reducing overhead, it may consume more memory, especially for a large number of UDP connections or when handling large volumes of data.



# UDP Network Writer Challenge

## Introduction
This is a hiring challenge for the intern role for core media team at our company. The goal of this challenge is to write a UDP network writer that performs significantly better than the baseline implementation.

## Challenge Description
In the `benchs_test.go` file, you will find a benchmark test that measures the performance of the UDP network writer. Your task is to improve the performance of the writer so that it is at least 2 times faster than the baseline implementation.

## Getting Started
To run the benchmark test, use the following command:
```bash
go test -benchmem -bench BenchmarkConnections
```

## Submission
To submit your solution, please create a pull request to this repository with your implementation. Make sure to include a description of your solution and any trade-offs you made.

## Evaluation
The golang benchmark test logs are of the following structure:
```
BenchmarkConnections-8   	<number of iterations per second>	       <time for processing each iteration> ns/op    <bytes allocated per iteration> B/op   <number of allocations per iteration> allocs/op
```
Your solution will be evaluated based on the following criteria (in this order):
- Number of terations per second
- Time for processing each iteration
- Code quality, E.g. readability, maintainability, and performance
- Explanation of the solution and trade-offs
- Bytes allocated per iteration
- Number of allocations per iteration
