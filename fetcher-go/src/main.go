package main

import (
	"fetcher-go/src/common"
	"fetcher-go/src/ledger"
)

func initApp() {
	common.InitRedis()
	ledger.NewLedger()
}

func main() {
	initApp()

	select {}
}
