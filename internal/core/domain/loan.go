package domain

import (
	"math"

	"github.com/rs/zerolog/log"
)

// Structure to represent a user with multiple loans
type User struct {
	UserName string `json:"user_name"`
	Loans    []Loan `json:"loans"`
}

// Structure to represent a loan
type Loan struct {
	LoanID          string    `json:"loan_id"`
	LoanName        string    `json:"loan_name"`
	Amount          float64   `json:"amount"`           // Initial loan amount
	RemainingAmount float64   `json:"remaining_amount"` // Remaining amount to be paid
	TotalPaid       float64   `json:"total_paid"`       // Total amount paid
	Interest        float64   `json:"interest"`         // Interest rate
	MonthlyPayment  float64   `json:"monthly_payment"`  // Estimated Monthly payment amount
	TimePaidOff     float64   `json:"time_paid_off"`    // Time to pay off the loan
	Payments        []Payment `json:"payments"`         // Payment history
}

// Structure for each payment in the history
type Payment struct {
	DateTime    string  `json:"date_time"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

func NewUser(userName string) User {
	return User{
		UserName: userName,
		Loans:    []Loan{},
	}
}

func NewLoan(loanID, loanName string, amount, interest, monthlyPayment float64) Loan {
	return Loan{
		LoanID:          loanID,
		LoanName:        loanName,
		Amount:          amount,
		RemainingAmount: amount,
		TotalPaid:       0,
		Interest:        interest,
		MonthlyPayment:  monthlyPayment,
		TimePaidOff:     0,
		Payments:        []Payment{},
	}
}

func (u *User) AddLoan(loan Loan) {
	u.Loans = append(u.Loans, loan)
}

func (u *User) RemoveLoan(loanID string) {
	for i, loan := range u.Loans {
		if loan.LoanID == loanID {
			u.Loans = append(u.Loans[:i], u.Loans[i+1:]...)
			return
		}
	}
}

func (u *User) GetLoans() []Loan {
	return u.Loans
}
func (u *User) GetLoan(loanID string) *Loan {
	for i, loan := range u.Loans {
		if loan.LoanID == loanID {
			return &u.Loans[i]
		}
	}
	return nil
}

func (l *Loan) AddPayment(payment Payment) {
	l.Payments = append(l.Payments, payment)
	l.RemainingAmount -= payment.Amount
	l.TotalPaid += payment.Amount

	l.recalculatePayOff()
	log.Info().Str("loan_id", l.LoanID).Float64("amount", l.Amount).Msg("Payment added")
}

func (l *Loan) GetPayments() []Payment {
	return l.Payments
}

func (l *Loan) GetPayment(paymentDate string) *Payment {
	for i, payment := range l.Payments {
		if payment.DateTime == paymentDate {
			return &l.Payments[i]
		}
	}
	return nil
}

func (l *Loan) RemovePayment(paymentDate string) {
	for i, payment := range l.Payments {
		if payment.DateTime == paymentDate {
			paymentIndex := i
			payment := l.Payments[paymentIndex]

			l.RemainingAmount += payment.Amount
			l.TotalPaid -= payment.Amount
			l.Payments = append(l.Payments[:paymentIndex], l.Payments[paymentIndex+1:]...)
			l.recalculatePayOff()

			log.Info().Str("loan_id", l.LoanID).Int("payment_index", paymentIndex).Msg("Payment removed")
			return
		}
	}
}

func (l *Loan) ModifyPayment(paymentDate string, newAmount float64, newDescription string) {
	for i, payment := range l.Payments {
		if payment.DateTime == paymentDate {
			paymentIndex := i

			// As we are modifying the payment, we need to UPDATE
			payment := &l.Payments[paymentIndex]

			// Remove the old amount from the total paid and remaining amount
			l.RemainingAmount += payment.Amount
			l.TotalPaid -= payment.Amount

			// Update the payment amount and description
			payment.Amount = newAmount
			payment.Description = newDescription

			// Update the total paid and remaining amount
			l.RemainingAmount -= newAmount
			l.TotalPaid += newAmount

			l.recalculatePayOff()

			log.Info().Str("loan_id", l.LoanID).Int("payment_index", paymentIndex).Float64("new_amount", newAmount).Msg("Payment modified")
			return
		}
	}
	log.Warn().Str("loan_id", l.LoanID).Str("payment_date", paymentDate).Msg("Payment not found")
}

func (l *Loan) recalculatePayOff() {
	if l.Interest == 0 {

		months := l.RemainingAmount / l.MonthlyPayment
		l.TimePaidOff = float64(months + 0.9999) // We round up
		return
	}

	// Convert the annual interest rate to a monthly interest rate
	monthlyInterestRate := l.Interest / 12 / 100

	//Verify if the monthly payment is enough to cover the interest
	if l.MonthlyPayment <= l.RemainingAmount*monthlyInterestRate {
		log.Error().Msg("The monthly payment is too low to cover the interest.")
		return // Insufficient monthly payment
	}

	numerator := math.Log(l.MonthlyPayment / (l.MonthlyPayment - l.RemainingAmount*l.Interest))
	denominator := math.Log(1 + l.Interest)

	l.TimePaidOff = numerator / denominator
}
