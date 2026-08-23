package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
)

type Player struct {
	ID    int
	score int
}

func (p *Player) Play(
	ctx context.Context,
	myTurn <-chan struct{},
	nextTurn chan<- struct{},
	wg *sync.WaitGroup,
	cancelFunc context.CancelFunc,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case <-myTurn:
			p.score += rand.Intn(7)

			fmt.Printf(
				"Player %d played, new score %d\n",
				p.ID,
				p.score,
			)

			if p.score >= 50 {
				fmt.Printf(
					"Player %d reached goal, ending game\n",
					p.ID,
				)

				cancelFunc()
				return
			}

			select {
			case nextTurn <- struct{}{}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p1Turn := make(chan struct{})
	p2Turn := make(chan struct{})

	p1 := Player{ID: 1}
	p2 := Player{ID: 2}

	wg.Add(2)

	go p1.Play(ctx, p1Turn, p2Turn, &wg, cancel)
	go p2.Play(ctx, p2Turn, p1Turn, &wg, cancel)

	// Player 1 starts
	p1Turn <- struct{}{}

	wg.Wait()

	fmt.Println(p1)
	fmt.Println(p2)
}
