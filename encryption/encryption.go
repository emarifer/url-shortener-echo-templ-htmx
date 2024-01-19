package encryption

import (
	// "fmt"
	"crypto/rand"
	"errors"
)

// defaultAlphabet is the alphabet used for ID characters by default.
var defaultAlphabet = []rune("_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const (
	defaultSize = 21
)

// CreateSlug generates secure URL-friendly unique ID.
// Accepts optional parameter - length of the ID to be generated
// (21 by default).
func CreateSlug(l ...int) (string, error) {
	var size int
	switch {
	case len(l) == 0:
		size = defaultSize
	case len(l) == 1:
		size = l[0]
		if size < 0 {
			return "", errors.New("negative id length")
		}
	default:
		return "", errors.New("unexpected parameter")
	}

	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// fmt.Println(bytes)
	id := make([]rune, size)
	for i := 0; i < size; i++ {
		// fmt.Println(bytes[i], "-", bytes[i]&63)
		id[i] = defaultAlphabet[bytes[i]&63] // ↓ see note below ↓
	}

	return string(id[:size]), nil
}

/* NOTE:
The expression bytes[i]&63 calculates the modulus 64 of the uint8 value
(bytes[i]) so that it is always inside the defaultAlphabet rune array,
which is more efficient than calculating the modulus directly. SEE:

https://www.geeksforgeeks.org/find-the-remainder-when-n-is-divided-by-4-using-bitwise-and-operator/
https://stackoverflow.com/questions/3072665/bitwise-and-in-place-of-modulus-operator

ID Collision Calculator:
https://zelark.github.io/nano-id-cc/
*/
