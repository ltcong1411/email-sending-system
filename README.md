# Email-Sending-System
Console application to send emails using a template

## Usage
* Run test application
```
go test --cover
```

* Create customer list stored in CSV file (**customers.csv**), which has the following format:
```
TITLE,FIRST_NAME,LAST_NAME,EMAIL
Mr,John,Smith,john.smith@example.com
Mrs,Michelle,Smith,michelle.smith@example.com
Mr,Cong,Le,ltcong1411@gmail.com
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
***/path/to/email_template.json*** : the path to email template file

***/path/to/customers.csv*** : the path to customer list file

***/path/to/output_emails/*** : the path to folder where emails are saved

***/path/to/errors.csv*** : the path to the file that stores the customer information has the wrong or doesn't have email address


* Run application with the built file
```
go build
./email-sending-system /path/to/email_template.json /path/to/customers.csv /path/to/output_emails/ /path/to/errors.csv
```

* Run application with docker
```
docker build -t email-sending-system .
docker run -v path/to/data/directory:/app/data email-sending-system data/email_template.json data/customers.csv data/output_emails/ data/errors.csv
```
***path/to/data/directory*** : the path to the directory containing the files **customers.csv** and **email_template.json**. After running the application the **output_emails** folder and the **errors.csv** file will also be created here



