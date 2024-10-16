# loanMgr - Loan Management Tool in Go

loanMgr is a command-line tool written in Go for managing loans, tracking payments, and calculating the time remaining until loans are fully paid off. It offers a simple interface to add new loans, track payments, view loan history, and perform calculations based on interest rates and payments.

## Features

- Add, track, and manage multiple loans.
- Log payments and calculate remaining balance.
- View loan payment history in a human-readable format.
- Calculate loan duration based on monthly payments and interest rate.
- Automatic saving and retrieval of loan data in JSON files.
- Deletion of completed loans with archiving to a separate directory.
- Handles loans with zero interest rates.
- User-friendly interface to interact with loans and payments.

## Installation

### Clone the repository and build

To install and run **loanMgr** locally:

```bash
git clone https://github.com/zapisanchez/loanMgr.git
cd loanMgr/cmd
go build -o loanMgr
```

### Install via Go

```bash
go install github.com/zapisanchez/loanMgr@latest
```

## Usage

```bash
./loanMgr
```

Once launched, you will see options to:

1) Show existing loans.
1) Create a new loan.
1) Add a payment to a loan.
1) View the payment history of a loan.
1) Exit.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
