package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (alphabet = "abcdefghijklmnopqrstuvwxyz"
 	numbers= "1234567890"
)
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int,wantNumString bool) string {
	var sb strings.Builder
	k := len(alphabet)
	if wantNumString {

		k=len(numbers)
		for i := 0; i < n; i++ {
			c := numbers[rand.Intn(k)]
			sb.WriteByte(c)
		}
		return sb.String()
	}
	
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6,false)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD,INR}
	n := len(currencies)
	return currencies[rand.Intn(n)]
} 
// RandomMobile generates a random mobile number
func RandomMobile() string {
	return RandomString(10,true)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6,false))
}

