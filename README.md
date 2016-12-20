[![Build Status](https://travis-ci.org/zensend/zensend_go_api.svg?branch=master)](https://travis-ci.org/zensend/zensend_go_api)
# ZenSend API for Go

## Installation

    go get github.com/zensend/zensend
    go install github.com/zensend/zensend

## Examples
First, make sure you import the zensend package
```go
import zensend
```
Then, create an instance of zensend.Client
```go
client := zensend.New("YOUR API KEY")
```

### Sending SMS
```go
message := &zensend.Message{
    Body: "The SMS Body",
    Originator: "FROM YOU",
    Numbers: []string{"447777777777", "448888888888"},
}

response, error := client.SendSMS(message)

if error == nil {
  log.Println("GUID:", response.TxGuid, "NUMBERS:", response.Numbers, "SMS Parts:", response.SmsParts, "Encoding:", response.Encoding)
} else {
  switch error := error.(type) {
  case zensend.ZenSendError:
    log.Println("Status Code:", error.StatusCode, "Failcode:", error.FailCode, "Parameter:", error.Parameter)
  default:
    log.Println(error)
  }
}
```

### Checking your balance
```go
balance, error = client.CheckBalance()

if error == nil {
  log.Println("Balance:", balance)
} else {
  switch error := error.(type) {
  case zensend.ZenSendError:
    log.Println("Status Code:", error.StatusCode, "Failcode:", error.FailCode, "Parameter:", error.Parameter)
  default:
    log.Println(error)
  }
}
```

### Retrieving pricing information
```go
prices, error = client.GetPrices()

if error == nil {
  for country := range prices {
    fmt.Printf("Country: %s. Price: %v\n", country, prices[country])
  }
} else {
  switch error := error.(type) {
  case zensend.ZenSendError:
    log.Println("Status Code:", error.StatusCode, "Failcode:", error.FailCode, "Parameter:", error.Parameter)
  default:
    log.Println(error)
  }
}
```

### Operator Lookup
```go
response, error := client.LookupOperator("441234567890")
if error == nil {
  fmt.Printf("MCC: %s MNC: %s Operator: %s CostInPence: %f NewBalanceInPence: %f", response.MCC, response.MNC, response.Operator, response.CostInPence, response.NewBalanceInPence)
} 

### Create Sub Account
```go
response, error := client.CreateSubAccount("sub account name")
if error == nil {
  fmt.Printf("Name: %s Api Key: %s", response.Name, response.ApiKey)
}
```
  
### Testing using the REPL

    gore
    :import zensend

