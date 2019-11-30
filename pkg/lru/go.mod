module ImageCutter/pkg/lru

go 1.12

require (
	ImageCutter/pkg/config v0.0.0
	ImageCutter/pkg/logger v0.0.0
	ImageCutter/pkg/models v0.0.0
	go.uber.org/zap v1.13.0
)

replace ImageCutter/pkg/models v0.0.0 => ../../pkg/models

replace ImageCutter/pkg/config v0.0.0 => ../../pkg/config

replace ImageCutter/pkg/logger v0.0.0 => ../../pkg/logger
