package main

import (
	"flag"
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"net"
	"strings"
	"time"
)

var (
	sources = flag.String("sources", "", "Comma separated list of source IP addresses")
	target  = flag.String("target", "", "Target hostname to ping")
)

func main() {
	flag.Parse()

	fmt.Printf("Resolving %s...", *target)
	targetHostRecords, err := net.LookupIP(*target)
	if err != nil {
		log.Fatalf("Unable to resolve hostname %s; %v", *target, err)
	}
	var targetHostRecordsStr []string
	for _, addr := range targetHostRecords {
		targetHostRecordsStr = append(targetHostRecordsStr, addr.String())
	}
	fmt.Println(strings.Join(targetHostRecordsStr, ", "))

	var targetIp string

	// If IPv6 source
	if strings.Contains(*sources, ":") {
		for _, addr := range targetHostRecords {
			if strings.Contains(addr.String(), ":") {
				targetIp = addr.String()
				break
			}
		}
	} else {
		for _, addr := range targetHostRecords {
			if strings.Contains(addr.String(), ".") {
				targetIp = addr.String()
				break
			}
		}
	}

	longestSourceAddress := 0
	var pingers []*ping.Pinger

	addresses := strings.Split(strings.ReplaceAll(*sources, " ", ""), ",")

	for _, address := range addresses {
		// Create pinger
		pinger, err := ping.NewPinger(targetIp)
		pinger.Source = address
		if err != nil {
			panic(err)
		}

		// Add the pinger pointer to the array
		pingers = append(pingers, pinger)

		// Update the longest source address if applicable
		if len(address) > longestSourceAddress {
			longestSourceAddress = len(address)
		}

		// Enable privileged mode
		pinger.SetPrivileged(true)

		//table[&ip].OnRecv = func(pkt *ping.Packet) {
		//	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		//}

		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s stping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}

		// Start a new pinger goroutine
		go func() {
			err = pinger.Run()
			if err != nil {
				panic(err)
			}
		}()
	}

	fmt.Printf("STPING %s (%s) from %d sources:\n", *target, targetIp, len(pingers))

	//// Listen for Ctrl-C
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//for _, pinger := range table {
	//	for range c {
	//		pinger.Stop()
	//	}
	//}

	fmt.Printf("%s  Sent  Loss  Min        Max        Avg\n", ("Source" + strings.Repeat(" ", longestSourceAddress))[:longestSourceAddress])
	for {
		for _, pinger := range pingers {
			if pinger.PacketsSent > 0 {
				source := (pinger.Source + strings.Repeat(" ", longestSourceAddress))[:longestSourceAddress]
				fmt.Printf("%v  %-3v %5.2f%%  %-v  %-v  %-v\n", source, pinger.PacketsSent, pinger.Statistics().PacketLoss, pinger.Statistics().MinRtt.Truncate(time.Microsecond), pinger.Statistics().MaxRtt.Truncate(time.Microsecond), pinger.Statistics().AvgRtt.Truncate(time.Microsecond))
			}
		}
		time.Sleep(time.Second)
	}
}
