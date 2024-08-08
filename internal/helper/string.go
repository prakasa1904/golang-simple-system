package helper

import (
	"math/rand"
	"strings"
	"time"
)

// random from timenano
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// whitelisted char to generate random string
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generate random string with specific length
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateUsernameFromEmail(email string) string {
	splEmail := strings.Split(email, "@")
	return strings.Replace(splEmail[0], ".", "", -1) + RandomString(5)
}
