package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
)

var (
	flagGo      bool
	flagTradeID bool
)

func main() {
	flag.BoolVar(&flagGo, `go`, false, `output go syntax`)
	flag.BoolVar(&flagTradeID, `trade`, false, `output trade ids not uuids`)
	flag.Parse()

	nUUIDs := 1

	if args := flag.Args(); len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}

		nUUIDs = n
	}

	for ; nUUIDs > 0; nUUIDs-- {
		var id string
		if flagTradeID {
			id, _ = createTradeID()
		} else {
			id = uuid.Must(uuid.NewV4()).String()
		}

		if flagGo {
			fmt.Printf("`%s`,", id)
		} else {
			fmt.Print(id)
		}

		if nUUIDs > 1 {
			fmt.Println()
		}
	}
}

var alphabet = `1234567abcdefghijkmnopqrstuvwxyz`

func createTradeID() (string, error) {
	var randBytes [32]byte
	n, err := rand.Read(randBytes[:])
	if n != 32 {
		return "", errors.New("failed to read enough random bytes from rand for trade id")
	} else if err != nil {
		return "", fmt.Errorf("failed to generate trade id: %w", err)
	}

	// Read 5 bits at a time, convert to alphabet, we use 51*5 = 255 bits
	// (discarding one bit of randomness we pulled)
	var str [51]byte
	for i := uint(0); i < 51; i++ {
		byteOffset := (i * 5) / 8
		bitOffset := (i * 5) % 8

		// [0] 0, 1, 2, 3, 4
		// [0] 5, 6, 7 [1] 0, 1
		// [1] 2, 3, 4, 5, 6
		// [1] 7 [2] 1, 2, 3, 4

		// Bits 0-3 are okay to start at, but if we are starting at
		// any bit past that then we need to straddle a byte
		var charIndex byte
		if bitOffset > 3 {
			var firstByteMask byte = (1 << (8 - bitOffset)) - 1
			var secondByteMask byte = (1 << (5 - (8 - bitOffset))) - 1
			charIndex = (randBytes[byteOffset] >> bitOffset) & firstByteMask
			charIndex |= (randBytes[byteOffset+1] & secondByteMask) << (8 - bitOffset)
		} else {
			charIndex = (randBytes[byteOffset] >> bitOffset) & 0x1F
		}

		str[i] = alphabet[charIndex]
	}

	return string(str[:]), nil
}
