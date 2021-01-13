package data

import (
	"context"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	minio "github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
	"github.com/pritunl/pritunl-cloud/database"
	"github.com/pritunl/pritunl-cloud/errortypes"
	"github.com/pritunl/pritunl-cloud/image"
	"github.com/pritunl/pritunl-cloud/storage"
	"github.com/pritunl/pritunl-cloud/utils"
	"github.com/sirupsen/logrus"
)

var (
	syncLock = utils.NewMultiTimeoutLock(1 * time.Minute)
)

func Sync(db *database.Database, store *storage.Storage) (err error) {
	if store.Endpoint == "" {
		return
	}

	lockId := syncLock.Lock(store.Id.Hex())
	defer syncLock.Unlock(store.Id.Hex(), lockId)

	client, err := minio.New(store.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(store.AccessKey, store.SecretKey, ""),
		Secure: !store.Insecure,
	})
	if err != nil {
		err = &errortypes.ConnectionError{
			errors.Wrap(err, "storage: Failed to connect to storage"),
		}
		return
	}

	images := []*image.Image{}
	signedKeys := set.NewSet()
	remoteKeys := set.NewSet()
	for object := range client.ListObjects(
		context.Background(),
		store.Bucket, minio.ListObjectsOptions{
			Recursive: true,
		},
	) {

		if object.Err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(object.Err, "storage: Failed to list objects"),
			}
			return
		}

		if strings.HasSuffix(object.Key, ".qcow2.sig") {
			signedKeys.Add(strings.TrimRight(object.Key, ".sig"))
		} else if strings.HasSuffix(object.Key, ".qcow2") {
			etag := image.GetEtag(object)
			remoteKeys.Add(object.Key)

			img := &image.Image{
				Storage:      store.Id,
				Key:          object.Key,
				Firmware:     image.Unknown,
				Etag:         etag,
				Type:         store.Type,
				LastModified: object.LastModified,
			}

			if store.IsOracle() {
				obj, e := client.StatObject(context.Background(),
					store.Bucket, object.Key, minio.StatObjectOptions{})
				if e != nil {
					err = &errortypes.ReadError{
						errors.Wrap(e, "storage: Failed to stat object"),
					}
					return
				}

				img.StorageClass = storage.ParseStorageClass(obj)
			} else {
				img.StorageClass = storage.ParseStorageClass(object)
			}

			images = append(images, img)
		}
	}

	for _, img := range images {
		img.Signed = signedKeys.Contains(img.Key)

		if img.Signed {
			if strings.Contains(img.Key, "_efi") ||
				strings.Contains(img.Key, "_uefi") {

				img.Firmware = image.Uefi
			} else {
				img.Firmware = image.Bios
			}
		}

		err = img.Sync(db)
		if err != nil {
			if _, ok := err.(*image.LostImageError); ok {
				logrus.WithFields(logrus.Fields{
					"bucket": store.Bucket,
					"key":    img.Key,
				}).Error("data: Ignoring lost image")
			} else {
				return
			}
		}
	}

	localKeys, err := image.Distinct(db, store.Id)
	if err != nil {
		return
	}

	removeKeysSet := set.NewSet()
	for _, key := range localKeys {
		removeKeysSet.Add(key)
	}
	removeKeysSet.Subtract(remoteKeys)

	removeKeys := []string{}
	for key := range removeKeysSet.Iter() {
		removeKeys = append(removeKeys, key.(string))
	}

	err = image.RemoveKeys(db, store.Id, removeKeys)
	if err != nil {
		return
	}

	return
}
