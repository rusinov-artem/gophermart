package order

import (
	"fmt"
	"math/rand"
	"strconv"
)

func LuhnCheckDigit(s string) int {
	number, _ := strconv.Atoi(s)
	checkNumber := LuhnChecksum(number)

	if checkNumber == 0 {
		return 0
	}

	return 10 - checkNumber
}

func LuhnChecksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func OrderNr() string {
	ds := DigitString(5, 15)
	ids, _ := strconv.Atoi(ds)
	cs := LuhnChecksum(ids)
	cd := 0
	if cs > 0 {
		cd = 10 - cs
	}
	return fmt.Sprintf("%s%d", ds, cd)
}

func DigitString(minLen, maxLen int) string {
	var letters = "0123456789"

	slen := rand.Intn(maxLen-minLen) + minLen

	s := make([]byte, 0, slen)
	i := 0
	for len(s) < slen {
		idx := rand.Intn(len(letters) - 1)
		char := letters[idx]
		s = append(s, char)
		i++
	}

	return string(s)
}
