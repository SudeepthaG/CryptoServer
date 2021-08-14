# CryptoServer

# Execution Steps

1. Download the git folder
2. Open terminal
3. cd into the git folder
4. Use command: go run main.go structures.go
5. Open localhost:8085/currency/all or localhost:8085/currency/{symbol} in terminal

## External Libraries used and their purpose

- "encoding/json"  : To decode external json file. I have used an external json file as the coding question requires the input symbols to be configurable. Any new symbols can be added to the json file.
- "fmt" : To print any errors
- "io" : To read and write data from streams
- "log" : To log any errors
- "net/http" :  to make an http connection, send requests and receive responses
- "os" : To open the external files
- "strings" :  to manipulate strings
- "sync" :  to provide mutex while using goroutines and channels to avoid race condition
- "golang.org/x/net/websocket" :  to start a socket to maintain a continuous connection to avoid wastage of time during sending requests and responses
