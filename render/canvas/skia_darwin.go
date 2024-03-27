//go:build skia && darwin

package canvas

/*
#cgo LDFLAGS: -framework Foundation
#cgo LDFLAGS: -framework CoreFoundation
#cgo LDFLAGS: -framework CoreGraphics
#cgo LDFLAGS: -framework CoreText
#cgo LDFLAGS: -framework ImageIO
#cgo LDFLAGS: -framework Carbon
#cgo LDFLAGS: -framework Metal
#
#cgo LDFLAGS: -mmacosx-version-min=14.0
*/
import "C"
