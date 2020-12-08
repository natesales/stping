package main

import (
	"flag"
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var release = "" // Set by the build process

var (
	sources = flag.String("s", "", "Comma separated list of source IP addresses")
	target  = flag.String("t", "", "Target hostname to ping")
)

func main() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage for stping (%s) https://github.com/natesales/stping:\n", release)
		flag.PrintDefaults()
	}

	if *sources == "" || *target == "" {
		flag.Usage()
		os.Exit(1)
	}

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

	if len(targetHostRecordsStr) == 1 {
		targetIp = targetHostRecordsStr[0]
	} else {
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
	}

	longestSourceAddress := 0
	var pingers []*ping.Pinger

	addresses := strings.Split(strings.ReplaceAll(*sources, " ", ""), ",")

	for _, address := range addresses {
		// Create pinger
		pinger, err := ping.NewPinger(targetIp)
		pinger.Source = address
		if err != nil {
			log.Fatalf("Creating pinger: %v", err)
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

		// Start a new pinger goroutine
		go func() {
			err = pinger.Run()
			if err != nil {
				log.Fatalf("Running pinger: %v", err)
			}
		}()
	}

	// Listen for Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			for _, pinger := range pingers {
				fmt.Printf("\n--- %s stping statistics source %s ---\n", pinger.Statistics().Addr, pinger.Source)
				fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n", pinger.Statistics().PacketsSent, pinger.Statistics().PacketsRecv, pinger.Statistics().PacketLoss)
				fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", pinger.Statistics().MinRtt, pinger.Statistics().AvgRtt, pinger.Statistics().MaxRtt, pinger.Statistics().StdDevRtt)
				pinger.Stop()
			}

			os.Exit(0)
		}
	}()

	fmt.Printf("STPING %s (%s) from %d sources:\n", *target, targetIp, len(pingers))
	fmt.Printf("%s  Sent     Loss  Min        Max        Avg        Dev\n", ("Source" + strings.Repeat(" ", longestSourceAddress))[:longestSourceAddress])
	for {
		for _, pinger := range pingers {
			if pinger.PacketsSent > 0 {
				source := (pinger.Source + strings.Repeat(" ", longestSourceAddress))[:longestSourceAddress]
				fmt.Printf("%v  %-3v   %6.2f%%  %-9v  %-9v  %-9v  %-9v\n", source, pinger.PacketsSent, pinger.Statistics().PacketLoss, pinger.Statistics().MinRtt.Truncate(time.Microsecond), pinger.Statistics().MaxRtt.Truncate(time.Microsecond), pinger.Statistics().AvgRtt.Truncate(time.Microsecond), pinger.Statistics().StdDevRtt.Truncate(time.Microsecond))
			}
		}
		time.Sleep(time.Second)
	}
}
