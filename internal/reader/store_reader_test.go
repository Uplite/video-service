package reader

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uplite/video-service/internal/config"
)

type storeMock struct {
	data map[string]io.Reader
}

func (s *storeMock) Put(ctx context.Context, key, contentType string, data io.Reader) error {
	s.data[key] = data
	return nil
}

func (s *storeMock) Head(ctx context.Context, key string) error {
	return nil
}

func (s *storeMock) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *storeMock) List(ctx context.Context, prefix string) ([]string, error) {
	return []string{"video_1", "video_2"}, nil
}

func TestStoreReader(t *testing.T) {
	s := &storeMock{data: make(map[string]io.Reader)}
	r := NewStoreReader(s)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	t.Run("should read one", func(t *testing.T) {
		key := "video_1"

		s, err := r.ReadOne(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, s, appendDelimiter(config.GetCloudFrontURL())+key)
	})
	t.Run("should read many", func(t *testing.T) {
		prefix := "user1"

		us, err := r.ReadMany(ctx, prefix)
		assert.NoError(t, err)
		assert.Len(t, us, 2)

		assert.Equal(t, us[0], appendDelimiter(config.GetCloudFrontURL())+"video_1")
		assert.Equal(t, us[1], appendDelimiter(config.GetCloudFrontURL())+"video_2")
	})
}
