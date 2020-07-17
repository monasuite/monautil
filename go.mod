module github.com/monasuite/monautil

go 1.13

require (
	github.com/aead/siphash v1.0.1
	github.com/alecthomas/gometalinter v3.0.0+incompatible // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/btcsuite/golangcrypto v0.0.0-20150304025918-53f62d9b43e8
	github.com/davecgh/go-spew v1.1.1
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/kkdai/bstream v1.0.0
	github.com/monasuite/monad v0.21.0-beta
	github.com/monasuite/monautil/psbt v1.0.1 // indirect
	github.com/shopspring/decimal v1.2.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20191105091915-95d230a53780 // indirect
)

replace github.com/monasuite/monad => ../monad
