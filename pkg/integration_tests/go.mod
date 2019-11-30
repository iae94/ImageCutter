module ImageCutter/pkg/integration_tests

go 1.12

require (
	ImageCutter/pkg/config v0.0.0
	ImageCutter/pkg/cropper v0.0.0
	ImageCutter/pkg/logger v0.0.0
	ImageCutter/pkg/lru v0.0.0
	ImageCutter/pkg/models v0.0.0
	ImageCutter/pkg/services/cutter v0.0.0
	github.com/DATA-DOG/godog v0.7.13
	github.com/disintegration/imaging v1.6.2 // indirect
	go.uber.org/zap v1.13.0
)

replace (
	ImageCutter/pkg/config v0.0.0 => ../../pkg/config
	ImageCutter/pkg/cropper v0.0.0 => ../../pkg/cropper
	ImageCutter/pkg/logger v0.0.0 => ../../pkg/logger
	ImageCutter/pkg/lru v0.0.0 => ../../pkg/lru
	ImageCutter/pkg/models v0.0.0 => ../../pkg/models
	ImageCutter/pkg/services/cutter v0.0.0 => ../../pkg/services/cutter
)
