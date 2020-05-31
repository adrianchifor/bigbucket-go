# bigbucket-go

Golang client for the [Bigbucket](https://github.com/adrianchifor/Bigbucket) database.

## Install

```
go get -u github.com/adrianchifor/bigbucket-go
```

## Usage

```go
import (
  "fmt"
  "log"

  "github.com/adrianchifor/bigbucket-go"
)

client := bigbucket.NewClient()

// --------------------------
//       Show tables
// --------------------------
tables, err := client.GetTables()
if err != nil {
  log.Fatalln(err.Error())
}
fmt.Println(tables)
// [table1, table2]

// --------------------------
//        Use table
// --------------------------
tbl := client.UseTable("table1")

// --------------------------
//       List columns
// --------------------------
columns, err := tbl.ListColumns()
if err != nil {
  log.Fatalln(err.Error())
}
fmt.Println(columns)
// [col1 col2 col3 col4]

// --------------------------
//       Count rows
// --------------------------
rows, err := tbl.CountRows()
if err != nil {
  log.Fatalln(err.Error())
}
fmt.Println(rows)
// 10

// --------------------------
//      List row keys
// --------------------------
rowKeys, err := tbl.ListRows()
if err != nil {
  log.Fatalln(err.Error())
}
fmt.Println(rowKeys)
// [key0 key1 key2 key3 key4 key5 key6 key7 key8 key9]

// --------------------------
//     Read single row
// --------------------------
row, err := tbl.ReadRow("key1", bigbucket.WithColumns([]string{"col1", "col2"}))
if err != nil {
  log.Fatalln(err.Error())
}
fmt.Println(row)
// map[col1:value col2:value]

// --------------------------
//      Read all rows
// --------------------------
allRows, err := tbl.ReadRows()
if err != nil {
  log.Fatalln(err.Error())
}
fmt.Println(allRows)
// map[
//   key0: map[col1:value col2:value col3:value col4:value]
//   key1: map[col1:value col2:value col3:value col4:value]
//   key2: map[col1:value col2:value col3:value col4:value]
//   ...
// ]

// --------------------------
//  Read rows with prefix
// --------------------------
prefixRows, err := tbl.ReadRows(bigbucket.WithPrefix("key"),
  bigbucket.WithColumns([]string{"col3"}),
  bigbucket.WithLimit(2))
if err != nil {
  log.Fatalln(err.Error())
}
fmt.Println(prefixRows)
// map[
//   key0: map[col3:value]
//   key1: map[col3:value]
// ]

// --------------------------
//         Set row
// --------------------------
updates := map[string]string{
  "col1": "new value",
  "col4": "new value",
}
err = tbl.SetRow("key1", updates)
if err != nil {
  log.Fatalln(err.Error())
}

updatedRow, _ := tbl.ReadRow("key1")
fmt.Println(updatedRow)
// map[col1:new value col2:value col3:value col4:new value]

// --------------------------
//    Delete single row
// --------------------------
err = tbl.DeleteRow("key1")
if err != nil {
  log.Fatalln(err.Error())
}

// --------------------------
//  Delete rows with prefix
// --------------------------
err = tbl.DeleteRows("key")
if err != nil {
  log.Fatalln(err.Error())
}

// --------------------------
//  Delete column (async)
// --------------------------
err = tbl.DeleteColumn("col1")
if err != nil {
  log.Fatalln(err.Error())
}

// --------------------------
//  Delete table (async)
// --------------------------
err = tbl.DeleteTable()
if err != nil {
  log.Fatalln(err.Error())
}
```

## Client configuration

Defaults

```go
client := bigbucket.NewClient(bigbucket.WithAddress("http://localhost:8080"),
  bigbucket.WithTimeout(30),
  bigbucket.WithGcpAuth(false),
  bigbucket.WithRequestHeaders(map[string]string{}))
```

```
WithAddress - The Bigbucket API endpoint
              e.g. port 80      > http://endpoint
                   port 443 TLS > https://endpoint
                   custom port  > http://endpoint:port

WithTimeout - Requests timeout (seconds, int)

WithGcpAuth - If true it generates GCP JWT tokens for address and adds "Authorization" header to requests.
              Useful when running clients in GCP and need to authenticate against private Cloud Run services.
              (docs -> https://cloud.google.com/run/docs/authenticating/service-to-service)

WithRequestHeaders - Custom headers to add to requests
```

## Running Bigbucket

See https://github.com/adrianchifor/Bigbucket#running