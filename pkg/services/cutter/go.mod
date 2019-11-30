module ImageCutter/pkg/services/cutter

go 1.12

require (
	github.com/gorilla/mux v1.7.3
	go.uber.org/zap v1.13.0
	ImageCutter/pkg/config v0.0.0
	ImageCutter/pkg/cropper v0.0.0
	ImageCutter/pkg/lru v0.0.0
	ImageCutter/pkg/models v0.0.0
)

replace (
	ImageCutter/pkg/config v0.0.0 => ../../../pkg/config
	ImageCutter/pkg/models v0.0.0 => ../../../pkg/models
	ImageCutter/pkg/models v0.0.0 => ../../../pkg/models
	ImageCutter/pkg/models v0.0.0 => ../../../pkg/models
	ImageCutter/pkg/models v0.0.0 => ../../../pkg/models
)
