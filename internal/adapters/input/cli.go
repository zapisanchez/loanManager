package input

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zapisanchez/loanMgr/internal/core/domain"

	"github.com/inancgumus/screen"
)

func ClearScreen() {
	screen.Clear()
	screen.MoveTopLeft()
}

// GetUserInput prompts the user for input and returns the entered string.
func GetUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1] // Remove the newline at the end
}

// GetUserName prompts the user for a username and returns it.
func GetUserName() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the username:")
	scanner.Scan()
	return scanner.Text()
}

// GetLoanName prompts the user for the loan name.
func GetLoanName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the loan name: ")
	loanName, _ := reader.ReadString('\n')
	return loanName[:len(loanName)-1] // Remove the newline character
}

// GetPaymentAmount prompts the user for a payment amount.
func GetPaymentAmount() float64 {
	var amount float64
	fmt.Println("Enter the payment amount:")
	fmt.Scanln(&amount)
	return amount
}

// GetPaymentDescription prompts the user for a payment description.
func GetPaymentDescription() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the payment description:")
	desc, _ := reader.ReadString('\n')
	return desc
}

// GetInitialLoanAmount prompts the user for the initial loan amount.
func GetInitialLoanAmount() float64 {
	var amount float64
	fmt.Println("Enter the initial loan amount:")
	fmt.Scanln(&amount)
	return amount
}

// GetMonthlyPaymentAmount prompts the user for the monthly payment amount.
func GetMonthlyPaymentAmount() float64 {
	var amount float64
	fmt.Println("Enter the monthly payment amount:")
	fmt.Scanln(&amount)
	return amount
}

// GetInitialLoanAmount prompts the user for the initial loan amount.
func GetInterestRate() float64 {
	var rate float64
	fmt.Println("Enter the interest rate:")
	fmt.Scanln(&rate)
	return rate
}

// GetUserChoice prompts the user for their choice from the menu.
func GetUserChoice() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your choice: ")
	choice, _ := reader.ReadString('\n')
	return choice[:len(choice)-1] // Remove the newline character
}

// GetLoanSelection prompts the user to select a loan by index.
func GetLoanSelection(loans []domain.Loan) string {
	for {
		fmt.Println("Available LoanIDs:")
		for _, loan := range loans {
			fmt.Printf("LoanID: %s, Name: %s\n", loan.LoanID, loan.LoanName)
		}

		fmt.Println("Enter a LoanID to select a loan or type 'exit' to return to the main menu:")
		var selection string
		fmt.Scanf("%s", &selection)

		if selection == "exit" {
			return "" // Return to the main menu
		}

		// Check if the entered LoanID is valid
		for _, loan := range loans {
			if loan.LoanID == selection {
				return selection // return the selected LoanID
			}
		}

		fmt.Println("Invalid LoanID. Please try again.")
	}
}

// GetPaymentSelection prompts the user to select a payment by index. Return datetime of the selected payment.
func GetPaymentSelection(payments []domain.Payment) string {
	for {
		fmt.Println("Payment History:")
		for i, payment := range payments {
			fmt.Printf("%d. Date: %s, Description: %s, Amount: %.2f\n", i+1, payment.DateTime, payment.Description, payment.Amount)
		}

		fmt.Println("Enter the number of the payment to select it or type 'exit' to return to the main menu:")
		var selection int
		fmt.Scanf("%d", &selection)

		if selection == 0 {
			return "" // Return to the main menu
		}

		if selection < 1 || selection > len(payments) {
			fmt.Println("Invalid selection. Please try again.")
			continue
		}

		return payments[selection-1].DateTime
	}
}
