package jjj

import (
	"fmt"
	"test_test/rep"
)

func DoABad() error {
	err := fmt.Errorf("not good")
	return rep.EWL(err)
}
