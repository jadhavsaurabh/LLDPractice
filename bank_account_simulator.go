package main

import (
	"fmt"
	"sync"
	"time"
)

type Account struct {
	balance int
	mu      sync.Mutex
}

func (ac *Account) Deposit(amount int) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	fmt.Printf("\n Depositing amount %d with curr %d", amount, ac.balance)
	ac.balance += amount
	fmt.Printf("\n After deposit balance: %d", ac.balance)
}

func (ac *Account) Withdraw(amount int) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	fmt.Printf("\n Withdraw amount %d with curr %d", amount, ac.balance)
	if ac.balance-amount > 0 {
		ac.balance -= amount
		fmt.Printf("\n After withdrawing amount balance %d", ac.balance)
	} else {
		fmt.Printf("\n Insufficient balance %d", ac.balance)
	}
}

func main() {
	account := Account{
		balance: 3,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// deposit
	go func() {
		defer wg.Done()

		amounts := []int{4, 5, 1, 3, 8, 9}

		for _, amount := range amounts {
			account.Deposit(amount)
			time.Sleep(5 * time.Second)
		}
	}()

	// Withdraw
	go func() {
		defer wg.Done()

		amounts := []int{1, 4, 1, 3, 8, 10, 10, 10}

		for _, amount := range amounts {
			account.Withdraw(amount)
			time.Sleep(10 * time.Second)
		}
	}()

	wg.Wait()
}
