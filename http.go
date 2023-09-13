package gohue

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func debugHttpResponse(ctx context.Context, resp *http.Response) error {
	logger := getDebugValue(ctx)
	if logger == nil {
		return nil
	}

	bodyCpy := bytes.NewBuffer(nil)
	_, err := io.Copy(bodyCpy, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying response body: %w", err)
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bodyCpy)

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return fmt.Errorf("error dumping response: %w", err)
	}

	// TODO: Improve this
	logger.Info(string(dump))
	return nil
}
