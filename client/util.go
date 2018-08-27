package client

import (
	"fmt"
	"regexp"
	"strconv"
)

var reSendAnglesData = regexp.MustCompile(`\w{1} (.?\d+.\d+)`)

func parceXYZ(data string) (X, Y, Z float64, err error) {
	submatch := reSendAnglesData.FindAllStringSubmatch(data, -1)

	if len(submatch) != 3 {
		err = fmt.Errorf("error parseXYZ with data %s", data)
		return
	}
	X, err = strconv.ParseFloat(submatch[0][1], 64)
	if err != nil {
		return
	}
	Y, err = strconv.ParseFloat(submatch[1][1], 64)
	if err != nil {
		return
	}
	Z, err = strconv.ParseFloat(submatch[2][1], 64)

	return
}
