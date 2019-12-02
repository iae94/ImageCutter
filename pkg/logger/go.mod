module github.com/iae94/ImageCutter/pkg/logger

go 1.12

require (
	ImageCutter/pkg/config v0.0.0
	go.uber.org/zap v1.13.0
)

replace ImageCutter/pkg/config v0.0.0 => ../../pkg/config
