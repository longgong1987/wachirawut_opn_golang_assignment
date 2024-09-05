package usecases

import (
	"challenge_go/internal/entities"
	"challenge_go/internal/repositories"
	"errors"
	"reflect"
	"testing"
	"time"
)

var (
	mock    = &repositories.MockRepository{}
	usecase = Usecases{
		Repo: mock,
	}
)

func TestCalculateTamboonSummary(t *testing.T) {

	location := "th"
	timeNow := time.Now().UTC()

	tests := []struct {
		name                  string
		file                  string
		songPahPaDonors       []*entities.SongPahPaDonor
		songPahPaDonorsErr    error
		donorCardTokenInfo    *entities.DonorCardTokenInfo
		donorCardTokenInfoErr error
		charged               *entities.DonorChargedInfo
		chargedErr            error
		summary               *entities.TamboonSummary
		err                   error
	}{
		{
			name: "Successfully donation.",
			file: "file1.csv",
			songPahPaDonors: []*entities.SongPahPaDonor{
				{
					Name:           "Mr. Grossman R Oldbuck",
					AmountSubunits: 2879410,
					CCNumber:       "5375543637862918",
					Cvv:            "488",
					ExpMonth:       11,
					ExpYear:        2025,
				},
			},
			songPahPaDonorsErr: nil,
			donorCardTokenInfo: &entities.DonorCardTokenInfo{
				Object:  "1234",
				ID:      "test_tk_1234",
				Name:    "Mr. Grossman R Oldbuck",
				Created: time.Now().UTC(),
				Amount:  2879410,
			},
			donorCardTokenInfoErr: nil,
			charged: &entities.DonorChargedInfo{
				Object:      "1234",
				ID:          "test_tk_1234",
				Live:        false,
				Location:    location,
				Created:     timeNow,
				Status:      "successful",
				Amount:      2879410,
				FailureCode: "",
				Currency:    "thb",
			},
			chargedErr: nil,
			summary: &entities.TamboonSummary{
				TotalReceived:       28794.10,
				DonatedSuccessfully: 28794.10,
				DonationFaulty:      0,
				AveragePerPerson:    28794.10,
				TopDonors: []string{
					"Mr. Grossman R Oldbuck",
				},
			},
			err: nil,
		},
		{
			name: "Failed create token.",
			file: "file2.csv",
			songPahPaDonors: []*entities.SongPahPaDonor{
				{
					Name:           "Mr. Ferdinand H Took-Brandybuck",
					AmountSubunits: 2879410,
					CCNumber:       "5375543637862918",
					Cvv:            "488",
					ExpMonth:       11,
					ExpYear:        2025,
				},
			},
			songPahPaDonorsErr:    nil,
			donorCardTokenInfo:    &entities.DonorCardTokenInfo{},
			donorCardTokenInfoErr: errors.New("credit card has expired"),
			charged:               &entities.DonorChargedInfo{},
			chargedErr:            nil,
			summary: &entities.TamboonSummary{
				TotalReceived:       28794.10,
				DonatedSuccessfully: 0,
				DonationFaulty:      0,
				AveragePerPerson:    0,
				TopDonors:           []string{},
			},
			err: nil,
		},
		{
			name: "Failed create charge.",
			file: "file3.csv",
			songPahPaDonors: []*entities.SongPahPaDonor{
				{
					Name:           "Ms. Estella R Boffin",
					AmountSubunits: 1245821,
					CCNumber:       "5375543637862918",
					Cvv:            "488",
					ExpMonth:       11,
					ExpYear:        2025,
				},
			},
			songPahPaDonorsErr: nil,
			donorCardTokenInfo: &entities.DonorCardTokenInfo{
				Object:  "1234",
				ID:      "test_tk_1234",
				Name:    "Ms. Estella R Boffin",
				Created: time.Now().UTC(),
				Amount:  1245821,
			},
			donorCardTokenInfoErr: nil,
			charged:               &entities.DonorChargedInfo{},
			chargedErr:            errors.New("can not charge this token."),
			summary: &entities.TamboonSummary{
				TotalReceived:       12458.21,
				DonatedSuccessfully: 0,
				DonationFaulty:      12458.21,
				AveragePerPerson:    0,
				TopDonors:           []string{},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.On("ReadFile", tt.file).Return(tt.songPahPaDonors, tt.songPahPaDonorsErr).Once()

			for _, v := range tt.songPahPaDonors {
				mock.On("CreateCreditCardToken", v).Return(tt.donorCardTokenInfo, tt.donorCardTokenInfoErr).Once()
				mock.On("CreateChargeDonation", tt.donorCardTokenInfo).Return(tt.charged, tt.chargedErr).Once()
			}

			result, err := usecase.CalculateTamboonSummary(tt.file)
			if err != tt.err {
				t.Errorf("error = %v, wantErr %v", err, tt.err)
				return
			}
			if reflect.DeepEqual(result, tt.summary) == false {
				t.Errorf("summary not equal = %v, want %v", result, tt.summary)
			}
		})
	}

}
