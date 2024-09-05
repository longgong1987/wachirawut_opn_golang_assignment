package repositories

import (
	"bufio"
	"bytes"
	"challenge_go/cipher"
	"challenge_go/internal/entities"
	"challenge_go/services"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Repositories struct {
	Opn *services.OpnClient
}

func NewRepositories(opnClient *services.OpnClient) *Repositories {
	return &Repositories{
		Opn: opnClient,
	}
}

func (repo *Repositories) ReadFile(file string) ([]*entities.SongPahPaDonor, error) {

	var dataListWithOutHeader []*entities.SongPahPaDonor
	// read file
	tamboonBytes, err := os.Open(file)
	if err != nil {
		return dataListWithOutHeader, err
	}
	defer tamboonBytes.Close()

	scanner := bufio.NewScanner(tamboonBytes)

	var dataList string
	if scanner.Scan() {
		dataList = scanner.Text()
	}

	if dataList == "" {
		return dataListWithOutHeader, errors.New("no tamboon data")
	}

	rot128Reader, err := cipher.NewRot128Reader(bytes.NewBufferString(dataList))
	if err != nil {
		return dataListWithOutHeader, err
	}

	fileSize := len(dataList)

	buffer := make([]byte, fileSize)
	_, err = rot128Reader.Read(buffer)
	if err != nil {
		return dataListWithOutHeader, err
	}

	decryptDataList := strings.Split(string(buffer), "\n")

	for index, value := range decryptDataList[1:] {

		donorInfo := strings.Split(value, ",")

		if len(donorInfo) > 0 && donorInfo[0] != "" {

			// fmt.Printf("%v\n\n", donorInfo)

			amount, err := strconv.ParseInt(donorInfo[1], 10, 64)
			if err != nil {
				fmt.Printf("Error to convert amount sub unit row %d string to int64 %v", index, err)
			}

			ccNumber := donorInfo[2]
			cvv := donorInfo[3]

			expMonth, err := strconv.Atoi(donorInfo[4])
			if err != nil {
				fmt.Printf("Error to convert exp month row %d string to int64 %v", index, err)
			}

			expYear, err := strconv.Atoi(donorInfo[5])
			if err != nil {
				fmt.Printf("Error to convert exp year row %d string to int64 %v", index, err)
			}

			dataListWithOutHeader = append(dataListWithOutHeader, &entities.SongPahPaDonor{
				Name:           donorInfo[0],
				AmountSubunits: amount,
				CCNumber:       ccNumber,
				Cvv:            cvv,
				ExpMonth:       expMonth,
				ExpYear:        expYear,
			})

		}
	}

	return dataListWithOutHeader, nil

}

func (repo *Repositories) CreateCreditCardToken(donorCardInfo *entities.SongPahPaDonor) (*entities.DonorCardTokenInfo, error) {

	var donorCardTokenInfo entities.DonorCardTokenInfo

	payload := &services.CreateTokenPayload{
		Name:            donorCardInfo.Name,
		Number:          donorCardInfo.CCNumber,
		ExpirationMonth: donorCardInfo.ExpMonth,
		ExpirationYear:  donorCardInfo.ExpYear,
		SecurityCode:    donorCardInfo.Cvv,
	}
	card, err := repo.Opn.CreateToken(payload)
	if err != nil {
		return &donorCardTokenInfo, err
	}

	donorCardTokenInfo.ID = card.ID
	donorCardTokenInfo.Object = card.Object
	donorCardTokenInfo.Name = card.Name
	donorCardTokenInfo.Created = card.Created
	donorCardTokenInfo.Amount = donorCardInfo.AmountSubunits

	// clear memory allocated
	card = nil

	return &donorCardTokenInfo, nil
}

func (repo *Repositories) CreateChargeDonation(donorToken *entities.DonorCardTokenInfo) (*entities.DonorChargedInfo, error) {

	var donorChargedInfo entities.DonorChargedInfo

	payload := &services.CreateChargePayload{
		Amount:    donorToken.Amount,
		Currency:  "thb",
		CardToken: donorToken.ID,
	}
	chargedInfo, err := repo.Opn.CreateCharge(payload)
	if err != nil {
		return &donorChargedInfo, err
	}

	donorChargedInfo.Amount = chargedInfo.Amount
	donorChargedInfo.Created = chargedInfo.Created
	donorChargedInfo.ID = chargedInfo.ID
	donorChargedInfo.Object = chargedInfo.Object
	donorChargedInfo.Live = chargedInfo.Live
	donorChargedInfo.Location = *chargedInfo.Location
	donorChargedInfo.Status = string(chargedInfo.Status)
	donorChargedInfo.FailureCode = *chargedInfo.FailureCode
	donorChargedInfo.Currency = chargedInfo.Currency
	donorChargedInfo.Description = chargedInfo.Description

	// clear memory allocated
	chargedInfo = nil

	return &donorChargedInfo, nil
}
