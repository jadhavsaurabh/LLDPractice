import (
	"fmt"
	"strings"
	"sync"
)

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func (s *SafeMap) set(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = s.data[key] + 1
}

func NewSafeMap() *SafeMap {
	mp := map[string]int{}
	return &SafeMap{data: mp}
}

func main() {
	stringsInput := [3]string{"abc ded", "pqr abc", "abc abc ded"}

	var wg sync.WaitGroup

	res := NewSafeMap()
	for _, currStr := range stringsInput {
		wg.Add(1)
		go processString(currStr, res, &wg)
	}

	wg.Wait()
	fmt.Print(res)
}

func processString(currStr string, res *SafeMap, wg *sync.WaitGroup) {
	defer wg.Done()
	listStr := strings.Fields(currStr)
	fmt.Println(listStr)
	for _, str := range listStr {
		res.set(str)
	}
}
