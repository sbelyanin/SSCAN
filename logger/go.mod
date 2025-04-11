module github.com/sbelyanin/SSCAN/logger

go 1.24.2

replace github.com/sbelyanin/SSCAN/config => ../config

require (
	github.com/sbelyanin/SSCAN/config v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
)

require (
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
