package s3

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)
type S3 struct {
	Endpoint  string
	Accesskey string
	Secretkey string
	UseSSL bool
	S3Client *minio.Client
}

// 新建s3 连接
func NewS3Conn(endpoint string, accessKeyID string,secretAccessKey string,useSSL bool)  *S3{
	Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	//log.Printf("%#v\n", Client)
	return &S3{
		Endpoint: endpoint,
		Accesskey: accessKeyID,
		Secretkey: secretAccessKey,
		S3Client: Client,
		UseSSL: useSSL,
	}
}

//上传对象
func (s *S3)CopyObject()  {

}
//列出桶或者对象
func (s *S3)List()  {
	// 列出所有桶
	s.ListBuckets()
}

func (s *S3)ListBuckets()  {
	buckets, err := s.S3Client.ListBuckets(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, bucket := range buckets {
		fmt.Println(bucket.CreationDate.Format("2006-01-02 15:04:05"),"		",bucket.Name)
	}
}

func (s *S3)ListObjects(bucketName string)  {
	//ctx, cancel := context.WithCancel(context.Background())
	//
	//defer cancel()
	//
	//objectCh := s.S3Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
	//	Prefix: "myprefix",
	//	Recursive: true,
	//})
	//for object := range objectCh {
	//	if object.Err != nil {
	//		fmt.Println(object.Err)
	//		return
	//	}
	//	fmt.Println(object)
	//}
}
//创建桶
func (s *S3)MakeBucket(bucketName string)  {
	err := s.S3Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
	fmt.Println("Successfully created mybucket.")
	if err != nil {
		log.Println("创建bucket错误: ",err)
		exists, _ := s.S3Client.BucketExists(context.Background(), bucketName)
		if exists {
			log.Printf("bucket: %s已经存在",bucketName)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}
//删除桶
func (s *S3)RemoveBucket()  {

}