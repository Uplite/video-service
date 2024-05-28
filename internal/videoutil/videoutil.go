package videoutil

import "github.com/uplite/video-service/api/pb"

type ContentType = string

const (
	ContentTypeMp4  ContentType = "video/mp4"
	ContentTypeWebM ContentType = "video/webm"
	ContentTypeOgg  ContentType = "video/ogg"
)

func ContentTypeFrom(contentType pb.VideoContentType) ContentType {
	switch contentType {
	case pb.VideoContentType_VIDEO_CONTENT_TYPE_UNDEFINED:
		return ""
	case pb.VideoContentType_VIDEO_CONTENT_TYPE_MP4:
		return ContentTypeMp4
	case pb.VideoContentType_VIDEO_CONTENT_TYPE_WEBM:
		return ContentTypeWebM
	case pb.VideoContentType_VIDEO_CONTENT_TYPE_OGG:
		return ContentTypeOgg
	default:
		return ""
	}
}
