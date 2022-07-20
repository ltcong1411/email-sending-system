package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Customer struct {
	TITLE      string
	FIRST_NAME string
	LAST_NAME  string
	EMAIL      string
}

func getCustumers(customerPath string) (customers []Customer, err error) {
	csvFile, err := os.Open(customerPath)
	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for i, line := range csvLines {
		if i == 0 {
			continue
		}

		customer := Customer{
			TITLE:      line[0],
			FIRST_NAME: line[1],
			LAST_NAME:  line[2],
			EMAIL:      line[3],
		}
		customers = append(customers, customer)
	}

	return
}

type EmailTemplate struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	MineType string `json:"mineType"`
	Body     string `json:"body"`
}

func getEmailTemplate(emailTemplatePath string) (emailTemplate EmailTemplate, err error) {
	// Open our jsonFile
	jsonFile, err := os.Open(emailTemplatePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &emailTemplate)

	return
}

func fillInfo(customer Customer, emailTemplate EmailTemplate) (emailInfo EmailTemplate, err error) {
	body := strings.ReplaceAll(emailTemplate.Body, `{{TODAY}}`, time.Now().Format("02 January 2006"))
	body = strings.ReplaceAll(body, `{{`, `{{.`)

	tmp := template.New("simple")
	tmp, err = tmp.Parse(body)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	err = tmp.Execute(&b, &customer)
	if err != nil {
		log.Fatal(err)
	}

	emailInfo = emailTemplate
	emailInfo.To = customer.EMAIL
	emailInfo.Body = b.String()

	return
}

func sendEmail(emailInfo EmailTemplate) (err error) {
	fmt.Printf("emailInfo: %+v\n", emailInfo)

	return
}

func main() {
	customers, _ := getCustumers("customers.csv")
	fmt.Printf("customers: %+v\n", customers)

	emailTemplate, _ := getEmailTemplate("email_template.json")
	fmt.Printf("emailTemplate: %+v\n", emailTemplate)

	for _, customer := range customers {
		emailInfo, _ := fillInfo(customer, emailTemplate)
		sendEmail(emailInfo)
	}

}
