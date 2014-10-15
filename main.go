package main

import (
	"fmt"
	"github.com/aprilsky/goutils/systool"
	"github.com/aprilsky/goutils/timetool"
	"time"
)

func main() {
	fmt.Println(systool.IntranetIP())
	fmt.Println(timetool.DateFormat(time.Now(), "YYYY-MM-DD"))
}
