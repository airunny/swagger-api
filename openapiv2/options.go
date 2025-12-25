package openapiv2

import "github.com/go-kratos/grpc-gateway/v2/protoc-gen-openapiv2/generator"

type options struct {
	generatorOptions []generator.Option
	services         []string
}

type HandlerOption func(opt *options)

func WithGeneratorOptions(opts ...generator.Option) HandlerOption {
	return func(opt *options) {
		opt.generatorOptions = opts
	}
}

func WithServices(services ...string) HandlerOption {
	return func(opt *options) {
		opt.services = services
	}
}
