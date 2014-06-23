Nametime
========

Measure performance of authoritative name servers

# Building

    Install go: http://golang.org/doc/install
    chmod +x build
    ./build

# Run as command-line tool

    ./nametime -server="ns1.p20.dynect.net:53" -domain="soundcloud.com" -timesToCheck=25

# Export data as JSON

    ./nametime -server="ns1.yahoo.com:53" -domain="yahoo.com" -jsonout

# Export data to a OpenTSDB endpoint

    ./nametime -server="pdns1.ultradns.net:53" -domain="tumblr.com" -tsdbServer="opentsdb.example.com:4242" -tsdput
