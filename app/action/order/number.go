package order

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
)

//nolint:gomnd // luhn algorithm obvious constants
func LuhnCheckDigit(s string) int {
	number, _ := strconv.Atoi(s)
	checkNumber := LuhnChecksum(number)

	if checkNumber == 0 {
		return 0
	}

	return 10 - checkNumber
}

//nolint:gomnd // luhn algorithm obvious constants
func LuhnChecksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur *= 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number /= 10
	}
	return luhn % 10
}

//nolint:gomnd // luhn algorithm obvious constants
func Number() string {
	ds := DigitString(5, 15)
	ids, _ := strconv.Atoi(ds)
	cs := LuhnChecksum(ids)
	cd := 0
	if cs > 0 {
		cd = 10 - cs
	}
	return fmt.Sprintf("%s%d", ds, cd)
}

func DigitString(minLen, maxLen int64) string {
	var letters = "0123456789"

	slen, _ := rand.Int(rand.Reader, big.NewInt(maxLen-minLen))
	slen = slen.Add(big.NewInt(minLen), slen)

	s := make([]byte, 0, slen.Int64())
	i := 0
	for int64(len(s)) < slen.Int64() {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))-1))
		char := letters[idx.Int64()]
		s = append(s, char)
		i++
	}

	return string(s)
}

var rule = regexp.MustCompile(`^\d+$`)

func ValidateOrderNr(orderNr string) error {
	if orderNr == "" {
		return fmt.Errorf("empty order")
	}

	if !rule.MatchString(orderNr) {
		return fmt.Errorf("orderNr has invalid format")
	}

	v := LuhnCheckDigit(orderNr[:len(orderNr)-1])
	l := orderNr[len(orderNr)-1]
	lv, _ := strconv.Atoi(string(l))

	if v != lv {
		return fmt.Errorf("invalid checksum")
	}

	return nil
}
