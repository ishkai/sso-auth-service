package tests

import (
	"fmt"
	"sso/tests/suite"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func TestGRPC_HealthCheck(t *testing.T) {
	ctx, st := suite.New(t)

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", st.Cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	healthClient := healthgrpc.NewHealthClient(conn)

	resp, err := healthClient.Check(ctx, &healthgrpc.HealthCheckRequest{
		Service: "",
	})
	require.NoError(t, err)
	require.Equal(t, healthgrpc.HealthCheckResponse_SERVING, resp.Status)

	resp, err = healthClient.Check(ctx, &healthgrpc.HealthCheckRequest{
		Service: "auth.Auth",
	})
	require.NoError(t, err)
	require.Equal(t, healthgrpc.HealthCheckResponse_SERVING, resp.Status)
}

func TestGRPC_ReflectionAvailable(t *testing.T) {
	ctx, st := suite.New(t)

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", st.Cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	refClient := reflectionpb.NewServerReflectionClient(conn)

	stream, err := refClient.ServerReflectionInfo(ctx)
	require.NoError(t, err)

	err = stream.Send(&reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_ListServices{
			ListServices: "",
		},
	})
	require.NoError(t, err)

	resp, err := stream.Recv()

	services := resp.GetListServicesResponse().GetService()

	names := make([]string, 0, len(services))
	for _, service := range services {
		names = append(names, service.GetName())
	}

	assert.Contains(t, names, "auth.Auth")
	assert.Contains(t, names, "grpc.health.v1.Health")
}
