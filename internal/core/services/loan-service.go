package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zapisanchez/loanMgr/internal/adapters/input"
	"github.com/zapisanchez/loanMgr/internal/core/domain"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
)

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
