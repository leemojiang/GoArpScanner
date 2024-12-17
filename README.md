# LAN Scanner in Go
Go网络编程练习

## Description
基于[Go Scan](https://github.com/timest/goscan)项目,修改代码以适配在windows上运行.

## Features
 * Scan the whole IPv4 address space
 * Scan your local network with ARP packets
 * Display the IP address, MAC address, ~~hostname~~ and vendor associated
 * ~~Using SMB(Windows devices) and mDNS(Apple devices) to detect hostname~~

## Modifications
 * ```handle.SetBPFFilter("arp or rarp")``` 在windows下失效,使用过滤器后无法接收到包.
 * `handle, err := pcap.OpenLive(pcapName, 1024, false, 20*time.Second)` 网卡打开名称在windows下需要修改.

## Usage
```sh
# install dependencies
$ go mod tidy

# build
$ go build

# or run
$ go run .
```
