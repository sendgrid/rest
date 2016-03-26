HTTP REST client, simplified for Go

Try out an example:

`go get github.com/sendgrid/rest`

Navigate to your go workspace + `/src/github.com/sendgrid/rest`

`echo "export SENDGRID_API_KEY='YOUR_API_KEY'" > sendgrid.env`

`echo "sendgrid.env" >> .gitignore`

`source ./sendgrid.env`

`go run example/example.go`
