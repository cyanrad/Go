package main

import (
	"fmt"
	"math"
)

const (
	R1       float64 = 99
	R2       float64 = 462
	RsTheory float64 = 264
)

func RsCurrent(IA, IB float64) float64 {
	return (R2 / ((IA / IB) - 1)) - R1
}

func RsVoltage(VA, VB float64) float64 {
	return (R1 * ((VB / VA) - 1)) / (1 - (VB/VA)*(R1/(R1+R2)))
}

func getErrPer(real, theory float64) float64 {
	return (math.Abs(real-theory) / math.Abs(theory)) * 100
}

func main() {
	RsC := RsCurrent(16.98*math.Pow(10, -3), 7.18*math.Pow(10, -3))
	RsV := RsVoltage(1.73, 4.31)

	fmt.Println("Rs from current measurements: ", RsC)
	fmt.Println("Error Value: ", getErrPer(RsC, RsTheory))
	fmt.Println("Rs from voltage measurements: ", RsV)
	fmt.Println("Error Value: ", getErrPer(RsV, RsTheory))
}
