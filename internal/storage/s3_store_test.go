package storage

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func (m *mockClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func (m *mockClient) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
}

func (m *mockClient) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func (m *mockClient) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.CreateBucketOutput), args.Error(1)
}

func TestS3Store(t *testing.T) {
	m := new(mockClient)

	m.On("PutObject", mock.Anything, mock.Anything).Return(&s3.PutObjectOutput{}, nil)
	m.On("GetObject", mock.Anything, mock.Anything).Return(&s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte{0, 1, 2, 3, 4, 5})),
	}, nil)

	s := NewS3Store(m, "test-bucket")

	t.Run("should put object", func(t *testing.T) {
		err := s.Put(context.Background(), "test-key", nil)
		assert.NoError(t, err)
	})

	t.Run("should get object", func(t *testing.T) {
		reader, err := s.Get(context.Background(), "test-key")
		assert.NoError(t, err)
		assert.NotNil(t, reader)

		data, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0, 1, 2, 3, 4, 5}, data)
	})

	t.Run("should delete object", func(t *testing.T) {
		// TODO @gebhartn: implement delete object test
		t.Skip()
	})

	t.Run("should list objects", func(t *testing.T) {
		// TODO @gebhartn: implement list objects test
		t.Skip()
	})
}
