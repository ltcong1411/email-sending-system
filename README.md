# Email-Sending-CLI
Console application to send emails using a template

## Usage
* Create customer list stored in CSV file (**customer.csv**), which has the following format:
```
TITLE,FIRST_NAME,LAST_NAME,EMAIL
Mr,Cong,Le,ltcong1411@gmail.com
Mrs,Thu,Vuong,vtathu32@gmail.com
Mr,Danh,Le,ltdanh0805@gmail.com
Mr,A,B,xyz
Mrs,X,Y,
```

* Create email template stored in JSON file (**email_template.json**), which has the following format:

```json
{
    "from": "The Marketing Team<marketing@example.com",
    "subject": "A new product is being launched soon...",
    "mineType": "text/plain",
    "body": "Hi {{TITLE}} {{FIRST_NAME}} {{LAST_NAME}},\nToday, {{TODAY}}, we would like to tell you that... Sincerely,\nThe Marketing Team"
}
```

* Run application
```
go run main.go /path/to/email_template.json /path/to/customers.csv /path/to/output_emails/ /path/to/errors.csv
```

* Run application with the built file
```
go build
email-sending-system /path/to/email_template.json /path/to/customers.csv /path/to/output_emails/ /path/to/errors.csv
```

* Run application with docker
```
docker build -t email-sending-system .
docker run email-sending-system /path/to/email_template.json /path/to/customers.csv /path/to/output_emails/ /path/to/errors.csv
```