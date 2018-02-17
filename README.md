# s3bolt
Golang's bolt db wrapper that automatically loads and backup to s3

## Purpose
s3bolt allows you to operate with bolt db with completelly transparent backup and restore to s3. It is intended for non-critical systems, due that:
- S3 is eventually-consistent. This means that after restoring from s3 you might not be loading last version.
- Every write operation triggers a push against s3, which clearly have an strong impact on performance.

## Usage
Use s3bolt is pretty simple. Here you have a simplified example:

```go
import "github.com/alanbover/s3bolt"

s3client, _ := s3bolt.NewS3Client(&s3bolt.SessionParameters{
    Region: "eu-west-1",
})
s3bolt := s3bolt.New(s3client, &s3bolt.Config{
    S3bucket: "someS3Bucket",
    S3prefix: "someS3Prefix",
})
db, _ := s3bolt.Open("/tmp/mydb", 0600, nil)
db.Update(func(tx *bolt.Tx) error {
    bucket, _ := tx.CreateBucketIfNotExists([]byte("someBucket"))
    bucket.Put([]byte("key"), []byte("value"))
    return nil
})
```
