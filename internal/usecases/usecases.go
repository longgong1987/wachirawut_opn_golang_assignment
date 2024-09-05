package usecases

import (
	"challenge_go/internal/entities"
	"challenge_go/internal/repositories"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Usecases struct {
	Repo entities.Repositories
}

func NewUsecases(repo *repositories.Repositories) *Usecases {
	return &Usecases{
		Repo: repo,
	}
}

func (uc *Usecases) CalculateTamboonSummary(file string) (*entities.TamboonSummary, error) {

	var (
		summary              entities.TamboonSummary
		totalReceived        float64 = 0
		donatedSuccessfully  float64 = 0
		donationFaulty       float64 = 0
		successDonationCount int64   = 0
	)

	// get all Song-pah-pa donors info
	songPahPaDonors, err := uc.Repo.ReadFile(file)
	if err != nil {
		return &summary, err
	}

	var wg sync.WaitGroup
	channel := make(chan struct{}, 5)

	for i, v := range songPahPaDonors {

		if i%2 == 1 {
			time.Sleep(2000 * time.Millisecond)
		}

		totalReceived += float64(v.AmountSubunits) / 100.0

		channel <- struct{}{}
		wg.Add(1)

		go func(donor *entities.SongPahPaDonor) {
			defer wg.Done()

			// create donor token
			donorCardTokenInfo, err := uc.Repo.CreateCreditCardToken(donor)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			} else {
				if donorCardTokenInfo.ID != "" {

					// create charge from token
					charged, err := uc.Repo.CreateChargeDonation(donorCardTokenInfo)
					if err != nil {
						fmt.Printf("%s\n", err.Error())
					}

					// charged status
					// - confirmed_amount_mismatch
					// - failed_fraud_check
					// - failed_processing
					// - insufficient_balance / insufficient_fund
					// - invalid_account_number / invalid_account
					// - payment_cancelled
					// - payment_rejected
					// - stolen_or_lost_card
					// - timeout
					if charged.Status == "successful" && charged.FailureCode == "" {
						// add successfully charged
						donatedSuccessfully += float64(donor.AmountSubunits) / 100.0
						successDonationCount += 1
					} else {
						// add fault charged
						donationFaulty += float64(donor.AmountSubunits) / 100.0
					}

				}
			}

			<-channel
		}(v)

	}

	wg.Wait()

	summary.TotalReceived = totalReceived
	summary.DonatedSuccessfully = donatedSuccessfully
	summary.DonationFaulty = donationFaulty

	var (
		topDonors    []string = []string{}
		avgPerPerson float64  = 0
	)
	if donatedSuccessfully > 0 {
		avgPerPerson = uc.getAverageDonation(donatedSuccessfully, successDonationCount)
		topDonors = uc.getTopTreeDonors(songPahPaDonors)

	}
	summary.AveragePerPerson = avgPerPerson
	summary.TopDonors = topDonors

	songPahPaDonors = nil

	return &summary, nil
}

func (uc *Usecases) getAverageDonation(donatedSuccess float64, donationCount int64) float64 {
	if donatedSuccess == 0 || donationCount == 0 {
		return 0
	}
	return donatedSuccess / float64(donationCount)
}

func (uc *Usecases) getTopTreeDonors(songPhapa []*entities.SongPahPaDonor) []string {

	var topDonors []string

	sort.Slice(songPhapa, func(i, j int) bool {
		return songPhapa[i].AmountSubunits < songPhapa[j].AmountSubunits
	})
	for i, v := range songPhapa {
		if i < 3 {
			// fmt.Println(v)
			topDonors = append(topDonors, v.Name)
		} else {
			break
		}
	}

	return topDonors
}
