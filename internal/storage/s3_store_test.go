package storage

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

func (m *mockClient) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.HeadObjectOutput), args.Error(1)
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
	m.On("HeadObject", mock.Anything, mock.Anything).Return(&s3.HeadObjectOutput{}, nil)
	m.On("DeleteObject", mock.Anything, mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)
	m.On("ListObjectsV2", mock.Anything, mock.Anything).Return(&s3.ListObjectsV2Output{
		Contents: []types.Object{
			{
				Key: aws.String("user1/test_video1.mp4"),
			},
			{
				Key: aws.String("user1/test_video2.mp4"),
			},
		},
	}, nil)

	s := NewS3Store(m, "test-bucket")

	t.Run("should put object", func(t *testing.T) {
		err := s.Put(context.Background(), "test-key", "video/mp4", nil)
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

	t.Run("should head object", func(t *testing.T) {
		err := s.Head(context.Background(), "test-key")
		assert.NoError(t, err)
	})

	t.Run("should delete object", func(t *testing.T) {
		err := s.Delete(context.Background(), "test-key")
		assert.NoError(t, err)
	})

	t.Run("should list objects", func(t *testing.T) {
		list, err := s.List(context.Background(), "user1/")
		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})
}
