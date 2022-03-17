#!/bin/bash
sudo sed -i 's/^dns=dnsmasq/#&/' /etc/NetworkManager/NetworkManager.conf
sudo service network-manager restart
sudo service networking restart
killall dnsmasq

myip=`ip route get 8.8.8.8 | head -1 | cut -d' ' -f8`
#[pdns]
docker run --name pdns -d -it -e API_KEY=1234 -p 8081:8081 -p $myip:53:53 -p $myip:53:53/udp pbertera/pdns

# example
#curl -X POST --data '{"name":"example.org.", "kind": "Native", "masters": [], "nameservers": ["ns1.example.org.", "ns2.example.org."]}' -v -H 'X-API-Key: openmcp' http://127.0.0.1:8081/api/v1/servers/localhost/zones
#curl -X PATCH --data '{"rrsets": [ {"name": "test.example.org.", "type": "A", "ttl": 86400, "changetype": "REPLACE", "records": [ {"content": "192.0.5.4", "disabled": false } ] } ] }' -H 'X-API-Key: openmcp' http://127.0.0.1:8081/api/v1/servers/localhost/zones/example.org
#dig @host -p 8853 test.example.org


#[pdns-admin] http://host:9191
docker run --name pdns-admin -d \
    -v pda-data:/data \
    -p 9191:80 \
    ngoduykhanh/powerdns-admin:latest
