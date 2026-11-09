<!-- Thanks for the PR. Please fill in what you can — the "why" matters most. -->

## What does this change?

<!-- One or two sentences. Link the issue with "Fixes #N" if there is one. -->

## Why is this the right change?

## How was it checked?

```
gofmt -l .
go vet ./...
go test ./...
GOOS=windows GOARCH=amd64 go build .
```

## Notes for the reviewer

<!-- Did you keep platform code behind the build tags? Did any confirmation prompt change? -->
