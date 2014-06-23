Nametime
========

Measure performance of authoritative name servers

# Building

    Install go: http://golang.org/doc/install
    chmod +x build
    ./build

# Run as command-line tool

    ./nametime -server="ns1.p20.dynect.net:53" -domain="soundcloud.com" -timesToCheck=25
    
Output:

    Resolved soundcloud.com with server ns1.p20.dynect.net:53 in avgConnTime 0.002s and avgTime 0.015s.

# Export data as JSON

    ./nametime -server="ns1.yahoo.com:53" -domain="yahoo.com" -jsonout
    
Output:

    {"Domain":"yahoo.com","DnsServer":"ns1.yahoo.com:53","AvgConnTime":0.0011037932000000001,"AvgTime":0.013623636160000003}

# Export data to a OpenTSDB endpoint

    ./nametime -server="pdns1.ultradns.net:53" -domain="tumblr.com" -tsdbServer="opentsdb.example.com:4242" -tsdput
