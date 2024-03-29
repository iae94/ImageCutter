module github.com/iae94/ImageCutter

go 1.12

require (
	github.com/iae94/ImageCutter/pkg/logger v0.0.0 // indirect
	github.com/iae94/ImageCutter/pkg/models v0.0.0 // indirect
	github.com/DATA-DOG/godog v0.7.13
)
replace (
	ImageCutter/pkg/config v0.0.0 => /pkg/config
	ImageCutter/pkg/cropper v0.0.0 => /pkg/cropper
	ImageCutter/pkg/cutter v0.0.0 => /pkg/cutter
	github.com/iae94/ImageCutter/pkg/logger v0.0.0 => /pkg/logger
	ImageCutter/pkg/models v0.0.0 => /pkg/models
)
