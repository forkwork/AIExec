package aliyun

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	aiexec_oss "aiexec-plugin/internal/oss"
)

type AliyunOSSStorage struct {
	client *oss.Client
	bucket *oss.Bucket
	path   string
}

func NewAliyunOSSStorage(
	region string,
	endpoint string,
	accessKeyID string,
	accessKeySecret string,
	authVersion string,
	path string,
	bucketName string,
) (*AliyunOSSStorage, error) {
	// options
	var options []oss.ClientOption

	// set region (required for v4)
	if region != "" {
		options = append(options, oss.Region(region))
	}

	// set auth-version
	if authVersion == "v1" {
		options = append(options, oss.AuthVersion(oss.AuthV1))
	} else if authVersion == "v4" {
		options = append(options, oss.AuthVersion(oss.AuthV4))
	} else {
		// default use v4
		options = append(options, oss.AuthVersion(oss.AuthV4))
	}

	// create client
	var client *oss.Client
	var err error

	client, err = oss.New(endpoint, accessKeyID, accessKeySecret, options...)

	if err != nil {
		return nil, fmt.Errorf("failed to create AliyunOSS client: %w", err)
	}

	// get specified bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket %s: %w", bucketName, err)
	}

	// normalize path: remove leading slash, ensure trailing slash
	path = strings.TrimPrefix(path, "/")
	if path != "" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return &AliyunOSSStorage{
		client: client,
		bucket: bucket,
		path:   path,
	}, nil
}

// combine full object path
func (s *AliyunOSSStorage) fullPath(key string) string {
	return path.Join(s.path, key)
}

func (s *AliyunOSSStorage) Save(key string, data []byte) error {
	fullPath := s.fullPath(key)
	err := s.bucket.PutObject(fullPath, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to put object to Aliyun OSS: %w", err)
	}
	return nil
}

func (s *AliyunOSSStorage) Load(key string) ([]byte, error) {
	fullPath := s.fullPath(key)
	object, err := s.bucket.GetObject(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from Aliyun OSS: %w", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data from Aliyun OSS: %w", err)
	}
	return data, nil
}

func (s *AliyunOSSStorage) Exists(key string) (bool, error) {
	fullPath := s.fullPath(key)
	exist, err := s.bucket.IsObjectExist(fullPath)
	if err != nil {
		return false, fmt.Errorf("failed to check if object exists in Aliyun OSS: %w", err)
	}
	return exist, nil
}

func (s *AliyunOSSStorage) State(key string) (aiexec_oss.OSSState, error) {
	fullPath := s.fullPath(key)
	meta, err := s.bucket.GetObjectMeta(fullPath)
	if err != nil {
		return aiexec_oss.OSSState{}, fmt.Errorf("failed to get object meta from Aliyun OSS: %w", err)
	}

	// Get content length
	size := int64(0)
	contentLength := meta.Get("Content-Length")
	if contentLength != "" {
		fmt.Sscanf(contentLength, "%d", &size)
	}

	// Get last modified time
	lastModified := time.Time{}
	lastModifiedStr := meta.Get("Last-Modified")
	if lastModifiedStr != "" {
		lastModified, _ = time.Parse(time.RFC1123, lastModifiedStr)
	}

	return aiexec_oss.OSSState{
		Size:         size,
		LastModified: lastModified,
	}, nil
}

func (s *AliyunOSSStorage) List(prefix string) ([]aiexec_oss.OSSPath, error) {
	// combine given prefix with path
	fullPrefix := path.Join(s.path, prefix)

	// Ensure the prefix ends with a slash for directories
	if !strings.HasSuffix(fullPrefix, "/") {
		fullPrefix = fullPrefix + "/"
	}

	var paths []aiexec_oss.OSSPath
	marker := ""
	for {
		lsRes, err := s.bucket.ListObjects(oss.Marker(marker), oss.Prefix(fullPrefix), oss.Delimiter("/"))
		if err != nil {
			return nil, fmt.Errorf("failed to list objects in Aliyun OSS: %w", err)
		}

		// Add files
		for _, object := range lsRes.Objects {
			if object.Key == fullPrefix {
				continue
			}
			// remove path and prefix from full path, only keep relative path
			key := strings.TrimPrefix(object.Key, fullPrefix)
			// Skip empty keys
			if key == "" {
				continue
			}
			paths = append(paths, aiexec_oss.OSSPath{
				Path:  key,
				IsDir: false,
			})
		}

		// Add directories
		for _, commonPrefix := range lsRes.CommonPrefixes {
			if commonPrefix == fullPrefix {
				continue
			}
			// remove path and prefix from full path, only keep relative path
			dirPath := strings.TrimPrefix(commonPrefix, fullPrefix)
			dirPath = strings.TrimSuffix(dirPath, "/")
			if dirPath == "" {
				continue
			}
			paths = append(paths, aiexec_oss.OSSPath{
				Path:  dirPath,
				IsDir: true,
			})
		}

		// Check if there are more results
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	return paths, nil
}

func (s *AliyunOSSStorage) Delete(key string) error {
	fullPath := s.fullPath(key)
	err := s.bucket.DeleteObject(fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete object from Aliyun OSS: %w", err)
	}
	return nil
}

func (s *AliyunOSSStorage) Type() string {
	return aiexec_oss.OSS_TYPE_ALIYUN_OSS
}
