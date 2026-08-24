package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	MaxChairs = 3
	CutTime   = 100 * time.Millisecond
)

type Customer struct {
	id   int
	done chan struct{}
}

type BarberShop struct {
	mu            sync.Mutex
	waitingChairs int

	customerReady chan Customer
	quit          chan struct{}
	barberDone    chan struct{}
}

func (s *BarberShop) barber() {
	defer close(s.barberDone)

	for {
		fmt.Println("Barber is sleeping...")

		select {
		case customer := <-s.customerReady:
			s.mu.Lock()
			s.waitingChairs--
			s.mu.Unlock()

			fmt.Printf("Barber is cutting customer %d\n", customer.id)

			time.Sleep(CutTime)

			fmt.Printf("Haircut done for customer %d\n", customer.id)
			close(customer.done)

		case <-s.quit:
			fmt.Println("Barber is going home")
			return
		}
	}
}

func (s *BarberShop) customer(id int) {
	fmt.Printf("Customer %d entered\n", id)

	s.mu.Lock()

	if s.waitingChairs >= MaxChairs {
		fmt.Printf("Customer %d left - no chair available\n", id)
		s.mu.Unlock()
		return
	}

	s.waitingChairs++
	fmt.Printf("Customer %d is waiting (%d/%d)\n",
		id, s.waitingChairs, MaxChairs)

	s.mu.Unlock()

	c := Customer{
		id:   id,
		done: make(chan struct{}),
	}

	select {
	case s.customerReady <- c:
	case <-s.quit:
		s.mu.Lock()
		s.waitingChairs--
		s.mu.Unlock()
		return
	}

	<-c.done

	fmt.Printf("Customer %d left\n", id)
}

func main() {
	shop := &BarberShop{
		customerReady: make(chan Customer),
		quit:          make(chan struct{}),
		barberDone:    make(chan struct{}),
	}

	go shop.barber()

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			shop.customer(id)
		}(i)

		time.Sleep(60 * time.Millisecond)
	}

	wg.Wait()

	close(shop.quit)
	<-shop.barberDone

	fmt.Println("Shop closed")
}
