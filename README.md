## Screenshot of the results attached
![image](https://github.com/Grim-R3ap3r/dyte/assets/62543734/4130bfa1-4080-4774-b1e6-11d1d22c166c)

### Description of the solution
- ✅ Buffer Pooling: Instead of allocating a new buffer for each write operation, buffers are obtained from the pool and returned after use. This reduces the overhead of memory allocation and garbage collection.
- ✅ Parallelization: Write operations are parallelized across multiple goroutines, allowing for concurrent writing to each port. This leverages the concurrency features of Go to maximize throughput.
- ✅ Minimized Synchronization: Synchronization is minimized by using channels and wait groups to coordinate the completion of write operations. Goroutines communicate through channels to signal when they have finished writing, and a wait group is used to ensure all goroutines have completed before proceeding.
- ✅ Reuse of UDP connections: UDP connections are reused within each goroutine to avoid the overhead of establishing new connections for each write operation. This reduces the latency associated with connection setup and teardown.

### Trade-offs
- ✅ Resource Utilization: By parallelizing write operations and reusing UDP connections, the solution may consume more system resources, such as CPU and memory. 
- ✅ Potential for Resource Contentions: With increased parallelism, there is a higher potential for resource contentions, such as contention for access to the buffer pool or contention for network resources. This could potentially lead to performance degradation under high load or in multi-threaded environments.



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
