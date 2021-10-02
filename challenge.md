# Technical Challenge (Golang)

## Overview

Using Golang, write a server ("application") that opens a socket and restricts input to at most 5
concurrent clients. Clients will connect to the application and write any number of 9 digit numbers,
and then close the connection. The application must write a de-duplicated list of these numbers to
a log file in no particular order.

## Primary Considerations

1. The application should work correctly as defined below in Requirements.
2. The overall structure of the application should be simple.
3. The code of the Application should be descriptive and easy to read, and the build method and
   runtime parameters must be well described and work.
4. The design should be resilient with regard to data loss.
5. The application should be optimized for maximum throughput, weighed along with the other
   primary considerations and the requirements below.
6. The solution should be able to be build and run from the command line. Include specific
   instructions on dependencies, build, test, and run instructions. 
   
## Requirements

1. The application must accept input from at most 5 concurrent clients on TCP/IP port 4000.
2. Input lines presented to the application via its socket must either be composed of exactly nine
   decimal digits (e.g. 314159265 or 007007009) immediately followed by a server-native newline
   sequence; or a termination sequence as detailed below.
3. Numbers presented to the application must include leading zeros as necessary to ensure they
   are each 9 decimal digits.
4. The log file, to be named "numbers.log", must be created anew and/or cleared when the
   application starts.
5. Only numbers may be written to the log file. Each number must be followed by a server-native
   newline sequence.
6. No duplicate numbers may be written to the log file.
7. Any data that does not conform to a valid line of input should be discarded and the client
   connection terminated immediately and without comment.
8. Every 10 seconds, the application must print a report to standard output:
   a. The difference since the last report of the count of new unique numbers that have been
   received.
   b. The difference since the last report of the count of new duplicate numbers that have been
   received.
   c. The total number of unique numbers received for this run of the application.
   d. Example text: Received 50 unique numbers, 2 duplicates. Unique total: 567231
9. If any connected client writes a single line with only the word "terminate" followed by a server-native newline sequence, the Application must disconnect all clients and perform a clean
   shutdown as quickly as possible.
10. Clearly state all of the assumptions you made in completing the Application. 
    
## Notes

1. Ensure your application is executable from the command line.
2. Distribute your code with all the necessary instructions to build it and run it.
3. You should write tests that confirm the functionality of the application to the best of your
   ability.
4. Your Application may not for any part of its operation use or require the use of external
   systems, for example Apache Ka`a or Redis.
5. Leading zeroes should be stripped when writing to the log file and console.
6. Robust implementations of the application typically handle more than 2M numbers per 10-
   second reporting period on a modern MacBook Pro.