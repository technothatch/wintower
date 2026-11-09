# Contributing to wintower

Thanks for helping improve a small Windows TUI.

## Ground rules

- Search [existing issues](https://github.com/technothatch/wintower/issues) before opening a new one.
- One problem per issue, please.
- Be kind. This is maintained in spare time; assume good intent.

## Reporting a bug

Use the **Bug report** template. What matters most:

1. `wintower --version`
2. Your Windows version and architecture (x64 / ARM64)
3. What you pressed and what happened instead — a screenshot of the terminal
   (or the text of it) is ideal.

If a service or process action fails, say whether you were elevated: Windows
refuses those calls without admin rights, and that distinction determines
whether it's a bug here or expected behaviour.

## Suggesting a feature

Use the **Feature request** template. Say which tab it belongs on and what
you're trying to find out or change. The bar for a new key is that it must not
break the "every mutating action asks for confirmation first" promise.

## Sending a pull request

1. Fork and branch from `main`.
2. Keep platform code behind the existing build tags
   (`//go:build windows` in `internal/procs`, `internal/services`,
   `internal/evtlog`) so the package still builds elsewhere.
3. Run the checks locally:

   ```
   gofmt -l .
   go vet ./...
   go test ./...
   GOOS=windows GOARCH=amd64 go build .
   ```

4. Open a PR that references the issue it fixes and explains *why* the change
   is correct, not just what it does.

Small, focused PRs merge fastest.

## Releases

Maintainers cut releases and attach the prebuilt `wintower-amd64.exe` /
`wintower-arm64.exe` binaries. Contributors don't need to worry about versioning.
