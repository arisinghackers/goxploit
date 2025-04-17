# goxploit

**GoXploit** is a modular Go client for interacting with the [Metasploit RPC API](https://docs.rapid7.com/metasploit/standard-api-methods-reference).  
It includes a code generator that scrapes official documentation and produces strongly-typed API wrappers.

## Features

- RPC client for Metasploit
- Automatic method generation from Rapid7's API docs
- Easy integration with other Go projects
- Simple structure for extending or contributing

## Installation

```bash
go get github.com/arisinghackers/goxploit
```

## Usage 
Authenticate and send a request
```go
client := msfrpc.NewMsfRpcClient("your_password", "false", "your_username", "127.0.0.1", 55552, "/api")
token, err := client.MsfAuth()
if err != nil {
    log.Fatalf("Auth error: %v", err)
}
fmt.Println("Token:", token)

response, err := client.AuthenticatedRequest([]any{"core.version"})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Core Version: %+v\n", response)
```
##### See /examples folder

## Contributing
Contributions are welcome. Please keep changes modular and follow idiomatic Go practices.
