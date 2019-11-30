module ImageCutter/cmd/cutter

go 1.12

require (
	ImageCutter/pkg/config v0.0.0
	ImageCutter/pkg/cropper v0.0.0
	ImageCutter/pkg/logger v0.0.0
	ImageCutter/pkg/lru v0.0.0
	ImageCutter/pkg/models v0.0.0
	ImageCutter/pkg/services/cutter v0.0.0
	github.com/disintegration/imaging v1.6.2 // indirect
)

replace (
	ImageCutter/pkg/config v0.0.0 => ../../pkg/config
	ImageCutter/pkg/cropper v0.0.0 => ../../pkg/cropper
	ImageCutter/pkg/logger v0.0.0 => ../../pkg/logger
	ImageCutter/pkg/lru v0.0.0 => ../../pkg/lru
	ImageCutter/pkg/models v0.0.0 => ../../pkg/models
	ImageCutter/pkg/services/cutter v0.0.0 => ../../pkg/services/cutter
)
