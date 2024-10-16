package services

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/zapisanchez/loanMgr/internal/adapters/input"
	"github.com/zapisanchez/loanMgr/internal/core/domain"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
)

// AddPayment adds a new payment to the loan.
func AddPayment(loan *domain.Loan, amount float64, description string) {
	payment := domain.Payment{
		DateTime:    time.Now().Format("2006-01-02 15:04:05"),
		Amount:      amount,
		Description: description,
	}
	loan.Payments = append(loan.Payments, payment)
	loan.RemainingAmount -= amount
	loan.TotalPaid += amount

	RecalulatePayOff(loan)

	log.Info().Str("loan_id", loan.LoanID).Float64("amount", amount).Msg("Payment added")
}

// RemovePayment removes a payment from the loan.
func RemovePayment(loan *domain.Loan, paymentDate string) {

	// Find the payment index
	for i, payment := range loan.Payments {
		if payment.DateTime == paymentDate {
			paymentIndex := i
			payment := loan.Payments[paymentIndex]

			loan.RemainingAmount += payment.Amount
			loan.TotalPaid -= payment.Amount
			loan.Payments = append(loan.Payments[:paymentIndex], loan.Payments[paymentIndex+1:]...)
			RecalulatePayOff(loan)

			log.Info().Str("loan_id", loan.LoanID).Int("payment_index", paymentIndex).Msg("Payment removed")
			return
		}
	}
	log.Warn().Str("loan_id", loan.LoanID).Str("payment_date", paymentDate).Msg("Payment not found")
}

// ModifyPayment modifies a payment from the loan.
func ModifyPayment(loan *domain.Loan, paymentDate string, newAmount float64, newDescription string) {

	// Find the payment index
	for i, payment := range loan.Payments {
		if payment.DateTime == paymentDate {

			paymentIndex := i

			// As we are modifying the payment, we need to UPDATE
			payment := &loan.Payments[paymentIndex]

			// Remove the old amount from the total paid and remaining amount
			loan.RemainingAmount += payment.Amount
			loan.TotalPaid -= payment.Amount

			// Update the payment amount and description
			payment.Amount = newAmount
			payment.Description = newDescription

			// Update the total paid and remaining amount
			loan.RemainingAmount -= newAmount
			loan.TotalPaid += newAmount

			RecalulatePayOff(loan)

			log.Info().Str("loan_id", loan.LoanID).Int("payment_index", paymentIndex).Float64("new_amount", newAmount).Msg("Payment modified")
			return
		}
	}
	log.Warn().Str("loan_id", loan.LoanID).Str("payment_date", paymentDate).Msg("Payment not found")
}

// PrintLoanSummary prints the summary of a loan.
func PrintLoanSummary(loan domain.Loan) {
	fmt.Println("Loan ID:", loan.LoanID)
	fmt.Println("Initial Loan Amount:", loan.Amount)
	fmt.Println("Remaining Loan Amount:", loan.RemainingAmount)
	fmt.Println("Total Paid:", loan.TotalPaid)
	fmt.Println("Monthly Payment:", loan.MonthlyPayment)
	fmt.Println("Payments:")
	for _, payment := range loan.Payments {
		fmt.Printf(" - Date: %s, Amount: %s\n", payment.DateTime, strconv.FormatFloat(payment.Amount, 'f', 2, 64))
	}
}

// PrintAllLoans prints all loans for the user
func PrintAllLoans(loans []domain.Loan) {
	if len(loans) == 0 {
		input.ClearScreen()
		fmt.Println("No loans found.")
		log.Info().Msg("Press 'Enter' to go back to the main")
		_ = input.GetUserInput()
		input.ClearScreen()
		return
	}

	input.ClearScreen()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Loan Name",
		"Loan ID",
		"Amount",
		"Remaining Amount",
		"Total Paid",
		"Interest Rate",
		"Monthly Payment",
		"Months to Pay Off",
		"Years to Pay Off",
	})
	for _, loan := range loans {

		table.Append([]string{
			loan.LoanName,
			loan.LoanID,
			fmt.Sprintf("%.2f", loan.Amount),
			fmt.Sprintf("%.2f", loan.RemainingAmount),
			fmt.Sprintf("%.2f", loan.TotalPaid),
			fmt.Sprintf("%.2f", loan.Interest),
			fmt.Sprintf("%.2f", loan.MonthlyPayment),
			fmt.Sprintf("%.2f", loan.TimePaidOff),
			fmt.Sprintf("%.2f", loan.TimePaidOff/12),
		})

		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
			tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold, tablewriter.BgBlackColor},
			tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
		)

	}

	table.SetAutoFormatHeaders(true)
	// table.SetBorder(false)
	table.Render()

	// Ask the user if they want to go back to the main menu or exit
	log.Info().Msg("Press 'Enter' to go back to the main")
	_ = input.GetUserInput()
	input.ClearScreen()
}

// PrintPaymentHistory prints the payment history for a specific loan.
func PrintPaymentHistory(loan domain.Loan) {
	if len(loan.Payments) == 0 {
		fmt.Println("No payment history found for this loan.")
		return
	}

	input.ClearScreen()

	fmt.Printf("Payment history for Loan: %s (%s)\n", loan.LoanName, loan.LoanID)
	fmt.Println()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Description", "Amount"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor})

	for _, payment := range loan.Payments {
		table.Append([]string{payment.DateTime, payment.Description, fmt.Sprintf("%.2f €", payment.Amount)})
	}

	table.SetAutoFormatHeaders(true)
	table.SetFooter([]string{"", "Total Paid", fmt.Sprintf("%.2f €", loan.TotalPaid)})
	table.SetFooterColor(
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiRedColor})

	// table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetBorder(false)

	table.Render()

	totalTable := tablewriter.NewWriter(os.Stdout)
	totalTable.SetHeader([]string{"Total Paid", "Remaining Balance"})
	totalTable.Append([]string{fmt.Sprintf("%.2f €", loan.TotalPaid), fmt.Sprintf("%.2f €", loan.RemainingAmount)})
	totalTable.SetAutoFormatHeaders(true)
	totalTable.SetAlignment(tablewriter.ALIGN_RIGHT)
	totalTable.Render()

	fmt.Println()
}

func RecalulatePayOff(loan *domain.Loan) {

	if loan.Interest == 0 {

		months := loan.RemainingAmount / loan.MonthlyPayment
		loan.TimePaidOff = float64(months + 0.9999) // We round up
		return
	}

	// Convert the annual interest rate to a monthly interest rate
	monthlyInterestRate := loan.Interest / 12 / 100

	//Verify if the monthly payment is enough to cover the interest
	if loan.MonthlyPayment <= loan.RemainingAmount*monthlyInterestRate {
		log.Error().Msg("The monthly payment is too low to cover the interest.")
		return // Insufficient monthly payment
	}

	numerator := math.Log(loan.MonthlyPayment / (loan.MonthlyPayment - loan.RemainingAmount*loan.Interest))
	denominator := math.Log(1 + loan.Interest)

	loan.TimePaidOff = numerator / denominator
}
