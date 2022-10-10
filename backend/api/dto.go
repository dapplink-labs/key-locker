package api

import (
	"encoding/json"
	"time"
)

type KeyResp struct {
	KeyUuid   string    `json:"key_uuid"`
	KeySecret []string  `json:"key_secret"`
	Date      time.Time `json:"date" copier:"Timestamp"`
}

func (t KeyResp) MarshalJSON() ([]byte, error) {
	type Alias KeyResp
	return json.Marshal(&struct {
		Alias
		Date string `json:"date"`
	}{
		Alias: (Alias)(t),
		Date:  t.Date.Format("2006-01-02 15:04:05"),
	})
}
