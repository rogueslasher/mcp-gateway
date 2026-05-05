/*
Package clients provides a set of clients for use with the gateway code
*/
package clients

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Kuadrant/mcp-gateway/internal/config"
	"github.com/Kuadrant/mcp-gateway/internal/tests/server2"
	"github.com/stretchr/testify/require"
)

func getFreePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port
}

func TestInitialize(t *testing.T) {
	testCases := []struct {
		name               string
		gatewayHost        string
		routerKey          string
		conf               *config.MCPServer
		passThroughHeaders map[string]string
		expectedError      bool
	}{
		{
			name:        "standard initialization",
			gatewayHost: "%invalid",
			routerKey:   "router-key-123",
			conf: &config.MCPServer{
				Name:       "test-server",
				ToolPrefix: "test_",
				Hostname:   "test.mcp.local",
			},
			passThroughHeaders: map[string]string{},
			expectedError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := Initialize(context.Background(), tc.gatewayHost, tc.routerKey, tc.conf, tc.passThroughHeaders, false)
			if tc.expectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, client)
		})
	}
}

func TestInitialize_WithMockServer(t *testing.T) {
	port := getFreePort(t)
	startup, shutdown, err := server2.RunServer("http", fmt.Sprintf("%d", port))
	require.NoError(t, err)

	go func() {
		_ = startup()
	}()
	t.Cleanup(func() {
		_ = shutdown()
	})

	time.Sleep(100 * time.Millisecond)

	conf := &config.MCPServer{
		Name:       "test-server",
		ToolPrefix: "test_",
		Hostname:   "test.mcp.local",
		URL:        fmt.Sprintf("http://localhost:%d/mcp", port),
	}

	client, err := Initialize(context.Background(), fmt.Sprintf("localhost:%d", port), "router-key-123", conf, map[string]string{}, false)
	require.NoError(t, err)
	require.NotNil(t, client)
}
