//go:build WithFrontEnd

package services

import "embed"

//go:embed static
var staticFolder embed.FS

const HasFrontEnd = true
