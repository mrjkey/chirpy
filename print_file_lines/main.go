package main

import (
	"fmt"

	"test_test/jjj"
	"test_test/rep"
)

func main() {
	err := doingMoreThings()
	if err != nil {
		fmt.Println(err)
	}
	err = doSomething()
	if err != nil {
		fmt.Println(err)
	}
	err = jjj.DoABad()
	if err != nil {
		fmt.Println(err)
	}
}

func doSomething() error {
	return rep.EWL("something went wrong")
}

func doingMoreThings() error {
	str := "I have a very odd feeling"
	for i, c := range str {
		if i%2 == 0 {
			if c > 1 {
				return rep.EWL("Something bad")
			}
		}
	}
	return nil
}
