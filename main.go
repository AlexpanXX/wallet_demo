package main

import (
	"fmt"
	"os"
)

func main() {
	cfg := config{
		Magic: 20180627,
		SeedList: []string{
			"node-regtest-201.elastos.org",
			"node-regtest-202.elastos.org",
		},
		DefaultPort:22866,
	}

	wallet, err := newWallet(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	wallet.Start()

	select {}
}

//"Magic": 20180627,
//"PrintLevel": 0,
//"DefaultPort": 22866,
//"JsonRpcPort": 20477,
//"Foundation": "8ZNizBf4KhhPjeJRGpox6rPcHE5Np6tFx3",
//"SeedList": [
//"node-regtest-201.elastos.org",
//"node-regtest-202.elastos.org",
//"node-regtest-203.elastos.org",
//"node-regtest-204.elastos.org",
//"node-regtest-205.elastos.org",
//"node-regtest-206.elastos.org",
//"node-regtest-207.elastos.org"
//]
