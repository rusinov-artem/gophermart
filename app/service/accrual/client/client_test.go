package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinov-artem/gophermart/test/utils"
)

func Test_ClientCloseBody(t *testing.T) {
	ctx := context.Background()
	c := New(ctx, "localhost:8080")
	bodySpy := utils.NewReadCloserSpy()
	c.client = &fakeHTTPClient{bodySpy: bodySpy}

	_, _ = c.GetSingleOrder("orderNr")

	assert.True(t, bodySpy.IsClosed)
}

func Test_ClientHandleError(t *testing.T) {
	ctx := context.Background()
	c := New(ctx, "localhost:8080")
	c.client = &fakeHTTPClient{
		err: fmt.Errorf("network error"),
	}

	_, err := c.GetSingleOrder("orderNr")

	assert.ErrorContains(t, err, "unable to fetch order info")
}

func Test_ClientHandleBadJson(t *testing.T) {
	ctx := context.Background()
	c := New(ctx, "localhost:8080")
	body := utils.NewReadCloserSpy()
	body.Data = bytes.NewBufferString("}InvalidJson{")
	c.client = &fakeHTTPClient{
		bodySpy: body,
	}

	_, err := c.GetSingleOrder("orderNr")

	assert.ErrorContains(t, err, "unable to fetch order info")
}

type fakeHTTPClient struct {
	bodySpy *utils.ReadCloserSpy
	err     error
}

func (f *fakeHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{Body: f.bodySpy, StatusCode: http.StatusOK}, f.err
}
