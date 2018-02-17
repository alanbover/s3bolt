package s3bolt_test

import (
	"bytes"
	"encoding/json"
	mock "github.com/alanbover/go-mockaws"
	"github.com/alanbover/s3bolt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/boltdb/bolt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

type ComplexObject struct {
	SomeArray []string
}

var testBucket = "somes3bucket"
var testKey = "somekey"
var s3prefix = "someprefix"
var dbpath = "/tmp/somedbpath"

func TestS3BoltWitMarshal(t *testing.T) {
	Convey("When create one empty s3bolt wrapper", t, func() {
		mock := initializeMock()
		s3bolt := s3bolt.New(mock, &s3bolt.Config{
			S3bucket: testBucket,
			S3prefix: s3prefix,
		})

		Convey("We should be able to write to and read from the database", func() {
			db, err := s3bolt.Open(dbpath, 0600, nil)
			So(err, ShouldBeNil)

			someobj := ComplexObject{
				SomeArray: []string{"somevalue1", "somevalue2"},
			}
			boltWriteValue(db, someobj, []byte(testKey), []byte(testBucket))
			value := boltReadValue(db, []byte(testKey), []byte(testBucket))
			So(value.SomeArray, ShouldHaveLength, 2)
			db.Close()

			Convey("We should be able to recover previous state", func() {
				os.Remove(dbpath)
				db, err := s3bolt.Open(dbpath, 0600, nil)
				So(err, ShouldBeNil)
				So(db, ShouldNotBeNil)
				value = boltReadValue(db, []byte(testKey), []byte(testBucket))
				So(value.SomeArray, ShouldHaveLength, 2)
				db.Close()
			})
		})
	})
}

func boltWriteValue(db *s3bolt.Db, value ComplexObject, key, bucket []byte) {
	encodedValue, err := json.Marshal(value)
	So(err, ShouldBeNil)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		err = bucket.Put(key, encodedValue)
		if err != nil {
			return err
		}
		return nil
	})
	So(err, ShouldBeNil)
}

func boltReadValue(db *s3bolt.Db, key, bucket []byte) ComplexObject {
	val := []byte{}
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		val = bucket.Get(key)
		return nil
	})
	So(err, ShouldBeNil)

	var decoded ComplexObject
	err = json.Unmarshal(val, &decoded)
	So(err, ShouldBeNil)

	return decoded
}

func initializeMock() *mock.MockS3 {
	mock := mock.NewMockS3()
	buffer := make([]byte, 1)
	// This creates an "empty" bucket.
	// Without it, open s3bolt.open will fail because bucket doesn't exist
	mock.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(testBucket),
		Key:    aws.String("someNonUsedKey"),
		Body:   bytes.NewReader(buffer),
	})
	return mock
}
