//go:build tools
// +build tools

package tools

//go:generate windres build/windows/app.rc -O coff -o cmd/app.syso
