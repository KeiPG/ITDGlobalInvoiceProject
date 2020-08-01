package main

import (

	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"math"
)
//employee data
type Customer struct{
	ID int
	name string
	address string
	currency_ID int
}
type Currency struct{
	Date string
	FromCurrency string
	ToCurrency []ToCurrency
}
type ToCurrency []struct {
	Currency string  
	Rate     float64 
}
type currencyRates struct {
	Rates struct {
		CAD float64 `json:"CAD"`
		HKD float64 `json:"HKD"`
		ISK float64 `json:"ISK"`
		PHP float64 `json:"PHP"`
		DKK float64 `json:"DKK"`
		HUF float64 `json:"HUF"`
		CZK float64 `json:"CZK"`
		AUD float64 `json:"AUD"`
		RON float64 `json:"RON"`
		SEK float64 `json:"SEK"`
		IDR float64 `json:"IDR"`
		INR float64 `json:"INR"`
		BRL float64 `json:"BRL"`
		RUB float64 `json:"RUB"`
		HRK float64 `json:"HRK"`
		JPY float64 `json:"JPY"`
		THB float64 `json:"THB"`
		CHF float64 `json:"CHF"`
		SGD float64 `json:"SGD"`
		PLN float64 `json:"PLN"`
		BGN float64 `json:"BGN"`
		TRY float64 `json:"TRY"`
		CNY float64 `json:"CNY"`
		NOK float64 `json:"NOK"`
		NZD float64 `json:"NZD"`
		ZAR float64 `json:"ZAR"`
		USD float64 `json:"USD"`
		MXN float64 `json:"MXN"`
		ILS float64 `json:"ILS"`
		GBP float64 `json:"GBP"`
		KRW float64 `json:"KRW"`
		MYR float64 `json:"MYR"`
	} `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}
//main function 
func main() {


	response, err := http.Get("https://api.exchangeratesapi.io/latest")
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	
	var Rates currencyRates

	err = json.Unmarshal(body,&Rates)
	// this was my attempt at the setting up the system to auto have the correct values for the currency rates 
	// var temp = [3]ToCurrency
	// temp[0].Currency = 'USD'
	// temp[0].Rate = 1.31
	// temp[1].Currency = 'EUR'
	// temp[1].Rate = 1.11
	// temp[2].Currency = 'JPY'
	// temp[2].Rate = 138.12
	// var Currencies [5]Currency
	// Currencies[1].date = Now()
	// Currencies[1].FromCurrency = "GBP"
	// Currencies[1].ToCurrency = temp
	// for i :=1; i <=4; i++ {
	// 	fmt.Println(Currencies[i].GBP)
	// 	fmt.Println(Currencies[i].EUR)
	// 	fmt.Println(Currencies[i].USD)
	// 	fmt.Println(Currencies[i].JPY)
	// 	// currentCurrency := '{"GBP":%1,"USD":%2,"EUR":%3,"JPY":%4}',Currencies[i].GBP,Currencies[i].USD,Currencies[i].EUR,Currencies[i].JPY
	// 	// stmt.QueryRow(currentCurrency,i)
	// }

	
	
	//Variables required for setup
	/*
	user= (using default user for postgres database)
	dbname= (using default database that comes with postgres)
	password = (password used during initial setup)
	host = (IP Address of server)
	sslmode = (must be set to disabled unless using SSL. This is not covered during tutorial)
	*/

	//DO NOT SAVE PASSWORD AS TEXT IN A PRODUCTION ENVIRONMENT. TRY USING AN ENVIRONMENT VARIABLE
	connStr := "user=postgres dbname=invoice_db password= host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Print(err)
	}
	var choice string
	fmt.Println("1.Update the Amount on an Invoice ")
	fmt.Println("2.Get all invoices by Customer ID ")
	fmt.Println("3.Get all Invoices by Currency Name")
	fmt.Print("enter the number of the option you want:")
	fmt.Scanf("%s",&choice)
	if choice == "1"{
		statement :="UPDATE Invoices set Amount = $1 WHERE ID = $2;"
		//prepare statement for sql
		stmt , err := db.Prepare(statement)
		if err != nil {
			fmt.Print(err)
		}
		var invID string
		var AMOUNT string
		defer stmt.Close()
		fmt.Print("Please Enter a Invoice ID: ")
		fmt.Scanf("%s",&invID)
		fmt.Print("Please Enter The Amount: ")
		fmt.Scanf("%s",&AMOUNT)
		stmt.QueryRow(AMOUNT,invID)
	}



	// this allowed me to add alot of Customers quickly
	// statement :="INSERT INTO Customers(Name, Address, Currency_ID) VALUES($1, $2, $3)"
	// stmt , err := db.Prepare(statement)
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// defer stmt.Close()
	// eCustomer := Customer{}
	// for i :=0; i <3; i++ {
	// 	fmt.Print("Name: ")
	// 	fmt.Scanf("%s",&eCustomer.name)
	// 	fmt.Print("Address: ")
	// 	fmt.Scanf("%s",&eCustomer.address)
	// 	fmt.Print("Currency: ")
	// 	fmt.Scanf("%s",&eCustomer.currency_ID)
	// 	stmt.QueryRow(eCustomer.name,eCustomer.address,eCustomer.currency_ID)
	// }

	if choice == "2"{
		rows, err := db.Query("SELECT* FROM invoices JOIN Customers ON invoices.customer_id = Customers.id JOIN Currencies ON Customers.Currency_ID = CUrrencies.id ORDER BY customer_id;")
		if err != nil {
			fmt.Print(err)
		}
		defer rows.Close()
	
		fmt.Println("---------------------------------------------------------------------")
		//loop through all employee results
		for rows.Next(){
			//assign values to variables
			var inv_ID int
			var invoice_creation_date string
			var Cus_ID int
			var amount float64
			var cus_id int
			var name string
			var address string
			var currency_ID int
			var cur_id int
			var cur_name string
			var exchange_rate string
			err := rows.Scan(&inv_ID, &invoice_creation_date, &Cus_ID, &amount,&cus_id,&name,&address,&currency_ID,&cur_id,&cur_name,&exchange_rate)
			if err != nil {
				fmt.Print(err)
			}
			//print results to console
			invoice_creation_date = invoice_creation_date[:10]
			response, err := http.Get("https://api.exchangeratesapi.io/"+invoice_creation_date)
			if err != nil {
				log.Fatal(err)
			}
			body, _ := ioutil.ReadAll(response.Body)
			
			var RatesToDate currencyRates

			err = json.Unmarshal(body,&RatesToDate)
			var amountInGBP float64
			var amountInGBPToDate float64
			if cur_name == "USD"{
				amountInGBPToDate = (amount/Rates.Rates.USD)*Rates.Rates.GBP
				amountInGBPToDate = math.Floor(amountInGBPToDate*100)/100
				amountInGBP = (amount/RatesToDate.Rates.USD)*RatesToDate.Rates.GBP
				amountInGBP = math.Floor(amountInGBP*100)/100
				fmt.Printf("Amount in GBP on Creation Date %v\n",amountInGBP)
				fmt.Printf("Amount in GBP today %v\n",amountInGBPToDate)
			}
			if cur_name == "EUR"{
				amountInGBPToDate = amount*Rates.Rates.GBP
				amountInGBPToDate = math.Floor(amountInGBPToDate*100)/100
				amountInGBP = amount*RatesToDate.Rates.GBP
				amountInGBP = math.Floor(amountInGBP*100)/100
				fmt.Printf("Amount in GBP on Creation Date %v\n",amountInGBP)
				fmt.Printf("Amount in GBP today %v\n",amountInGBPToDate)
			}
			if cur_name == "JPY"{
				amountInGBPToDate = (amount/Rates.Rates.JPY)*Rates.Rates.GBP
				amountInGBPToDate = math.Floor(amountInGBPToDate*100)/100
				amountInGBP = (amount/RatesToDate.Rates.JPY)*RatesToDate.Rates.GBP
				amountInGBP = math.Floor(amountInGBP*100)/100
				fmt.Printf("Amount in GBP on Creation Date %v\n",amountInGBP)
				fmt.Printf("Amount in GBP today %v\n",amountInGBPToDate)
			}
			fmt.Printf("Invoice ID:%v Date:%s Customer ID:%v Amount:%v\n",inv_ID,invoice_creation_date,Cus_ID,amount)
		}
	}
	
	if choice == "3"{
		rowss, err := db.Query("SELECT* FROM invoices JOIN Customers ON invoices.customer_id = Customers.id JOIN Currencies ON Customers.Currency_ID = CUrrencies.id ORDER BY currency_name;")
		if err != nil {
			fmt.Print(err)
		}
		defer rowss.Close()

		fmt.Println("---------------------------------------------------------------------")
		//loop through all employee results
		for rowss.Next(){
			//assign values to variables
			var inv_ID int
			var invoice_creation_date string
			var Cus_ID int
			var amount int
			var cus_id int
			var name string
			var address string
			var currency_ID int
			var cur_id int
			var cur_name string
			var exchange_rate string
			err := rowss.Scan(&inv_ID, &invoice_creation_date, &Cus_ID, &amount,&cus_id,&name,&address,&currency_ID,&cur_id,&cur_name,&exchange_rate)
			if err != nil {
				fmt.Print(err)
			}

			invoice_creation_date = invoice_creation_date[:10]
			response, err := http.Get("https://api.exchangeratesapi.io/"+invoice_creation_date)
			if err != nil {
				log.Fatal(err)
			}
			body, _ := ioutil.ReadAll(response.Body)
			
			var RatesToDate currencyRates

			err = json.Unmarshal(body,&RatesToDate)
			var amountInGBP float64
			var amountInGBPToDate float64
			if cur_name == "USD"{
				amountInGBPToDate = (amount/Rates.Rates.USD)*Rates.Rates.GBP
				amountInGBPToDate = math.Floor(amountInGBPToDate*100)/100
				amountInGBP = (amount/RatesToDate.Rates.USD)*RatesToDate.Rates.GBP
				amountInGBP = math.Floor(amountInGBP*100)/100
				fmt.Printf("Amount in GBP on Creation Date %v\n",amountInGBP)
				fmt.Printf("Amount in GBP today %v\n",amountInGBPToDate)
			}
			if cur_name == "EUR"{
				amountInGBPToDate = amount*Rates.Rates.GBP
				amountInGBPToDate = math.Floor(amountInGBPToDate*100)/100
				amountInGBP = amount*RatesToDate.Rates.GBP
				amountInGBP = math.Floor(amountInGBP*100)/100
				fmt.Printf("Amount in GBP on Creation Date %v\n",amountInGBP)
				fmt.Printf("Amount in GBP today %v\n",amountInGBPToDate)
			}
			if cur_name == "JPY"{
				amountInGBPToDate = (amount/Rates.Rates.JPY)*Rates.Rates.GBP
				amountInGBPToDate = math.Floor(amountInGBPToDate*100)/100
				amountInGBP = (amount/RatesToDate.Rates.JPY)*RatesToDate.Rates.GBP
				amountInGBP = math.Floor(amountInGBP*100)/100
				fmt.Printf("Amount in GBP on Creation Date %v\n",amountInGBP)
				fmt.Printf("Amount in GBP today %v\n",amountInGBPToDate)
			}
			//print results to console
			fmt.Printf("Invoice ID:%v Date:%s Customer ID:%v Amount:%v\n",inv_ID,invoice_creation_date,Cus_ID,amount)
		}
	}
	
}//end of main function
