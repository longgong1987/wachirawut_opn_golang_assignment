package entities

import "time"

type SongPahPaDonor struct {
	Name           string
	AmountSubunits int64
	CCNumber       string
	Cvv            string
	ExpMonth       int
	ExpYear        int
}

type DonorCardTokenInfo struct {
	Object  string
	ID      string
	Name    string
	Created time.Time
	Amount  int64
}

type DonorChargedInfo struct {
	Object      string
	ID          string
	Live        bool
	Location    string
	Created     time.Time
	Status      string
	Amount      int64
	FailureCode string
	Currency    string
	Description *string
}

// model for send summary to presentation layer(controller)
type TamboonSummary struct {
	TotalReceived       float64
	DonatedSuccessfully float64
	DonationFaulty      float64
	AveragePerPerson    float64
	TopDonors           []string
}
