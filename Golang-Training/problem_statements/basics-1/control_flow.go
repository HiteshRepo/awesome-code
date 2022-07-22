package basics_1

import (
	"fmt"
	"time"
)

func FindDayDuration() {
	if h := time.Now().Hour(); h < 12 {
		fmt.Println("Now is AM time.")
	} else if h > 19 {
		fmt.Println("Now is evening time.")
	} else {
		fmt.Println("Now is afternoon time.")
		h := h // shadow
		_ = h
	}
}
