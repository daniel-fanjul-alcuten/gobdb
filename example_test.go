package gobdb

import (
	"fmt"
)

func ExampleMemBurstRepository() {

	// the testRoot object keeps a counter
	bursts := NewMemBurstRepository()
	dispatcher := NewDefaultBurstDispatcher(bursts)
	defer dispatcher.Close()
	database := NewDefaultDatabase(&testRoot{}, dispatcher)

	// the testWriter increments the counter
	result1, _ := database.Write(&testWriter{3})
	fmt.Println("first write:", result1)

	// the testWriter decrements the counter
	result2, _ := database.Write(&testWriter{-1})
	fmt.Println("second write:", result2)

	// the testReader reads the counter
	result3 := database.Read(&testReader{})
	fmt.Println("read:", result3)
	// Output: first write: 3
	// second write: 2
	// read: 2
}
