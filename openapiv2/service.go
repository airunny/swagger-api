package openapiv2

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/grpc-gateway/v2/protoc-gen-openapiv2/generator"
	"github.com/go-kratos/kratos/v2/api/metadata"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/pluginpb"
)

// Service is service
type Service struct {
	ser  *metadata.Server
	opts *options
}

// New service
func New(srv *grpc.Server, handlerOpts ...HandlerOption) *Service {
	opts := &options{
		// Compatible with default UseJSONNamesForFields is true
		generatorOptions: []generator.Option{generator.UseJSONNamesForFields(true)},
	}

	for _, o := range handlerOpts {
		o(opts)
	}
	return &Service{
		ser:  metadata.NewServer(srv),
		opts: opts,
	}
}

func (s *Service) hasServices(service string) bool {
	for _, serviceName := range s.opts.services {
		if strings.EqualFold(serviceName, service) {
			return true
		}
	}
	return false
}

// ListServices list services
func (s *Service) ListServices(ctx context.Context, in *metadata.ListServicesRequest) (*metadata.ListServicesReply, error) {
	rsp, err := s.ser.ListServices(ctx, &metadata.ListServicesRequest{})
	if err != nil {
		return nil, err
	}

	if len(s.opts.services) <= 0 {
		return rsp, nil
	}

	out := &metadata.ListServicesReply{
		Services: make([]string, 0, len(rsp.Services)),
		Methods:  make([]string, 0, len(rsp.Methods)),
	}

	for _, service := range rsp.Services {
		if !s.hasServices(service) {
			continue
		}
		out.Services = append(out.Services, service)
	}

	for _, method := range rsp.Methods {
		filter := true
		for _, service := range s.opts.services {
			if strings.HasPrefix(method, "/"+service) {
				filter = false
				break
			}
		}

		if filter {
			continue
		}
		out.Methods = append(out.Methods, method)
	}
	return out, nil
}

// GetServiceOpenAPI get service open api
func (s *Service) GetServiceOpenAPI(ctx context.Context, in *metadata.GetServiceDescRequest, onlyRPC bool) (string, error) {
	protoSet, err := s.ser.GetServiceDesc(ctx, in)
	if err != nil {
		return "", err
	}
	files := protoSet.FileDescSet.File
	var target string
	if len(files) == 0 {
		return "", fmt.Errorf("proto file is empty")
	}
	if files[len(files)-1].Name == nil {
		return "", fmt.Errorf("proto file's name is null")
	}
	target = *files[len(files)-1].Name

	req := new(pluginpb.CodeGeneratorRequest)
	req.FileToGenerate = []string{target}
	var para = ""
	req.Parameter = &para
	req.ProtoFile = files

	g := generator.NewGenerator(s.opts.generatorOptions...)
	resp, err := g.Gen(req, onlyRPC)
	if err != nil {
		return "", err
	}
	if len(resp.File) == 0 {
		return "{}", nil
	}
	return *resp.File[0].Content, nil
}
