package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"strings"
)

func main() {
	marshaller := jsonpb.Marshaler{}

	ctx := context.Background()

	cc, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	refClient := grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(cc))

	methodDesc := getMethodDescriptor(refClient, "helloworld.Greeter", "SayHello")

	msgFactory := dynamic.NewMessageFactoryWithDefaults()
	req := msgFactory.NewMessage(methodDesc.GetInputType())

	if err := jsonpb.Unmarshal(strings.NewReader(`{"name": "Alice"}`), req); err != nil {
		panic(err)
	}

	stub := grpcdynamic.NewStubWithMessageFactory(cc, msgFactory)

	res, err := stub.InvokeRpc(ctx, methodDesc, req)
	if err != nil {
		panic(err)
	}

	json, err := marshaller.MarshalToString(res)
	if err != nil {
		panic(err)
	}
	fmt.Println(json)
}

func getMethodDescriptor(refClient *grpcreflect.Client, serviceName, methodName string) *desc.MethodDescriptor {
	svcDesc, err := refClient.ResolveService(serviceName)
	if err != nil {
		panic(err)
	}

	fmt.Println(svcDesc)
	for _, method := range svcDesc.GetMethods() {
		fmt.Println("===")
		fmt.Printf("METHOD : %+v\n", method.GetFullyQualifiedName())

		if method.GetName() != methodName {
			continue
		}

		fmt.Printf("IN  : %+v\n", method.GetInputType())
		for _, field := range method.GetInputType().GetFields() {
			fmt.Printf("> %+v : %+v [%+v]\n", field.GetNumber(), field.GetName(), field.GetType())
		}

		fmt.Printf("OUT : %+v\n", method.GetOutputType())
		for _, field := range method.GetOutputType().GetFields() {
			fmt.Printf("< %+v : %+v [%+v]\n", field.GetNumber(), field.GetName(), field.GetType())
		}

		return method
	}

	return nil
}
