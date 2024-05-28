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

func (s *storeMock) Put(ctx context.Context, key, contentType string, data io.Reader) error {
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

func (s *storeMock) Head(ctx context.Context, prefix string) error {
	return nil
}

func (s *storeMock) Delete(ctx context.Context, key string) error {
	if _, ok := s.data[key]; !ok {
		return errors.New("key does not exist")
	}
	delete(s.data, key)
	return nil
}

func (s *storeMock) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, nil
}

func TestStoreWriter(t *testing.T) {
	s := &storeMock{data: make(map[string]io.ReadCloser)}
	w := NewStoreWriter(s)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	t.Run("should write", func(t *testing.T) {
		data := []byte{0, 1, 2, 3, 4, 5}
		key := "key_1"
		ct := "video/mp4"

		err := w.Write(ctx, key, ct, bytes.NewBuffer(data))
		assert.NoError(t, err, "unexpected error while writing")

		stored, err := w.store.Get(ctx, key)
		assert.NoError(t, err, "unexpected error while getting stored data")

		storedBytes, err := io.ReadAll(stored)
		assert.NoError(t, err, "unexpected error while reading stored data")

		assert.Equal(t, data, storedBytes, "written bytes do not match payload bytes")
	})

	t.Run("should delete", func(t *testing.T) {
		data := []byte{0, 1, 2, 3, 4, 5}
		key := "key_2"
		ct := "video/mp4"

		err := w.Write(ctx, key, ct, bytes.NewBuffer(data))
		assert.NoError(t, err, "unexpected error while writing")

		err = w.Delete(ctx, key)
		assert.NoError(t, err, "unexpected error while deleting stored data ")

		_, err = w.store.Get(ctx, key)
		assert.Error(t, err)
	})
}