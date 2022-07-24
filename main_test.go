package main

import (
	"strings"
	"testing"
	"time"
)

type testValidEmail struct {
	name     string
	email    string
	expected bool
}

func TestValidEmail(t *testing.T) {
	tests := []testValidEmail{
		{
			name:     "Invalid email address",
			email:    "abc",
			expected: false,
		},
		{
			name:     "Valid email address",
			email:    "abc@gmail.com",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := validEmail(tc.email)
			if actual != tc.expected {
				t.Errorf("expected '%v', got '%v'", tc.expected, actual)
			}
		})
	}
}

type testImportCustumers struct {
	name                   string
	file                   string
	lenOfCustomersExpected int
}

func TestImportCustumers(t *testing.T) {
	tests := []testImportCustumers{
		{
			name:                   "There are 4 customers on the list",
			file:                   "./testdata/customers.csv",
			lenOfCustomersExpected: 4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			actual := importCustumers(tc.file)
			if len(actual) != tc.lenOfCustomersExpected {
				t.Errorf("expected '%v', got '%v'", tc.lenOfCustomersExpected, len(actual))
			}
		})
	}
}

type testSaveCustomerErrorToFile struct {
	name     string
	customer Customer
	file     string
	expected error
}

func TestSaveCustomerErrorToFile(t *testing.T) {
	tests := []testSaveCustomerErrorToFile{
		{
			name: "File does not exist",
			customer: Customer{
				TITLE:      "Mr",
				FIRST_NAME: "A",
				LAST_NAME:  "B",
				EMAIL:      "xyz",
			},
			file:     "testdata/errors_1.csv",
			expected: nil,
		},
		{
			name: "File already exists",
			customer: Customer{
				TITLE:      "Mr",
				FIRST_NAME: "A",
				LAST_NAME:  "B",
				EMAIL:      "xyz",
			},
			file:     "testdata/errors.csv",
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := saveCustomerErrorToFile(tc.customer, tc.file)
			if actual != tc.expected {
				t.Errorf("expected '%v', got '%v'", tc.expected, actual)
			}
		})
	}
}

func TestImportEmailTemplate(t *testing.T) {
	emailTemplate := Email{
		From:     "The Marketing Team<marketing@example.com",
		Subject:  "A new product is being launched soon...",
		MineType: "text/plain",
		Body:     "Hi {{TITLE}} {{FIRST_NAME}} {{LAST_NAME}},\nToday, {{TODAY}}, we would like to tell you that... Sincerely,\nThe Marketing Team",
	}

	actual := importEmailTemplate("./testdata/email_template.json")
	if actual == nil {
		t.Errorf("expected '%v', got '%v'", emailTemplate, actual)
	}

}

type testFillInfoToEmailTemplate struct {
	name          string
	customer      Customer
	emailTemplate Email
	expected      Email
}

func TestFillInfoToEmailTemplate(t *testing.T) {
	tests := []testFillInfoToEmailTemplate{
		{
			name: "Fill customer information to email template",
			customer: Customer{
				TITLE:      "Mr",
				FIRST_NAME: "Cong",
				LAST_NAME:  "Le",
				EMAIL:      "ltcong1411@gmail.com",
			},
			emailTemplate: Email{
				From:     "The Marketing Team<marketing@example.com",
				Subject:  "A new product is being launched soon...",
				MineType: "text/plain",
				Body:     "Hi {{TITLE}} {{FIRST_NAME}} {{LAST_NAME}},\nToday, {{TODAY}}, we would like to tell you that... Sincerely,\nThe Marketing Team",
			},
			expected: Email{
				From:     "The Marketing Team<marketing@example.com",
				To:       "ltcong1411@gmail.com",
				Subject:  "A new product is being launched soon...",
				MineType: "text/plain",
				Body:     "Hi Mr Cong Le,\nToday, {{TODAY}}, we would like to tell you that... Sincerely,\nThe Marketing Team",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expected.Body = strings.ReplaceAll(tc.expected.Body, `{{TODAY}}`, time.Now().Format("02 January 2006"))
			actual, _ := fillInfoToEmailTemplate(tc.customer, tc.emailTemplate)
			if actual.Body != tc.expected.Body {
				t.Errorf("expected '%v', got '%v'", tc.expected, actual)
			}
		})
	}
}

type testSaveEmailInfoToFile struct {
	name     string
	email    Email
	expected error
}

func TestSaveEmailInfoToFile(t *testing.T) {
	outputEmailPath = "testdata/"
	tests := []testSaveEmailInfoToFile{
		{
			name: "Save email information to file",
			email: Email{
				From:     "The Marketing Team<marketing@example.com",
				To:       "ltcong1411@gmail.com",
				Subject:  "A new product is being launched soon...",
				MineType: "text/plain",
				Body:     "Hi Mr Cong Le,\nToday, 24 July 2022, we would like to tell you that... Sincerely,\nThe Marketing Team",
			},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := saveEmailInfoToFile(tc.email)
			if actual != tc.expected {
				t.Errorf("expected '%v', got '%v'", tc.expected, actual)
			}
		})
	}
}

type testSendEmail struct {
	name          string
	customers     []Customer
	emailTemplate Email
	expected      error
}

func TestSendEmail(t *testing.T) {
	tests := []testSendEmail{
		{
			name: "send email",
			customers: []Customer{
				{
					TITLE:      "Mr",
					FIRST_NAME: "Cong",
					LAST_NAME:  "Le",
					EMAIL:      "ltcong1411@gmail.com",
				},
				{
					TITLE:      "Mr",
					FIRST_NAME: "A",
					LAST_NAME:  "B",
					EMAIL:      "xyz",
				},
			},
			emailTemplate: Email{
				From:     "The Marketing Team<marketing@example.com",
				Subject:  "A new product is being launched soon...",
				MineType: "text/plain",
				Body:     "Hi {{TITLE}} {{FIRST_NAME}} {{LAST_NAME}},\nToday, {{TODAY}}, we would like to tell you that... Sincerely,\nThe Marketing Team",
			},
			expected: nil,
		},
	}

	emailTemplateFile = "testdata/email_template.json"
	customerFile = "testdata/customers.csv"
	outputEmailPath = "testdata/output_emails/"
	errorFile = "testdata/errors.csv"
	sendEmailVia = "file"

	// test save email to file
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sendEmail(tc.customers, tc.emailTemplate)
		})
	}

	// test save email via api
	sendEmailVia = "api"
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sendEmail(tc.customers, tc.emailTemplate)
		})
	}

	// test save email via smtp
	sendEmailVia = "smtp"
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sendEmail(tc.customers, tc.emailTemplate)
		})
	}
}
