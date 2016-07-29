package lib

import (
	"encoding/json"
	"fmt"
	"testing"
)

// IsEqualString checks strings for equalities and outputs the difference between them if not.
func IsEqualString(t *testing.T, expected string, actual string) {
	if expected == actual {
		return
	}
	expectedjs, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	actualjs, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}
	expr := []rune(expected)
	actr := []rune(actual)

	if len(expr) != len(actr) {
		t.Errorf("String lengths differ. expected: %d actual: %d\r\n", len(expr), len(actr))
	}
	sz := len(expr)
	if len(expr) > len(actr) {
		sz = len(actr)
	}
	for i := 0; i < sz; i++ {
		if expr[i] != actr[i] {
			fmt.Printf("chr #%d differs. e:'%c' a:'%c'\r\n", i, expr[i], actr[i])
		}
	}

	fmt.Printf("\r\n%s\r\n%s", map[string]string{"data": string(expectedjs)}, map[string]string{"data": string(actualjs)})
	t.Fail()
}
