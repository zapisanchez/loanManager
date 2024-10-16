package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zapisanchez/loanMgr/internal/adapters/input"
	"github.com/zapisanchez/loanMgr/internal/adapters/repository"
	"github.com/zapisanchez/loanMgr/internal/core/domain"
	"github.com/zapisanchez/loanMgr/internal/core/services"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	input.ClearScreen()
	// Configure zerolog for output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Get the username
	userName := input.GetUserName()
	log.Info().Str("username", userName).Msg("User entered")

	// Load user data
	user, err := repository.LoadUser(userName)
	if err != nil {
		log.Warn().Err(err).Msg("No user found. Do you want create a new one? y/n.")
		// user = domain.User{UserName: userName}
		ch := input.GetUserChoice()
		if ch == "y" {
			createNewUser(userName)
			user, _ = repository.LoadUser(userName)
		} else {
			log.Info().Msg("Exiting the program.")
			return
		}
	}

	input.ClearScreen()

	// Main menu loop
	for {
		fmt.Println("Select an option:")
		fmt.Println("======= Loans =======")
		fmt.Println("1) Show loans")
		fmt.Println("2) Create a new loan")

		fmt.Println()
		fmt.Println("======= Payments =======")
		fmt.Println("3) Add a payment")
		fmt.Println("4) Modify a payment")
		fmt.Println("5) View payment history")

		fmt.Println()
		fmt.Println("6) Exit")
		choice := input.GetUserChoice()

		switch choice {
		case "1":
			services.PrintAllLoans(user.Loans) // A function to print all loans
		case "2":
			createNewLoan(&user) // Function to create a new loan
		case "3":
			addPaymentToLoan(&user) // Function to add payment to an existing loan
		case "4":
			modifyPaymentFromLoan(&user) // New function to modify a payment
		case "5":
			viewPaymentHistory(&user) // New function to view payment history
		case "6":
			log.Info().Msg("Exiting the program.")
			return // Exit the program
		default:
			log.Warn().Msg("Invalid choice. Please try again.")
		}
	}
}

func createNewUser(userName string) {
	log.Info().Msg("Creating a new user.")

	user := domain.User{UserName: userName}

	if err := repository.SaveUser(user); err != nil {
		log.Error().Err(err).Msg("Error saving user data")
	} else {
		log.Info().Msg("User data saved successfully")
	}
}

func createNewLoan(user *domain.User) {
	log.Info().Msg("Creating a new loan.")

	// Get loan details from the user
	loanName := input.GetLoanName()
	initialLoan := input.GetInitialLoanAmount()
	monthlyPayment := input.GetMonthlyPaymentAmount()
	interest := input.GetInterestRate()

	// Generate a unique LoanID
	loanID := generateUniqueLoanID(*user)

	// Create the new loan
	loan := domain.Loan{
		LoanID:          loanID,
		LoanName:        loanName,
		Amount:          initialLoan,
		RemainingAmount: initialLoan,
		TotalPaid:       0, // No payments made yet
		Interest:        interest,
		MonthlyPayment:  monthlyPayment,
		// TimePaidOff:     calculateMonthsUntilLoanPaidOff(initialLoan, monthlyPayment, interest),
	}

	//init to 0 payment and calculte payOff
	services.AddPayment(&loan, 0, "Initial Loan")

	user.Loans = append(user.Loans, loan)

	if err := repository.SaveUser(*user); err != nil {
		log.Error().Err(err).Msg("Error saving user data")
	} else {
		log.Info().Msg("User data saved successfully")
	}

	log.Info().Msg("New loan created")
}

func modifyPaymentFromLoan(user *domain.User) {

	// Select a loan to modify a payment
	loanID := input.GetLoanSelection(user.Loans)

	// If the user selects "exit", return to the main menu
	if loanID == "" {
		return
	}

	// Look for the loan by LoanID
	var selectedLoan *domain.Loan
	for i := range user.Loans {
		if user.Loans[i].LoanID == loanID {
			selectedLoan = &user.Loans[i]
			break
		}
	}

	if selectedLoan == nil {
		log.Warn().Msg("Loan not found.")
		return
	}

	paymentDate := input.GetPaymentSelection(selectedLoan.Payments)
	newAmount := input.GetPaymentAmount()
	newDesc := input.GetPaymentDescription()

	services.ModifyPayment(selectedLoan, paymentDate, newAmount, newDesc)

	// Save the user data
	if err := repository.SaveUser(*user); err != nil {
		log.Error().Err(err).Msg("Error saving user data")
	} else {
		log.Info().Msg("User data saved successfully")
	}
}

func addPaymentToLoan(user *domain.User) {
	if len(user.Loans) == 0 {
		log.Warn().Msg("No loans available to add payments.")
		return
	}

	// Select a loan to add a payment
	loanID := input.GetLoanSelection(user.Loans)

	// If the user selects "exit", return to the main menu
	if loanID == "" {
		return
	}

	// Look for the loan by LoanID
	var selectedLoan *domain.Loan
	for i := range user.Loans {
		if user.Loans[i].LoanID == loanID {
			selectedLoan = &user.Loans[i]
			break
		}
	}

	if selectedLoan == nil {
		log.Warn().Msg("Loan not found.")
		return
	}

	if selectedLoan.RemainingAmount == 0 {
		input.ClearScreen()

		log.Warn().Msg("Cannot add more payments to a loan with a total amount of 0.")

		// Ask the user if they want to go back to the main menu or exit
		log.Info().Msg("Press 'Enter' to go back to the main menu or type 'exit' to exit:")
		inputStr := input.GetUserInput()

		if inputStr == "exit" {
			log.Info().Msg("Exiting the program.")
			return
		} else {
			input.ClearScreen()
			return
		}
	}

	amount := input.GetPaymentAmount()
	description := input.GetPaymentDescription()
	services.AddPayment(selectedLoan, amount, description)

	log.Info().Float64("amount", amount).Msg("Payment added")

	// Save the user data
	if err := repository.SaveUser(*user); err != nil {
		log.Error().Err(err).Msg("Error saving user data")
	} else {
		log.Info().Msg("User data saved successfully")
	}

	if selectedLoan.RemainingAmount == 0 {
		log.Info().Msg("🎉🎉 Loan fully paid 🎉🎉")
	}
}

// Generate a unique LoanID based on existing loans
func generateUniqueLoanID(user domain.User) string {
	return strconv.Itoa(len(user.Loans) + 1) // Simple increment based on the number of existing loans
}

func viewPaymentHistory(user *domain.User) {
	if len(user.Loans) == 0 {
		log.Warn().Msg("No loans available to view payment history.")
		return
	}

	for {
		log.Info().Msg("Select a loan to view payment history:")
		loanID := input.GetLoanSelection(user.Loans) // Function to get user selection

		// If the user selects "exit", return to the main menu
		if loanID == "" {
			return
		}

		// Find the Loan by LoanID
		var selectedLoan domain.Loan
		for _, loan := range user.Loans {
			if loan.LoanID == loanID {
				selectedLoan = loan
				break
			}
		}

		services.PrintPaymentHistory(selectedLoan) // Print payment history for the selected loan

		// Ask the user if they want to go back to the main menu or exit
		log.Info().Msg("Press 'Enter' to go back to the main menu or type 'exit' to exit:")
		inputStr := input.GetUserInput()

		if inputStr == "exit" {
			log.Info().Msg("Exiting the program.")
			return
		} else {
			input.ClearScreen()
		}
	}
}
