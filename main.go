package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

var threads uint64
var suffix string

func init() {
	flag.Uint64Var(&threads, "threads", 8, "parallel threads")
	flag.StringVar(&suffix, "suffix", "", "wallet address suffix")
}

func main() {
	flag.Parse()

	if suffix == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var counter uint64

	for t := uint64(0); t < threads; t++ {
		go func() {
			for {
				atomic.AddUint64(&counter, 1)

				seed := wallet.NewSeed()
				w, _ := wallet.FromSeed(nil, seed, wallet.V4R2)

				if strings.HasSuffix(w.Address().String(), suffix) {
					fmt.Printf("\nAddress: %v", w.Address().String())
					fmt.Printf("\nSeed phrase: %v\n\n", seed)
				}
			}
		}()
	}

	for {
		log.Printf("Searching, %v per second\n", counter)
		atomic.StoreUint64(&counter, 0)
		time.Sleep(1 * time.Second)
	}
}
