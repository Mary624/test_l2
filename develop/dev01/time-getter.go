package ntptime

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func WriteTime() {
	time, err := getTime()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(-1)
	}
	fmt.Println(time)
}

func getTime() (time.Time, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	return time, err
}
