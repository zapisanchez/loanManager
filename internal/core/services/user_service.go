package services

import (
	"errors"

	"github.com/zapisanchez/loanMgr/internal/core/domain"
)

type UserRepo interface {
	// Data from Memory
	GetUser(userName string) *domain.User
	AddUser(user *domain.User) error
	MoveUserToDeleted(userID string) error

	// Data from File
	PersistUserData() error
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(userName string) (*domain.User, error) {

	// return err if user already exists
	usr := s.repo.GetUser(userName)
	if usr != nil {
		return nil, errors.New("user already exists")
	}

	user := domain.NewUser(userName)

	// Create new user in memory
	err := s.repo.AddUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUser(userName string) *domain.User {
	return s.repo.GetUser(userName)
}

func (s *UserService) DeleteUser(userName string) error {

	// return err if user already exists
	usr := s.repo.GetUser(userName)
	if usr == nil {
		return errors.New("user not found")
	}

	// return err if user already exists
	err := s.repo.MoveUserToDeleted(userName)
	if err == nil {
		return errors.New("user not found")
	}

	err = s.repo.MoveUserToDeleted(userName)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AddLoanToUser(userName string, loan domain.Loan) error {
	user := s.repo.GetUser(userName)
	if user == nil {
		return errors.New("user not found")
	}
	user.AddLoan(loan)
	return nil
}

func (s *UserService) AddPaymentToLoan(userName string, loanID string, payment domain.Payment) error {
	user := s.repo.GetUser(userName)
	if user == nil {
		return errors.New("user not found")
	}

	selectedLoan := user.GetLoan(loanID)
	if selectedLoan == nil {
		return errors.New("loan not found")
	}

	if selectedLoan.RemainingAmount == 0 {
		return errors.New("loan fully paid")
	}

	selectedLoan.AddPayment(payment)

	return nil
}

func (s *UserService) ModifyPaymentFromLoan(userName string, loanID string, paymentDate string, newAmount float64, newDescription string) error {
	user := s.repo.GetUser(userName)
	if user == nil {
		return errors.New("user not found")
	}

	selectedLoan := user.GetLoan(loanID)
	if selectedLoan == nil {
		return errors.New("loan not found")
	}

	selectedLoan.ModifyPayment(paymentDate, newAmount, newDescription)
	return nil
}

func (s *UserService) Persist() error {
	return s.repo.PersistUserData()
}
