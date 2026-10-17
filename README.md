# wintower

**A terminal cockpit for Windows: services, processes and the System event log in one keyboard-driven TUI.**

Windows has `services.msc`, `taskmgr.exe` and `eventvwr.msc` — three heavy GUIs for three everyday jobs.
`wintower` puts all three in your terminal, keyboard-first, in a single 2 MB binary with no runtime dependencies.

```
wintower v0.1.0  Services  Processes  Event Log  About

NAME                          DISPLAY NAME                     STATUS      START
Windows Update                Windows Update                   Running     Automatic
Spooler                       Print Spooler                    Running     Automatic
Some Legacy Service           Some Legacy Service              Stopped     Disabled
...

start Spooler? press y to confirm
tab/h/l switch · 1-4 jump · r refresh · s start · t restart · x stop/kill · y confirm · q quit
```

## Install

**Prebuilt binaries:** download `wintower-amd64.exe` (or `arm64`) from [Releases](../../releases), drop it on your `PATH`.

**From source** (Go 1.27+):

```
go install github.com/technothatch/wintower@latest
```

## Keys

| Key | Action |
|-----|--------|
| `tab` / `h` / `l` | switch tab |
| `1` `2` `3` `4` | jump to tab |
| `r` | refresh current tab |
| `s` | start selected service |
| `t` | restart selected service |
| `x` | stop selected service / kill selected process |
| `y` | confirm the pending action |
| `q` / `ctrl+c` | quit |

Every mutating action asks for confirmation first. wintower never changes anything on its own.

## What's under the hood (no magic)

Honest technical notes, because tools should say what they actually do:

- **Services tab** talks to the Windows Service Control Manager through the native Win32 API
  (`golang.org/x/sys/windows/svc/mgr`). No `sc.exe` shelling out.
- **Processes tab** enumerates processes with `EnumProcesses` and terminates with `TerminateProcess`.
  Access-denied processes (system, other users) are listed but cannot be killed without elevation — that's Windows, not a bug.
- **Event Log tab** tails the newest 100 Critical/Error/Warning entries from the `System` log via
  `Get-WinEvent` (the one place a PowerShell helper is the sane implementation).
- On non-Windows systems the binary still builds and runs; the data tabs explain that they need Windows.
  CI runs the full test suite on both Linux and Windows runners.

## Development

```
go test ./...
go build .
```

Service/process smoke tests only run on Windows; the event-log JSON parser is tested everywhere.

## License

[MIT](LICENSE)
