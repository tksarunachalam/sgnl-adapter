// Copyright 2023 SGNL.ai, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net"
	"os"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	"github.com/sgnl-ai/adapter-template/pkg/adapter"
	"github.com/sgnl-ai/adapter-template/pkg/example_datasource"
	"google.golang.org/grpc"
)

const (
	// Port is the port at which the gRPC server will listen.
	//
	// SCAFFOLDING:
	// Modify this port as needed, or make it configurable.
	ServerPort = 8080
)

func main() {
	logger := log.New(os.Stdout, "adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerPort))
	if err != nil {
		logger.Fatalf("Failed to open server port: %v", err)
	}

	// SCAFFOLDING:
	// Pass options to configure TLS, etc.
	s := grpc.NewServer()

	// SCAFFOLDING:
	// This directly connects the adapter to an in-memory example datasource.
	adapter := adapter.NewAdapter(logger, example_datasource.NewClient())

	api_adapter_v1.RegisterAdapterServer(s, server.New(adapter))

	logger.Printf("Started adapter gRPC server on port %d", ServerPort)

	if err := s.Serve(listener); err != nil {
		logger.Fatalf("Failed to listen on server port: %v", err)
	}
}
