package metrics_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"manc/metrics"
)

func TestCounter(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	addr, ok := listener.Addr().(*net.TCPAddr)
	assert.True(t, ok)

	port := addr.Port
	require.NoError(t, err)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/increment", func(w http.ResponseWriter, r *http.Request) {
		metrics.Collector.IncrementCounter("my_counter", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"type": "example"}})
		metrics.Collector.IncrementCounter("my_counter", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"type": "example"}})
		metrics.Collector.IncrementCounter("my_counter_empty")
		metrics.Collector.IncrementCounter("my_counter_empty")
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	go func() {
		time.Sleep(time.Second)

		response, httpErr := http.Get(fmt.Sprintf("http://0.0.0.0:%d/increment", port))
		assert.NoError(t, httpErr)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		response, httpErr = http.Get(fmt.Sprintf("http://0.0.0.0:%d/metrics", port))
		assert.NoError(t, httpErr)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		body, httpErr := io.ReadAll(response.Body)
		assert.NoError(t, httpErr)

		// assertion
		want := fmt.Sprintf(`my_counter{type="example"} 2`)
		assert.Contains(t, string(body), want)

		// assertion
		want = fmt.Sprintf(`empty 2`)
		assert.Contains(t, string(body), want)

		assert.NoError(t, server.Shutdown(context.Background()))
	}()

	err = server.Serve(listener)
	require.Error(t, err)
	require.ErrorIs(t, err, http.ErrServerClosed)
}
