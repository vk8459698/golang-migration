package ledger

import (
	"fmt"
	"time"
)

type Ledger struct {
	quoteModelInterval time.Duration
}

func NewLedger() *Ledger {
	fmt.Println("[Ledger] Starting Ledger initialization")

	l := &Ledger{
		quoteModelInterval: 5 * time.Second,
	}

	l.startQuoteModel()

	fmt.Println("[Ledger] Completed Ledger initialization")
	return l
}

func (l *Ledger) startQuoteModel() {
	go func() {
		for {
			start := time.Now()
			fmt.Println("[Ledger] Starting pool data update cycle")

			err := l.updatePoolData()
			if err != nil {
				fmt.Println("[Ledger] Pool data update failed:", err)
			} else {
				fmt.Println("[Ledger] Completed pool data update cycle")
			}

			fmt.Printf("[Ledger] Pool data update duration: %v\n", time.Since(start))
			time.Sleep(l.quoteModelInterval)
		}
	}()
}

func (l *Ledger) updatePoolData() error {
	time.Sleep(10 * time.Second)
	return nil
}
