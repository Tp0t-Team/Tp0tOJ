//go:build !WithFrontEnd

package services

import "embed"

var staticFolder embed.FS

const HasFrontEnd = false
