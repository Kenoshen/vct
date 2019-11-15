package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
	"time"
)

func Initialize(current *Current) *Current {
	current.Data = &Data{
		MyInfo:   MyInfo{},
		VCs:      nil,
		Invoices: nil,
	}

	color.Yellow("Please fill out your personal info:")
	current.Data.MyInfo = *EditMyInfo(&current.Data.MyInfo, current.Reader)
	SaveData(*current.Data)

	color.Yellow("Please fill out the Victim's Compensation info:")
	current.Data.VCs = []VC{*EditVC(&VC{}, current.Reader)}

	return current
}

const menu_StartNewInvoice = "Start a new Invoice"
const menu_EditMyInfo = "Edit My Info"
const menu_EditVC = "Edit Victim's Compensation Info"
const menu_Exit = "Exit"

func Menu(current *Current) *Current {
	var options = []string{
		menu_StartNewInvoice,
		menu_EditMyInfo,
		menu_EditVC,
		menu_Exit,
	}
	var selectedIndex = captureOption("What would you like to do?", options, "Select one (leave blank to start invoice)", true, current.Reader)
	if selectedIndex < 0 {
		// default option
		return DisplayRecent(current)
	}

	var selected = options[selectedIndex]
	if selected == menu_StartNewInvoice {
		return DisplayRecent(current)
	} else if selected == menu_EditMyInfo {
		current.Data.MyInfo = *EditMyInfo(&current.Data.MyInfo, current.Reader)
		SaveData(*current.Data)
		return current
	} else if selected == menu_EditVC {
		return DisplayVCs(current)
	} else if selected == menu_Exit {
		os.Exit(0)
		return current
	} else {
		panic(fmt.Errorf("ERR: unknown option: %v", selected))
	}
}

func DisplayVCs(current *Current) *Current {
	var options []string
	for i := 0; i < len(current.Data.VCs); i++ {
		options = append(options, current.Data.VCs[i].Name)
	}
	var selectedIndex = captureOption("Victim's Compensations", options, "Select one (leave blank to create a new one)", true, current.Reader)
	if selectedIndex < 0 {
		var vc = EditVC(&VC{}, current.Reader)
		current.Data.VCs = append(current.Data.VCs, *vc)
	} else {
		current.Data.VCs[selectedIndex] = *EditVC(&current.Data.VCs[selectedIndex], current.Reader)
	}
	SaveData(*current.Data)
	return current
}

const recentInvoicesCount = 9

func DisplayRecent(current *Current) *Current {
	numOfInvoices := len(current.Data.Invoices)
	if numOfInvoices > 0 {
		var lastNInvoices []Invoice
		for i := numOfInvoices - 1; i >= 0; i-- {
			var inv = current.Data.Invoices[i]
			var add = true
			for k := 0; k < len(lastNInvoices); k++ {
				var kInv = lastNInvoices[k]
				if inv.Client.Name == kInv.Client.Name && inv.Client.ClaimNumber == kInv.Client.ClaimNumber {
					add = false
					break
				}
			}
			if add {
				lastNInvoices = append(lastNInvoices, inv)
			}
			if len(lastNInvoices) >= recentInvoicesCount {
				break
			}
		}

		var options []string
		for i := 0; i < len(lastNInvoices); i++ {
			options = append(options, lastNInvoices[i].Client.Name)
		}
		var option = captureOption("Recent", options, "Select one", true, current.Reader)
		if option < 0 {
			color.Yellow("Starting a new invoice...")
			return StartNewInvoice(current, Invoice{
				MyInfo:    current.Data.MyInfo,
				Client:    Client{},
				VC:        SelectVC(current.Data.VCs, current.Reader),
				Date:      time.Now(),
				Notes:     "",
				LineItems: nil,
			})
		} else {
			color.Yellow("Copying a previous invoice...")
			return StartNewInvoice(current, Invoice{
				MyInfo:    current.Data.MyInfo,
				Client:    lastNInvoices[option].Client,
				VC:        lastNInvoices[option].VC,
				Date:      time.Now(),
				Notes:     "",
				LineItems: nil,
			})
		}
	} else {
		color.Yellow("Starting a new invoice...")
		return StartNewInvoice(current, Invoice{
			MyInfo:    current.Data.MyInfo,
			Client:    Client{},
			VC:        SelectVC(current.Data.VCs, current.Reader),
			Date:      time.Now(),
			Notes:     "",
			LineItems: nil,
		})
	}
}

func StartNewInvoice(current *Current, invoice Invoice) *Current {
	EditClient(&invoice.Client, current.Reader)
	var addAnotherLineItem = true
	var lastLineItem *LineItem = nil
	for addAnotherLineItem {
		if lastLineItem == nil {
			lastLineItem = &LineItem{
				Date:        time.Now(),
				Description: "Personal therapy",
				Amount:      9000,
			}
		} else {
			lastLineItem = &LineItem{
				Date:        lastLineItem.Date,
				Description: lastLineItem.Description,
				Amount:      lastLineItem.Amount,
			}
		}
		EditLineItem(lastLineItem, current.Reader)
		invoice.LineItems = append(invoice.LineItems, *lastLineItem)
		addAnotherLineItem = captureBool("Add another", current.Reader)
	}

	color.Blue("Notes:")
	invoice.Notes = captureText("Enter notes", " ", current.Reader)

	var isOk = false

	for !isOk {
		PrintInvoice(invoice)
		isOk = captureBool("Is this ok?", current.Reader)
		if !isOk {
			invoice = *EditInvoice(&invoice, current.Reader)
		}
	}

	current.Data.Invoices = append(current.Data.Invoices, invoice)
	SaveData(*current.Data)

	SavePdf(invoice)
	color.Yellow("=================================================")
	color.Yellow("   This invoice has been saved to your Desktop")
	color.Yellow(fmt.Sprintf("            %v_%v_%v_%v", strings.Replace(invoice.Client.Name, " ", "", -1), invoice.Date.Year(), int(invoice.Date.Month()), invoice.Date.Day()))
	color.Yellow("=================================================")
	return current
}

func EditMyInfo(myInfo *MyInfo, reader *bufio.Reader) *MyInfo {
	color.Blue("My Info:")
	myInfo.BusinessName = captureText("Business name", myInfo.BusinessName, reader)
	myInfo.Name = captureText("Personal name", myInfo.Name, reader)
	myInfo.Street = captureText("Street", myInfo.Street, reader)
	myInfo.CityStateZip = captureText("City, State Zip", myInfo.CityStateZip, reader)

	return myInfo
}

func EditVC(vc *VC, reader *bufio.Reader) *VC {
	color.Blue("Victim's Compensation Organization:")
	vc.Name = captureText("Name", vc.Name, reader)
	vc.Street = captureText("Street", vc.Street, reader)
	vc.CityStateZip = captureText("City, State Zip", vc.CityStateZip, reader)

	return vc
}

func EditClient(client *Client, reader *bufio.Reader) *Client {
	color.Blue("Client:")
	client.Name = captureText("Client name", client.Name, reader)
	//client.Street = captureText("Street", client.Street, reader)
	//client.CityStateZip = captureText("City, State Zip", client.CityStateZip, reader)
	client.ClaimNumber = captureText("Claim number", client.ClaimNumber, reader)

	return client
}

func EditLineItem(lineItem *LineItem, reader *bufio.Reader) *LineItem {
	color.Blue("Line item:")
	lineItem.Date = captureDate("Date", lineItem.Date, reader)
	lineItem.Description = captureText("Description", lineItem.Description, reader)
	lineItem.Amount = captureMoney("Amount", lineItem.Amount, reader)

	return lineItem
}

func EditInvoice(invoice *Invoice, reader *bufio.Reader) *Invoice {
	invoice.Client = *EditClient(&invoice.Client, reader)
	invoice.LineItems = EditLineItems(invoice.LineItems, reader)
	invoice.Notes = captureText("Notes", invoice.Notes, reader)
	return invoice
}

func EditLineItems(lineItems []LineItem, reader *bufio.Reader) []LineItem {
	for true {
		var options []string
		for i := 0; i < len(lineItems); i++ {
			options = append(options, lineItems[i].ToString())
		}

		selectedIndex := captureOption("Line Items", options, "Which line item would you like to edit? (blank to skip)", true, reader)
		if selectedIndex < 0 {
			break
		} else {
			lineItems[selectedIndex] = *EditLineItem(&lineItems[selectedIndex], reader)
		}
	}
	return lineItems
}

func SelectVC(vcs []VC, reader *bufio.Reader) VC {
	var options []string
	for i := 0; i < len(vcs); i++ {
		options = append(options, vcs[i].Name)
	}
	var selectedIndex = captureOption("Victim's Compensation Organization", options, "Select one", false, reader)
	return vcs[selectedIndex]
}

//noinspection GoUnhandledErrorResult
func PrintInvoice(invoice Invoice) {
	c := color.New(color.BgHiBlue)
	c.Println()
	c.Println(invoice.Client.Name)
	c.Println(invoice.Client.ClaimNumber)
	c.Println()
	for i := 0; i < len(invoice.LineItems); i++ {
		var li = invoice.LineItems[i]
		c.Println(li.ToString())
	}
	if len(strings.Replace(invoice.Notes, " ", "", -1)) > 0 {
		c.Println()
		c.Println(invoice.Notes)
	}
	c.DisableColor()
	c.Println()
}
