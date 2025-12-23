package utils


import (
	"encoding/json"
	"fmt"
	"time"
	"math/rand/v2"
)


func NowUTCString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func LogPretty(title string, v interface{}) {

	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("* ", title)
	fmt.Println("------------------------------------------------------------------------")

	// Handle error types explicitly
	if err, ok := v.(error); ok {
		fmt.Println(err.Error())
		// return
	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal object:", err)
		return
	}
	fmt.Println(string(bytes))
	fmt.Println("_____")
}

func RandRange(min, max int) int {
    return rand.IntN(max-min) + min
}
func Uint64() uint64 {
    return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}
