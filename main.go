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

	emailTemplateFile = "email_template.json"
	customerFile      = "customers.csv"
	outputEmailPath   = "output_emails/"
	errorFile         = "errors.csv"
	sendEmailVia      = "file"
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
func saveCustomerErrorToFile(customerError Customer, file string) (err error) {
	var data [][]string

	// check file exist
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil && os.IsExist(err) {
		return
	}

	if os.IsNotExist(err) {
		// create file
		f, err = os.Create(file)
		if err != nil {
			return
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
func importEmailTemplate(file string) (emailTemplate *Email) {
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

	tmp := template.New("email-sending-system")
	tmp, err = tmp.Parse(body)
	if err != nil {
		return
	}

	var b bytes.Buffer
	err = tmp.Execute(&b, &customer)
	if err != nil {
		return
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

func prepareAndSendEmail(customers <-chan Customer, emailTemplate Email, results chan<- error) {
	for customer := range customers {
		// check valid of this email address
		if customer.EMAIL == "" || !validEmail(customer.EMAIL) {
			err := saveCustomerErrorToFile(customer, errorFile)
			if err != nil {
				results <- fmt.Errorf("saveCustomerErrorToFile failed - customer: %+v - err: %v", customer, err)
				return
			}

			fmt.Printf("can't send email to %v-%v-%v\n", customer.TITLE, customer.FIRST_NAME, customer.LAST_NAME)

			results <- nil
			return
		}

		emailInfo, err := fillInfoToEmailTemplate(customer, emailTemplate)
		if err != nil {
			results <- fmt.Errorf("could not fill information to email template - email address: %v - err: %v", customer.EMAIL, err)
			return
		}

		switch sendEmailVia {
		case "api":
			err = sendEmailViaAPI(emailInfo)
		case "smtp":
			err = sendEmailViaSMTP(emailInfo)
		default:
			err = saveEmailInfoToFile(emailInfo)
		}

		if err != nil {
			results <- fmt.Errorf("could not send email - email address: %v - err: %v", customer.EMAIL, err)
			return
		}

		fmt.Printf("sent an email to: %v-%v-%v with email address: %v\n", customer.TITLE, customer.FIRST_NAME, customer.LAST_NAME, customer.EMAIL)

		results <- nil
	}

}

func sendEmail(customers []Customer, emailTemplate Email) {
	customerChan := make(chan Customer, len(customers))
	results := make(chan error, len(customers))

	// start up 3 workers to run prepareAndSendEmail function
	for w := 1; w <= 3; w++ {
		go prepareAndSendEmail(customerChan, emailTemplate, results)
	}

	// send customer information to customer channel
	for _, customer := range customers {
		customerChan <- customer
	}
	close(customerChan)

	for i := 0; i < len(customers); i++ {
		err := <-results
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	// can use "flag" to manage the command-line arguments
	if len(os.Args) >= 5 {
		emailTemplateFile = os.Args[1]
		customerFile = os.Args[2]
		outputEmailPath = os.Args[3]
		errorFile = os.Args[4]

		// check the conditions for sending email via
		if len(os.Args) >= 6 {
			sendEmailVia = os.Args[5]
		}
	} else {
		log.Fatal("You need to enter more information about the path to the email template file, the customer list, the email output folder, the error customer list.\nFor example: go run main.go email_template.json customers.csv output_emails/ errors.csv")
	}

	customers := importCustumers(customerFile)
	if len(customers) == 0 {
		fmt.Println("empty customer list, please check again")
		return
	}

	emailTemplate := importEmailTemplate(emailTemplateFile)
	if emailTemplate == nil {
		fmt.Println("no email template, please check again")
		return
	}

	sendEmail(customers, *emailTemplate)
}
