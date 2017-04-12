package bosh

import (
	"encoding/base64"
	"fmt"
	"time"
)

func basicAuth(user, pass string) string {
	auth := user + ":" + pass
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func lapse(a, b time.Time) string {
	d := b.Sub(a)
	if f := d.Hours(); f >= 1.5 {
		return fmt.Sprintf("%dh", int(f))
	}
	if f := d.Minutes(); f >= 1.5 {
		return fmt.Sprintf("%dm", int(f))
	}
	if f := d.Seconds(); f >= 1.5 {
		return fmt.Sprintf("%ds", int(f))
	}
	return ""
}
