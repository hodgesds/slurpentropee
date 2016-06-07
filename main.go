package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "slurpentropee"
	app.Usage = "Slurps entropy"
	app.Action = run

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "max,m",
			Value: -1,
			Usage: "Max int for generating RNG ints",
		},
		cli.IntFlag{
			Name:  "bytes,b",
			Value: 32,
			Usage: "Bytes to read from RNG",
		},
		cli.StringFlag{
			Name:  "timeout,t",
			Value: "5ms",
			Usage: "Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {
	timeout, err := time.ParseDuration(c.String("timeout"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if timeout < 0 {
		fmt.Println("Timeout must be >0")
		os.Exit(1)
	}

	readBytes := c.Int("bytes")
	readMax := c.Int("max")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	stopChan := make(chan bool)

	if readBytes > 0 {
		go func() {
			byteTicker := time.NewTicker(timeout).C
			for {
				select {
				case <-byteTicker:
					b := make([]byte, readBytes)
					_, err := rand.Read(b)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				case <-stopChan:
					return
				}
			}
		}()
	}

	if readMax > 0 {
		go func() {
			intTicker := time.NewTicker(timeout).C
			bigMax := big.NewInt(int64(readMax))
			for {
				select {
				case <-intTicker:
					_, err := rand.Int(rand.Reader, bigMax)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					count += 1
				case <-stopChan:
					return
				}
			}
		}()
	}

	<-signalChan
	close(stopChan)
}
