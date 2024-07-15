# DNS Server

Stupid simple Domain Name System server.

## Features

1. Supports recursive lookups - starting from the root name servers.
2. Caches response to make subsequent queries faster. This way, query latency can be reduced to 0ms.
3. Zero dependencies.

## Example usage

Run the server locally on the `4321` port:

```
$ go run main.go
```

Open new terminal window and make a DNS query using `dig` command-line utility:

```
$ dig @localhost -p 4321 +noedns www.google.com
; <<>> DiG 9.18.26 <<>> @localhost -p 4321 +noedns www.google.com
; (2 servers found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 10586
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;www.google.com.			IN	A

;; ANSWER SECTION:
www.google.com.		300	IN	A	142.250.203.132

;; Query time: 96 msec
;; SERVER: ::1#4321(localhost) (UDP)
;; WHEN: Mon Jul 15 20:09:17 CEST 2024
;; MSG SIZE  rcvd: 48
```

Repeat the same query several times and observe lower latency:

```
$ dig @localhost -p 4321 +noedns www.google.com
...
;; Query time: 67 msec
...


$ dig @localhost -p 4321 +noedns www.google.com
...
;; Query time: 36 msec
...


$ dig @localhost -p 4321 +noedns www.google.com
...
;; Query time: 0 msec
...
```
