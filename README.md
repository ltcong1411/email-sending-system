# Email-Sending-CLI



## Usage
* Create customer list stored in CSV file (**customer.csv**), which has the following format:
```
TITLE,FIRST_NAME,LAST_NAME,EMAIL
Mr,Cong,Le,ltcong1411@gmail.com
Mrs,Thu,Vuong,vtathu32@gmail.com
Mr,Danh,Le,ltdanh0805@gmail.com
Mr,A,B,xyz
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

* Run
```
go run main.go email_template.json customers.csv output_emails/ errors.csv
```
