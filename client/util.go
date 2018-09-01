package client

import (
	"fmt"
	"strconv"
	"strings"
)

func boolTo1or0(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func parceEspData(dataStr string) (espData EspData, err error) {
	dataStr = strings.Replace(dataStr, "\x00", "", -1)
	numbersSliseStr := strings.Split(dataStr, "|")
	if len(numbersSliseStr) < 12 {
		return EspData{}, fmt.Errorf(`separators "|" must be 11 and elements 12`)
	}
	millis, err := strconv.ParseUint(numbersSliseStr[0], 10, 64)
	if err != nil {
		return espData, err
	}
	espData.Milis = millis

	numbersSliseFloat := make([]float64, 0, 11)
	for i := 1; i < 13; i++ {
		f, err := strconv.ParseFloat(numbersSliseStr[i], 64)
		if err != nil {
			return espData, err
		}
		numbersSliseFloat = append(numbersSliseFloat, f)
	}

	espData.AngleX = numbersSliseFloat[1]
	espData.AngleY = numbersSliseFloat[2]
	espData.AngleZ = numbersSliseFloat[3]
	espData.AccX = numbersSliseFloat[4]
	espData.AccY = numbersSliseFloat[5]
	espData.AccZ = numbersSliseFloat[6]
	espData.GyroX = numbersSliseFloat[7]
	espData.GyroY = numbersSliseFloat[8]
	espData.GyroZ = numbersSliseFloat[9]
	espData.AAngleX = numbersSliseFloat[10]
	espData.AAngleY = numbersSliseFloat[11]
	espData.AAngleZ = 0
	return espData, nil
}
