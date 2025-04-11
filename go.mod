module github.com/sbelyanin/SSCAN

go 1.24.2

replace github.com/sbelyanin/SSCAN/config => ./config

replace github.com/sbelyanin/SSCAN/metrics => ./metrics

replace github.com/sbelyanin/SSCAN/logger => ./logger

replace github.com/sbelyanin/SSCAN/scanner => ./scanner

replace github.com/sbelyanin/SSCAN/server => ./server

require (
	github.com/sbelyanin/SSCAN/config v0.0.0-20250411131609-dab5c159b8af
	github.com/sbelyanin/SSCAN/logger v0.0.0-00010101000000-000000000000
	github.com/sbelyanin/SSCAN/metrics v0.0.0-00010101000000-000000000000
	github.com/sbelyanin/SSCAN/scanner v0.0.0-00010101000000-000000000000
	github.com/sbelyanin/SSCAN/server v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/sys v0.30.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
