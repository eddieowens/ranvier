# Ranvier Go client
This package houses all of the code necessary for communicating with the Ranvier [server](../server).

## Installation
```bash
go get github.com/eddieowens/ranvier/client
```

## Usage
```go
package main

import (
  "fmt"
  "github.com/eddieowens/ranvier/client"
  "os"
)

func main() {
  c, err := client.NewClient(&client.ClientOptions{
    Hostname:        "localhost:8080",
    ConfigDirectory: os.TempDir(),
  })
  
  if err != nil {
  	panic(err)
  }
  
  // Listens to the Ranvier server for changes to the 'users' config file. If someone were to make an update to this 
  // config file, the client will receive a message, and update its internal cache for the next query. 
  _, err = c.Connect(&client.ConnOptions{
  	Names: []string{"users"},
  })
  
  if err != nil {
  	panic(err)
  }
  
  // Query Ranvier using a jsonpath query to retrieve the config.
  conf, err := c.Query(&client.QueryOptions{
  	Name: "users",
  	Query: "$.db",
  })
  
  if err != nil {
  	panic(err)
  }
  
  // Prints 'users'
  fmt.Println(conf.Name)
}
```