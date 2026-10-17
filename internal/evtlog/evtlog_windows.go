//go:build windows

package evtlog

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"unicode/utf16"
)

const psScript = `$ErrorActionPreference='Stop'
Get-WinEvent -FilterHashtable @{LogName='System'; Level=1,2,3} -MaxEvents 100 |
  Select-Object @{n='Time';e={$_.TimeCreated.ToString('o')}}, Id, LevelDisplayName, ProviderName, Message |
  ConvertTo-Json -Compress`

func encodeCommand(script string) string {
	utf16le := utf16.Encode([]rune(script))
	buf := make([]byte, len(utf16le)*2)
	for i, v := range utf16le {
		buf[i*2] = byte(v)
		buf[i*2+1] = byte(v >> 8)
	}
	return base64.StdEncoding.EncodeToString(buf)
}

func Events() ([]Event, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-EncodedCommand", encodeCommand(psScript))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("read event log: %w", err)
	}
	if len(out) == 0 {
		return nil, nil
	}
	return ParseJSON(out)
}
