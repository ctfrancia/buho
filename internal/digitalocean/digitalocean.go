package digitalocean

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// DigitalOceanSpacesClient manages uploads to DigitalOcean Spaces
type DigitalOceanSpacesClient struct {
	Client *minio.Client
	Bucket string
}

// NewDigitalOceanSpacesClient creates a new client for DigitalOcean Spaces
func NewDigitalOceanSpacesClient(endpoint, accessKeyID, secretAccessKey, bucket string) (*DigitalOceanSpacesClient, error) {
	// Create a new Minio client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true, // Use HTTPS
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Minio client: %v", err)
	}

	return &DigitalOceanSpacesClient{
		Client: client,
		Bucket: bucket,
	}, nil
}

// UploadFile uploads a file to DigitalOcean Spaces
func (d *DigitalOceanSpacesClient) UploadFile(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	// Upload the file with context
	info, err := d.Client.PutObject(ctx, d.Bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file to Spaces: %v", err)
	}

	// Construct the file URL
	fileURL := fmt.Sprintf("https://%s.%s/%s", d.Bucket, d.Client.EndpointURL().Host, info.Key)

	return fileURL, nil
}

// Additional utility methods
func (d *DigitalOceanSpacesClient) ListFiles(ctx context.Context) error {
	// List objects in the bucket
	objectCh := d.Client.ListObjects(ctx, d.Bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return fmt.Errorf("error listing objects: %v", object.Err)
		}
		log.Printf("Object: %s", object.Key)
	}

	return nil
}

func (d *DigitalOceanSpacesClient) DeleteFile(ctx context.Context, objectName string) error {
	// Delete a specific object from the bucket
	err := d.Client.RemoveObject(ctx, d.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("error deleting object: %v", err)
	}
	return nil
}
