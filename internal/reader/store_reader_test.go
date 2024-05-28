package reader

import (
	"context"
	"errors"
	"io"
	"testing"
)

type extendsReader struct {
	buf io.Reader
}

type storeMock struct {
	data map[string]io.ReadCloser
}

func (e *extendsReader) Close() error {
	return nil
}

func (e *extendsReader) Read(p []byte) (int, error) {
	return e.buf.Read(p)
}

func (s *storeMock) Put(ctx context.Context, key string, data io.Reader) error {
	s.data[key] = &extendsReader{buf: data}
	return nil
}

func (s *storeMock) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if exists, ok := s.data[key]; !ok {
		return nil, errors.New("not found")
	} else {
		return exists, nil
	}
}

func (s *storeMock) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *storeMock) List(ctx context.Context) ([]string, error) {
	return nil, nil
}

func TestStoreWriter(t *testing.T) {
	t.Run("should read one", func(t *testing.T) {
	})
	t.Run("should read many", func(t *testing.T) {
	})
}
