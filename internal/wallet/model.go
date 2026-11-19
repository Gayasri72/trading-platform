package wallet

import "time"

type Balance struct {
	Asset     string `json:"asset"`
	Available string `json:"available"`
	Locked    string `json:"locked"`
	Total     string `json:"total"`
}

type LedgerEntry struct {
	ID          int64     `json:"id"`
	Asset       string    `json:"asset"`
	Amount      string    `json:"amount"`
	EntryType   string    `json:"entry_type"`
	ReferenceID string    `json:"reference_id"`
	CreatedAt   time.Time `json:"created_at"`
}
