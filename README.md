# Simple SSH Checker
A simple ssh Checker to verify If my students have a good enough ssh client

## Requirements
- go 1.18+

## Getting started

run 
```
go run . -p 4222
``` 

(default port on 2222 )


If you need to store data, just run then:
```
ssh -p 4222 localhost WITH SOME DATA
```

If you want to see the data, run:
```
ssh -p 4222 decima@localhost
```
