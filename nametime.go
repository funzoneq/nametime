package main

import (
	"fmt"
	"os"
	"flag"
	"net"
	"time"
	"math/rand"
	"strings"
	"encoding/json"
)

var dnsServer string
var tsdbServer string
var domain string
var verbose, jsonout, tsdput bool
var dnsType uint16
var id uint16
var timesToCheck int
var result []lookupTimer

type lookupTimer struct {
	conntime, totaltime float64
}

type Message struct {
	Domain string
	DnsServer string
	AvgConnTime float64
	AvgTime float64
}

func init() {
	flag.StringVar(&dnsServer, "server", "8.8.8.8:53", "DNS server address (ip:port)")
	flag.StringVar(&domain, "domain", "tumblr.com", "The domain to resolv")
	flag.StringVar(&tsdbServer, "tsdserver", "opentsdb.example.com:4243", "Endpoint of OpenTSDB to report metrics")
	flag.IntVar(&timesToCheck, "timesToCheck", 25, "How many lookups do you want to base the avg on")
	flag.BoolVar(&jsonout, "jsonout", false, "Output json")
	flag.BoolVar(&tsdput, "tsdput", false, "Export metrics to OpenTSDB")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
}

func Average(xs []lookupTimer) (avgConnTime float64, avgTime float64) {
	// Set totals to 0
	total := float64(0)
	conntime := float64(0)
    
	// Foreach lookupTimer item
	for _, x := range xs {
		total += x.totaltime
		conntime += x.conntime
	}
	
	// Times of tries
	tries := float64(len(xs))
	
	// Get avg
	avgTime = total / tries
	avgConnTime = conntime / tries
	
	return
}

func tagsToString (tags map[string]string) string {
	slice := []string{}
	
	for k, v := range tags {
		slice = append(slice, fmt.Sprintf("%s=%s", k, v))
	}
	
	return strings.Join(slice, " ")
}

func tsd_put (c net.Conn, metric string, time int64, value float64, tags map[string]string) {
	// Convert tags to TSD format
	t := tagsToString(tags)
	
	// Format string for TSD PUT
	m := fmt.Sprintf("put %s %d %f %s\n", metric, time, value, t)
	
	// Send it to the server
	_, err := c.Write([]byte(m))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write(tcp): %s\n", err)
		os.Exit(1)
	}
}

func lookshitup (dnsServer string, domain string, id uint16, dnsType uint16) (conntime float64, ips []net.IP) {
	// Start connection timer
	t0 := time.Now()
	
	// Open connection, timeout in 1s
	c, err := net.DialTimeout("udp", dnsServer, time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bind(udp, %s): %s\n", dnsServer, err)
		os.Exit(1)
	}
	
	// End connection timer
	conntime = time.Now().Sub(t0).Seconds()
	
	// Make the query
	msg := packDns(domain, id, dnsType)
	
	// Send it to the server
	_, err = c.Write(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write(udp): %s\n", err)
		os.Exit(1)
	}
	
	// Make a result buffer
	buf := make([]byte, 4096)
	
	// Read the result
	n, err := c.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	
	// Close the connection
	c.Close()

	// Unpack the results, we only care about the ips
	_, _, ips = unpackDns(buf[:n], dnsType)
	
	return conntime, ips
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, strings.Join([]string{
			"\"resolve\" measures the performance of a DNS service provider.",
			"",
			"Usage: resolve [option ...]",
			"",
		}, "\n"))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	
	// Set a random seed
	rand.Seed(time.Now().UnixNano())
	
	// We want to check an A record
	dnsType = dnsTypeA
	
	for i := 0; i < timesToCheck; i++ {
		// Make a channel for comms
		c1 := make(chan lookupTimer, 1)
		
		go func() {
			t0 := time.Now()
			conntime, _ := lookshitup(dnsServer, domain, uint16(rand.Intn(100)), dnsType)
			c1 <- lookupTimer{ conntime, time.Now().Sub(t0).Seconds() }
		}()
		
		select {
		case res := <-c1:
			result = append(result, res)
		case <-time.After(time.Second * 1):
			result = append(result, lookupTimer{1, 1})
		}
	}
	
	avgConnTime, avgTime := Average(result)
	m := Message{Domain: domain, DnsServer: dnsServer, AvgConnTime: avgConnTime, AvgTime: avgTime}
	
	if jsonout {
		b, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	} else if tsdput {
		// Open connection, timeout in 5s
		c, err := net.DialTimeout("tcp", tsdbServer, time.Second * 5)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bind(tcp, %s): %s\n", tsdbServer, err)
			os.Exit(1)
		}
		
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		
		tags := map[string]string{}
		tags["host"] 		= hostname
		tags["dnsServer"] 	= strings.Split(dnsServer, ":")[0]
		tags["domain"] 		= domain
	
		tsd_put(c, "nametime.connTime", time.Now().Unix(), avgConnTime, tags)
		tsd_put(c, "nametime.resolveTime", time.Now().Unix(), avgTime, tags)
		
		// Close the connection
		c.Close()
	} else {
		fmt.Fprintf(os.Stderr, "Resolved %s with server %s in avgConnTime %.3fs and avgTime %.3fs.\n", m.Domain, m.DnsServer, m.AvgConnTime, m.AvgTime)
	}
}