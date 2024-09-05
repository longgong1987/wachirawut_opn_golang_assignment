package controllers

import (
	"challenge_go/internal/entities"
	"challenge_go/internal/usecases"
	"fmt"
	"strings"
)

type Controller struct {
	Uc entities.Usecase
}

func NewControllers(usecases *usecases.Usecases) *Controller {
	return &Controller{
		Uc: usecases,
	}
}

func (c *Controller) GetDonationSummary(filePath string) (string, error) {

	var summary *entities.TamboonSummary

	summary, err := c.Uc.CalculateTamboonSummary(filePath)
	if err != nil {
		return "", err
	}

	var summaryStr string

	summaryStr += fmt.Sprintf("      total received: THB %.2f\n", summary.TotalReceived)
	summaryStr += fmt.Sprintf("successfully donated: THB %.2f\n", summary.DonatedSuccessfully)
	summaryStr += fmt.Sprintf("     faulty donation: THB %.2f\n", summary.DonationFaulty)
	summaryStr += "\n"
	summaryStr += fmt.Sprintf("  average per person: THB %.2f\n", summary.AveragePerPerson)
	summaryStr += fmt.Sprintf("          top donors: %s\n", strings.Join(summary.TopDonors[:], "\n\t\t      "))

	return summaryStr, nil
}
