package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"strings"
)

var desktopDir = fmt.Sprintf("%v/Desktop", getHomeDir())

func SavePdf(invoice Invoice) {
	var fileName = fmt.Sprintf("%v/%v_%v.pdf", desktopDir, strings.Replace(invoice.Client.Name, " ", "", -1), invoice.Date.Format("2006_01_02"))
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	leftMargin, _, _, _ := pdf.GetMargins()
	pageWidth, pageHeight := pdf.GetPageSize()
	pdf.SetFont("Arial", "", 12)
	list(pdf, leftMargin, 20, 0, 5, []string{
		"From:",
		invoice.MyInfo.BusinessName,
		invoice.MyInfo.Name,
		invoice.MyInfo.Street,
		invoice.MyInfo.CityStateZip,
	}, "R")
	list(pdf, leftMargin, 20, 0, 5, []string{
		"To:",
		invoice.VC.Name,
		invoice.VC.Street,
		invoice.VC.CityStateZip,
	}, "L")
	list(pdf, leftMargin, 50, 0, 5, []string{
		"Client:",
		fmt.Sprintf("%v", invoice.Client.Name),
		//invoice.Client.Street,
		//invoice.Client.CityStateZip,
		fmt.Sprintf("Claim #: %v", invoice.Client.ClaimNumber),
	}, "L")


	pdf.Rect(pageWidth / 2.0, 50, pageWidth / 2.0 - leftMargin, 40, "D")
	pdf.Rect(leftMargin, 100, pageWidth - leftMargin * 2, pageHeight - 100 - 25, "D")
	pdf.Rect(leftMargin, 100, pageWidth - leftMargin * 2, 8, "D")
	pdf.Rect(leftMargin, pageHeight - 35, pageWidth - leftMargin * 2, 10, "D")

	var dates = []string{"Date"}
	var descriptions = []string{"Description"}
	var amounts = []string{"Amount"}
	var total = 0.0
	for i := 0; i < len(invoice.LineItems); i++ {
		var li = invoice.LineItems[i]
		dates = append(dates, li.Date.Format("01/02"))
		descriptions = append(descriptions, li.Description)
		var amt = float64(li.Amount) / 100.0
		amounts = append(amounts, fmt.Sprintf("$%.2f", amt))
		total += amt
	}

	list(pdf, leftMargin, 100, 0, 8, dates, "L")
	list(pdf, leftMargin + 20, 100, 0, 8, descriptions, "L")
	list(pdf, leftMargin, 100, 0, 8, amounts, "R")

	pdf.SetY(-32.5)
	pdf.SetFontStyle("B")
	pdf.WriteAligned(0, 5, fmt.Sprintf("Total: $%.2f", total), "R")

	pdf.SetY(50)
	pdf.SetX(pageWidth / 2.0)
	pdf.Cell(pageWidth / 2.0 - leftMargin, 5, "Notes:")
	pdf.SetFontStyle("")
	pdf.SetY(55)
	pdf.SetX(pageWidth / 2.0)
	pdf.MultiCell(pageWidth / 2.0 - leftMargin, 5, invoice.Notes, "", "L", false)

	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		panic(err)
	}
}

func list(pdf *gofpdf.Fpdf, x, y, w, h float64, items []string, alignment string) {
	pdf.SetY(y)
	for i := 0; i < len(items); i++ {
		if i == 0 {
			pdf.SetFontStyle("B")
		} else {
			pdf.SetFontStyle("")
		}
		pdf.SetX(x)
		pdf.WriteAligned(w, h, items[i], alignment)
		pdf.Ln(h)
	}
}
