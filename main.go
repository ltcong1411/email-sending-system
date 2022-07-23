package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"
)

var (
	smtpHost       = os.Getenv("SMTP_EMAIL_HOST")
	smtpUsername   = os.Getenv("SMTP_EMAiL_USERNAME")
	smtpPassword   = os.Getenv("SMTP_EMAIL_PASSWORD")
	smtpPortNumber = os.Getenv("SMTP_EMAIL_PORT")

	apiSend  = os.Getenv("API_SEND")
	apiToken = os.Getenv("API_TOKEN")

	emailTemplateFile = os.Args[1]
	customerFile      = os.Args[2]
	outputEmailPath   = os.Args[3]
	errorFile         = os.Args[4]
)

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Customer struct which contains the information of customer
type Customer struct {
	TITLE      string
	FIRST_NAME string
	LAST_NAME  string
	EMAIL      string
}

// importCustumers function which import list of customer is stored in a CSV file
func importCustumers(file string) (customers []Customer) {
	csvFile, err := os.Open(file)
	if err != nil {
		log.Fatalln("Could not open file", err)
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Fatalln("Could not read CSV", err)
	}

	for i, line := range csvLines {
		// skip first row
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

// saveCustomerErrorToFile function which store customer that doesn’t have an email address or invalid email address
func saveCustomerErrorToFile(customerError Customer) (err error) {
	var data [][]string

	// check file exist
	f, err := os.OpenFile(errorFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil && os.IsExist(err) {
		log.Fatal(err)
		return
	}

	if os.IsNotExist(err) {
		// create file
		f, err = os.Create(errorFile)
		if err != nil {
			log.Fatalln("Could not create file", err)
		}

		data = [][]string{{"TITLE", "FIRST_NAME", "LAST_NAME", "EMAIL"}}
	}

	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	row := []string{customerError.TITLE, customerError.FIRST_NAME, customerError.LAST_NAME, customerError.EMAIL}
	data = append(data, row)

	err = w.WriteAll(data)
	if err != nil {
		log.Fatalln("Could not write to file", err)
		return
	}

	return
}

// Email struct which contains the information of that email
type Email struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	MineType string `json:"mineType"`
	Body     string `json:"body"`
}

// importEmailTemplate function which import the email template is stored in a JSON file
func importEmailTemplate(file string) (emailTemplate Email) {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'emailTemplate' which we defined above
	err = json.Unmarshal(byteValue, &emailTemplate)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// fillInfoToEmailTemplate function which filled customer information with email template
func fillInfoToEmailTemplate(customer Customer, emailTemplate Email) (emailInfo Email, err error) {
	// replace {{TODAY}} in the email template with the date on which it runs with the format is "31 Dec 2020"
	body := strings.ReplaceAll(emailTemplate.Body, `{{TODAY}}`, time.Now().Format("02 January 2006"))

	// replace `{{` in the email template to `{{.` to match the syntax of html/template library
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

func saveEmailInfoToFile(emailInfo Email) (err error) {
	// if the directory does not exist, it will be created first
	if _, err := os.Stat(outputEmailPath); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(outputEmailPath, 0755)
		} else {
			return err
		}
	}

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetIndent("", "\t")
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(emailInfo)
	// https://stackoverflow.com/questions/24656624/how-to-display-a-character-instead-of-ascii
	// https://developpaper.com/the-solution-of-escaping-special-html-characters-in-golang-json-marshal/

	fileName := fmt.Sprintf("%s%s.json", outputEmailPath, emailInfo.To)

	err = ioutil.WriteFile(fileName, bf.Bytes(), 0644)
	if err != nil {
		return
	}

	return
}

func sendEmailViaAPI(emailInfo Email) (err error) {
	fmt.Printf("send email via API - email address: %v - api: %v - token: %v\n", emailInfo.To, apiSend, apiToken)
	return
}

func sendEmailViaSMTP(emailInfo Email) (err error) {
	fmt.Printf("send email via SMTP - email address: %v - host: %v - username: %v - password: %v - port: %v\n", emailInfo.To, smtpHost, smtpUsername, smtpPassword, smtpPortNumber)
	return
}

func sendEmail(emailInfo Email) (err error) {
	sendEmailVia := "file"
	if len(os.Args) > 5 {
		sendEmailVia = os.Args[5]
	}

	switch sendEmailVia {
	case "api":
		err = sendEmailViaAPI(emailInfo)
	case "smtp":
		err = sendEmailViaSMTP(emailInfo)
	default:
		err = saveEmailInfoToFile(emailInfo)
	}

	return
}

func prepareAndSendEmail(customers <-chan Customer, emailTemplate Email, results chan<- error) {
	for customer := range customers {
		// check valid of this email address
		if customer.EMAIL == "" || !validEmail(customer.EMAIL) {
			err := saveCustomerErrorToFile(customer)
			if err != nil {
				fmt.Printf("saveCustomerErrorToFile failed - customer: %+v - err: %v\n", customer, err)
				results <- err
				return
			}

			results <- nil
			return
		}

		emailInfo, err := fillInfoToEmailTemplate(customer, emailTemplate)
		if err != nil {
			fmt.Printf("Could not fill Information To Email Template - email address: %v - err: %v\n", customer.EMAIL, err)
			results <- err
			return
		}

		err = sendEmail(emailInfo)
		if err != nil {
			fmt.Printf("Could not send email - email address: %v - err: %v\n", customer.EMAIL, err)
			results <- err
			return
		}

		results <- nil
	}

}

func main() {
	customers := importCustumers(customerFile)
	emailTemplate := importEmailTemplate(emailTemplateFile)

	customerChan := make(chan Customer, len(customers))
	result := make(chan error, len(customers))

	for w := 1; w <= 3; w++ {
		go prepareAndSendEmail(customerChan, emailTemplate, result)
	}

	for _, customer := range customers {
		customerChan <- customer
	}
	close(customerChan)

	for i := 0; i < len(customers); i++ {
		<-result
	}
}
