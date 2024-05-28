package writer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	s := &storeMock{data: make(map[string]io.ReadCloser)}
	r := NewStoreWriter(s)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	t.Run("should write", func(t *testing.T) {
		data := []byte{0, 1, 2, 3, 4, 5}
		key := "key_1"

		err := r.Write(ctx, key, bytes.NewBuffer(data))
		assert.NoError(t, err, "unexpected error while writing")

		stored, err := r.store.Get(ctx, key)
		assert.NoError(t, err, "unexpected error while getting stored data")

		storedBytes, err := io.ReadAll(stored)
		assert.NoError(t, err, "unexpected error while reading stored data")

		assert.Equal(t, data, storedBytes, "written bytes do not match payload bytes")
	})
}
