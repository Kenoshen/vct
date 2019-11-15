package main

import (
	"bufio"
	"fmt"
	"time"
)

type MyInfo struct {
	Name         string `json:"name"`
	BusinessName string `json:"businessName"`
	Street       string `json:"street"`
	CityStateZip string `json:"cityStateZip"`
}

type Client struct {
	Name         string `json:"name"`
	ClaimNumber  string `json:"claimNumber"`
}

type VC struct {
	Name         string `json:"name"`
	Street       string `json:"street"`
	CityStateZip string `json:"cityStateZip"`
}

func (c *Client) PrintName() {
	fmt.Print(c.Name)
}

type LineItem struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
}

func (li *LineItem) ToString() string {
	return fmt.Sprintf("%v/%v - %v - $%.2f", int(li.Date.Month()), li.Date.Day(), li.Description, float64(li.Amount)/100.0)
}

type Invoice struct {
	MyInfo    MyInfo     `json:"myInfo"`
	VC        VC         `json:"vc"`
	Client    Client     `json:"client"`
	Date      time.Time  `json:"date"`
	Notes     string     `json:"notes"`
	LineItems []LineItem `json:"lineItems"`
}

type Data struct {
	MyInfo   MyInfo    `json:"myInfo"`
	VCs      []VC      `json:"vcs"`
	Invoices []Invoice `json:"invoices"`
}

type Current struct {
	Data   *Data
	Reader *bufio.Reader
}
