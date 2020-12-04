# stping

[![Go Report](https://goreportcard.com/badge/github.com/natesales/stping?style=for-the-badge)](https://goreportcard.com/report/github.com/natesales/stping) 
[![License](https://img.shields.io/github/license/natesales/stping?style=for-the-badge)](https://raw.githubusercontent.com/natesales/stping/main/LICENSE) 
[![Release](https://img.shields.io/github/v/release/natesales/stping?style=for-the-badge)](https://github.com/natesales/stping/releases) 

A small utility to test ICMP latency from multiple source IP simultaneously.

### Installation
stping is available as a debian package and x86 binary in the releases section of this repo. It's also available as an APT package by adding `deb [trusted=yes] https://apt.fury.io/natesales/ /` to your `/etc/apt/source.list` file.

### Usage
```
Usage for stping (devel) https://github.com/natesales/stping:
  -sources string
        Comma separated list of source IP addresses
  -target string
        Target hostname to ping
```