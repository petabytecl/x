package checkerx

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type serverInfoResponse struct {
	Version string `json:"version"`
	State   string `json:"state"`
}

// mockIstioProxy is a mock for Proxy interface.
type mockIstioProxy struct{}

// Wait mock
func (m *mockIstioProxy) Wait() error {
	return nil
}

// Stop mock
func (m *mockIstioProxy) Stop() error {
	return nil
}

// CheckHealth mock
func (m *mockIstioProxy) CheckHealth() error {
	return nil
}

// istioProxy implements the IstioProxy interface
type istioProxy struct {
	serverInfoAddress string
	cancel            context.CancelFunc

	client     *http.Client
	retryDelay time.Duration
	maxWait    time.Duration

	stopped <-chan struct{}
	healthy bool
}

// NewIstioChecker returns a new IstioProxy interface
// based on ISTIO_PROXY_ENABLED environment variable
// If "ISTIO_PROXY_ENABLED" != "true" a mock of the Proxy
// interface will be returned.
func NewIstioChecker(timeout, retryDelay, maxWait time.Duration) HealthChecker {
	var enabled = false
	val, ok := os.LookupEnv("ISTIO_PROXY_ENABLED")
	if ok {
		result, err := strconv.ParseBool(val)
		if err == nil {
			enabled = result
		}
	}

	if !enabled {
		return &mockIstioProxy{}
	}

	stopped := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	i := &istioProxy{
		serverInfoAddress: "http://localhost:15000/server_info",
		client:            &http.Client{Timeout: timeout},
		retryDelay:        retryDelay,
		maxWait:           maxWait,
		cancel:            cancel,
		stopped:           stopped,
	}

	go func() {
		var timer = time.NewTimer(i.retryDelay)
		defer func() {
			timer.Stop()
			close(stopped)
		}()

		for {
			if err := i.execute(); err != nil {
				i.healthy = true
				return
			}

			timer.Reset(i.retryDelay)
			select {
			case <-timer.C:
				if i.retryDelay < i.maxWait {
					i.retryDelay *= 2
					if i.retryDelay > i.maxWait {
						i.retryDelay = i.maxWait
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return i
}

// execute check
func (i *istioProxy) execute() error {
	response, err := i.client.Get(i.serverInfoAddress)
	if err != nil {
		return errors.Wrap(err, "unable to get response from the proxy")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response body")
	}

	serverResponse := serverInfoResponse{}
	if err := jsoniter.Unmarshal(responseBody, &serverResponse); err != nil {
		return errors.Wrap(err, "unable to unmarshal json response")
	}

	if serverResponse.State != "LIVE" {
		return errors.New("proxy is not ready")
	}

	return nil
}

// CheckHealth of the istio-proxy
func (i *istioProxy) CheckHealth() error {
	select {
	case <-i.stopped:
		if !i.healthy {
			return errors.New("stopped before becoming healthy")
		}
		return nil
	default:
		return errors.New("service unavailable")
	}
}

// Stop istio-proxy health check
func (i *istioProxy) Stop() error {
	i.cancel()
	<-i.stopped

	return nil
}

// Wait until the istio-proxy is ready
func (i *istioProxy) Wait() error {
	select {
	case <-i.stopped:
		if !i.healthy {
			return errors.New("stopped before becoming healthy")
		}
		return nil
	default:
		return errors.New("service unavailable")
	}
}
