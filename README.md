# stping

[![Go Report](https://goreportcard.com/badge/github.com/natesales/stping?style=for-the-badge)](https://goreportcard.com/report/github.com/natesales/stping) 
[![License](https://img.shields.io/github/license/natesales/stping?style=for-the-badge)](https://raw.githubusercontent.com/natesales/stping/main/LICENSE) 
[![Release](https://img.shields.io/github/v/release/natesales/stping?style=for-the-badge)](https://github.com/natesales/stping/releases) 

A small utility to test ICMP latency from multiple source IP simultaneously.

### Example
```
~ » sudo stping -s "2001:xxx:1:xxx::2,2a04:xxxx:2:xxxx::4,2001:6xx:xxx::1" -t he.net
Resolving he.net...2001:470:0:503::2, 216.218.236.2
STPING he.net (2001:470:0:503::2) from 3 sources:
Source               Sent  Loss  Min       Max       Avg
2001:xxx:1:xxx::2    1    0.00%  1.255ms   1.255ms   1.255ms 
2a04:xxxx:2:xxxx::4  1    0.00%  1.098ms   1.098ms   1.098ms 
2001:6xx:xxx::1      1    0.00%  820µs     820µs     820µs   
2001:xxx:1:xxx::2    2    0.00%  1.255ms   1.284ms   1.269ms 
2a04:xxxx:2:xxxx::4  2    0.00%  1.098ms   1.191ms   1.144ms 
2001:6xx:xxx::1      2    0.00%  637µs     820µs     729µs   
2001:xxx:1:xxx::2    3    0.00%  1.105ms   1.284ms   1.214ms 
2a04:xxxx:2:xxxx::4  3    0.00%  1.031ms   1.191ms   1.107ms 
2001:6xx:xxx::1      3    0.00%  451µs     820µs     636µs   
^C
--- 2001:470:0:503::2 stping statistics source 2001:xxx:1:xxx::2 ---
4 packets transmitted, 4 packets received, 0% packet loss
round-trip min/avg/max/stddev = 1.01878ms/1.165911ms/1.284092ms/108.742µs

--- 2001:470:0:503::2 stping statistics source 2a04:xxxx:2:xxxx::4 ---
4 packets transmitted, 4 packets received, 0% packet loss
round-trip min/avg/max/stddev = 1.010895ms/1.083099ms/1.191288ms/70.386µs

--- 2001:470:0:503::2 stping statistics source 2001:6xx:xxx::1 ---
4 packets transmitted, 4 packets received, 0% packet loss
round-trip min/avg/max/stddev = 376.797µs/571.659µs/820.798µs/172.4µs
```

### Installation
This project is available in my public code repositories. See https://github.com/natesales/repo for more info.

### Usage
```
Usage for stping https://github.com/natesales/stping:
  -T uint
        Seconds to wait for a response (default 60)
  -i uint
        Seconds between pings (default 1)
  -s string
        Comma separated list of source IP addresses
  -t string
        Target hostname to ping
  -u    Use UDP instead of ICMP
```
