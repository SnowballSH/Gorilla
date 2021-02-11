powershell ./withenv GOOS=darwin GOARCH=amd64 go build -o gorilla-darwin-amd64
powershell ./withenv GOOS=windows GOARCH=amd64 go build -o gorilla-windows-amd64.exe
powershell ./withenv GOOS=linux GOARCH=amd64 go build -o gorilla-linux-amd64