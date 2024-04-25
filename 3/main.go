package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	//"gonum.org/v1/plot/plotutil"
)

func psiATheory(psiAs, phi1 float64, x []float64) []float64 {
	psiA := make([]float64, len(x))
	for i := range x {
		psiA[i] = psiAs * math.Cosh(phi1*x[i]) / math.Cosh(phi1)
	}
	return psiA
}

func psiRTheory(psiAs, psiRs, phi1, phi2 float64, x []float64) []float64 {
	sigma := (phi1 * phi1) / (phi2 * phi2)
	psiR := make([]float64, len(x))
	for i := range x {
		coshPhi1 := math.Cosh(phi1*x[i]) / math.Cosh(phi1)
		coshPhi2 := math.Cosh(phi2*x[i]) / math.Cosh(phi2)
		psiR1 := psiRs * coshPhi2
		psiR2 := (sigma / (1.0 - sigma)) * psiAs * (coshPhi1 - coshPhi2)
		psiR[i] = psiR1 + psiR2
	}
	return psiR
}

func psiSTheory(psiA, psiR []float64) []float64 {
	psiS := make([]float64, len(psiA))
	for i := range psiA {
		psiS[i] = 1.0 - psiA[i] - psiR[i]
	}
	return psiS
}

func main() {
	// Define parameters
	psiAs := 0.6
	psiRs := 0.3
	phi1 := 4.0
	phi2 := 10.0
	zRange := make([]float64, 0)
	for z := -1.0; z < 1.0; z += 0.0001 {
		zRange = append(zRange, z)
	}

	// Calculate concentrations for different scenarios
	psiA1 := psiATheory(psiAs, phi1, zRange)
	psiR1 := psiRTheory(psiAs, psiRs, phi1, phi2, zRange)
	psiS1 := psiSTheory(psiA1, psiR1)

	createAndSavePlot(phi1, phi2, zRange, psiA1, psiR1, psiS1, 1)

	phi2 = 0.05
	psiA2 := psiATheory(psiAs, phi1, zRange)
	psiR2 := psiRTheory(psiAs, psiRs, phi1, phi2, zRange)
	psiS2 := psiSTheory(psiA2, psiR2)

	createAndSavePlot(phi1, phi2, zRange, psiA2, psiR2, psiS2, 2)

	phi1 = 4.0
	phi2 = 2.0
	psiA3 := psiATheory(psiAs, phi1, zRange)
	psiR3 := psiRTheory(psiAs, psiRs, phi1, phi2, zRange)
	psiS3 := psiSTheory(psiA3, psiR3)

	// Plotting

	createAndSavePlot(phi1, phi2, zRange, psiA3, psiR3, psiS3, 3)
}
func createAndSavePlot(phi1, phi2 float64, zRange, psiA, psiR, psiS []float64, fileNum int) {
	p := plot.New()

	p.Title.Text = "φ1=" + fmt.Sprintf("%v", phi1) + ", φ2=" + fmt.Sprintf("%v", phi2)
	p.X.Label.Text = "ξ [-]"
	p.Y.Label.Text = "ψ [-]"

	p.X.Min = -1
	p.X.Max = 1
	p.X.Tick.Marker = plot.DefaultTicks{}
	p.X.Padding = vg.Length(0)
	p.Y.Min = 0
	p.Y.Max = 1
	p.Y.Tick.Marker = plot.DefaultTicks{}
	p.Y.Padding = vg.Length(0)

	pts1 := make(plotter.XYs, len(zRange))
	pts2 := make(plotter.XYs, len(zRange))
	pts3 := make(plotter.XYs, len(zRange))
	for i := range zRange {
		pts1[i].X = zRange[i]
		pts1[i].Y = psiA[i]
		pts2[i].X = zRange[i]
		pts2[i].Y = psiR[i]
		pts3[i].X = zRange[i]
		pts3[i].Y = psiS[i]
	}

	lp1, _ := plotter.NewLine(pts1)
	lp1.LineStyle.Width = vg.Points(2)
	lp1.LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	lp2, _ := plotter.NewLine(pts2)
	//lp2.LineStyle.Dashes = []vg.Length{vg.Points(10), vg.Points(10)}
	lp2.LineStyle.Width = vg.Points(2)
	lp2.LineStyle.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}

	lp3, _ := plotter.NewLine(pts3)
	//lp3.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	lp3.LineStyle.Width = vg.Points(2)
	lp3.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}

	p.Add(lp1, lp2, lp3)

	p.Legend.Add("A", lp1)
	p.Legend.Add("R", lp2)
	p.Legend.Add("S", lp3)

	p.Save(12*vg.Inch, 12*vg.Inch, strconv.Itoa(fileNum)+".png")
}
