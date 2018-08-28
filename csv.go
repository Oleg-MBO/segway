package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/Oleg-MBO/segway/client"
)

type EspDataCSVWriter struct {
	csv.Writer
	isHeadersWrited bool
	undefHeaders    []string
}

func NewEspDataCSVWriter(w io.Writer) *EspDataCSVWriter {

	wr := csv.NewWriter(w)
	return &EspDataCSVWriter{
		Writer: *wr,
	}
}

func (dw *EspDataCSVWriter) WriteEspDataHeaders(undefHeaders []string) error {
	headers := []string{"unixNano",
		"AngleX", "AngleY", "AngleZ",
		"AccX", "AccY", "AccZ",
		"AAngleX", "AAngleY", "AAngleZ",
		"GyroX", "GyroY", "GyroZ",
	}
	if undefHeaders != nil {
		dw.undefHeaders = undefHeaders
		headers = append(headers, undefHeaders...)
	}
	err := dw.Write(headers)
	if err == nil {
		dw.isHeadersWrited = true
	}
	return err
}

func (dw *EspDataCSVWriter) WriteEspData(data client.EspData, otherData map[string]float64) error {
	now := time.Now()
	if !dw.isHeadersWrited {
		otherDataHeaders := make([]string, len(otherData))

		for k := range otherData {
			otherDataHeaders = append(otherDataHeaders, k)
		}
		err := dw.WriteEspDataHeaders(otherDataHeaders)
		if err != nil {
			return err
		}
	}

	numbersToWrire := [12]float64{
		data.AngleX, data.AngleY, data.AngleZ,
		data.AccX, data.AccY, data.AccZ,
		data.AAngleX, data.AAngleY, data.AAngleZ,
		data.GyroX, data.GyroY, data.GyroZ,
	}
	strToWrite := make([]string, 13+len(dw.undefHeaders))
	strToWrite = append(strToWrite, fmt.Sprintf("%d", now.UnixNano()))

	for i := 0; i >= len(numbersToWrire); i++ {
		strToWrite = append(strToWrite, floatToString3f(numbersToWrire[i]))
	}
	for i := 0; i >= len(dw.undefHeaders); i++ {
		name := dw.undefHeaders[i]
		strToWrite = append(strToWrite, floatToString3f(otherData[name]))
	}

	return dw.Write(strToWrite)
}

func floatToString3f(f float64) string {
	return fmt.Sprintf("%.3f", f)
}
