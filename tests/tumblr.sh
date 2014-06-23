#!/bin/bash

# UltraDNS
./nametime -timesToCheck=25 -server="pdns1.ultradns.net:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="pdns2.ultradns.net:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="pdns3.ultradns.org:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="pdns4.ultradns.org:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="pdns5.ultradns.info:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="pdns6.ultradns.co.uk:53" -domain="tumblr.com" -tsdput

# Dyn
./nametime -timesToCheck=25 -server="ns1.p03.dynect.net:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="ns2.p03.dynect.net:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="ns3.p03.dynect.net:53" -domain="tumblr.com" -tsdput
./nametime -timesToCheck=25 -server="ns4.p03.dynect.net:53" -domain="tumblr.com" -tsdput

# Route53
./nametime -timesToCheck=25 -server="ns-9.awsdns-01.com:53" -domain="photoset.com" -tsdput
./nametime -timesToCheck=25 -server="ns-952.awsdns-55.net:53" -domain="photoset.com" -tsdput
./nametime -timesToCheck=25 -server="ns-1747.awsdns-26.co.uk:53" -domain="photoset.com" -tsdput
./nametime -timesToCheck=25 -server="ns-1384.awsdns-45.org:53" -domain="photoset.com" -tsdput