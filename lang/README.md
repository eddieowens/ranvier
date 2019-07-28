# Ranvier lang
The "compiler" for converting the [config schema files](https://github.com/eddieowens/ranvier/wiki/Config-schema-files)
into config that an application can consume.

## Installation
```bash
go get github.com/eddieowens/ranvier/lang
```

## Usage
```go
package main

import (
  "fmt"
  "github.com/eddieowens/ranvier/lang"
  "github.com/eddieowens/ranvier/lang/compiler"
)

func main() {
  compiler := lang.NewCompiler()
  schema, err := compiler.Compile("path/to/config.json", compiler.CompileOptions{
  	ParseOptions: compiler.ParseOptions{
  		// The base directory housing the config.json file
  		Root: "/",
  	},
  })
  if err != nil {
    panic(err)
  }
  
  // Prints path-to-config
  fmt.Println(schema.Name)
}
```
The above will compile the `path/to/config.json` file and all of its dependencies into a single JSON file.