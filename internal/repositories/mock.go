package repositories

import (
	"challenge_go/internal/entities"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) ReadFile(file string) ([]*entities.SongPahPaDonor, error) {
	args := mock.Called(file)
	return args.Get(0).([]*entities.SongPahPaDonor), args.Error(1)
}

func (mock *MockRepository) CreateCreditCardToken(donorCardInfo *entities.SongPahPaDonor) (*entities.DonorCardTokenInfo, error) {
	args := mock.Called(donorCardInfo)
	return args.Get(0).(*entities.DonorCardTokenInfo), args.Error(1)
}

func (mock *MockRepository) CreateChargeDonation(donorToken *entities.DonorCardTokenInfo) (*entities.DonorChargedInfo, error) {
	args := mock.Called(donorToken)
	return args.Get(0).(*entities.DonorChargedInfo), args.Error(1)
}
