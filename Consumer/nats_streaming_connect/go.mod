module nats-streaming-connect

go 1.21.1

require (
	Consume/memcache v0.0.0-00010101000000-000000000000
	github.com/nats-io/nats.go v1.34.1
)

require (
	Consume/order_struct v0.0.0-00010101000000-000000000000 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace Consume/order_struct => ../order_struct

replace Consume/memcache => ../memcache
