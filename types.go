package main

import (
	"encoding/json"
	"strconv"
)

type Coin struct {
	Denom  string `json:"denom"`
	Amount uint64 `json:"amount,string"`
}

type BalanceResponse struct {
	Balances []Coin `json:"balances"`
}

type SizeResponse struct {
	Size uint64 `json:"size,string"`
}

type Count struct {
	Total uint64 `json:"total,string"`
}

type PageResponse struct {
	Pagination Count `json:"pagination"`
}

type StatsResponse struct {
	Purchased   uint64            `json:"purchased,string"`
	Used        uint64            `json:"used,string"`
	UsedRatio   float64           `json:"used_ratio,string"`
	ActiveUsers uint64            `json:"activeUsers,string"`
	TotalUsers  uint64            `json:"uniqueUsers,string"`
	UsersByPlan map[uint64]uint64 `json:"users_by_plan"`
}

func (s *StatsResponse) UnmarshalJSON(data []byte) error {
	// Create a raw structure to hold the fields.
	type Alias StatsResponse
	aux := &struct {
		UsersByPlan map[string]string `json:"users_by_plan"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	// Unmarshal into the temporary structure
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert map keys and values from string to uint64
	s.UsersByPlan = make(map[uint64]uint64)
	for k, v := range aux.UsersByPlan {
		key, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			return err
		}
		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		s.UsersByPlan[key] = value
	}

	return nil
}
