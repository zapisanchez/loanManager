package domain

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
