package entities

// usecase interface
type Usecase interface {
	CalculateTamboonSummary(file string) (*TamboonSummary, error)
}

// repositories interface
type Repositories interface {
	ReadFile(file string) ([]*SongPahPaDonor, error)
	CreateCreditCardToken(donorCardInfo *SongPahPaDonor) (*DonorCardTokenInfo, error)
	CreateChargeDonation(donorToken *DonorCardTokenInfo) (*DonorChargedInfo, error)
}
