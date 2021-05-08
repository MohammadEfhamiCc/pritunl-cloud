package image

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"time"

	"github.com/dropbox/godropbox/container/set"
	minio "github.com/minio/minio-go"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-cloud/database"
	"github.com/pritunl/pritunl-cloud/utils"
)

var (
	etagReg = regexp.MustCompile("[^a-zA-Z0-9]+")
)

func GetEtag(info minio.ObjectInfo) string {
	etag := info.ETag
	if etag == "" {
		modifiedHash := md5.New()
		modifiedHash.Write(
			[]byte(info.LastModified.Format(time.RFC3339)))
		etag = fmt.Sprintf("%x", modifiedHash.Sum(nil))
	}
	return etagReg.ReplaceAllString(etag, "")
}

func Get(db *database.Database, imgId primitive.ObjectID) (
	img *Image, err error) {

	coll := db.Images()
	img = &Image{}

	err = coll.FindOneId(imgId, img)
	if err != nil {
		return
	}

	return
}

func GetOrg(db *database.Database, orgId, imgId primitive.ObjectID) (
	img *Image, err error) {

	coll := db.Images()
	img = &Image{}

	err = coll.FindOne(db, &bson.M{
		"_id":          imgId,
		"organization": orgId,
	}).Decode(img)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetOrgPublic(db *database.Database, orgId, imgId primitive.ObjectID) (
	img *Image, err error) {

	coll := db.Images()
	img = &Image{}

	err = coll.FindOne(db, &bson.M{
		"_id": imgId,
		"$or": []*bson.M{
			&bson.M{
				"organization": orgId,
			},
			&bson.M{
				"organization": &bson.M{
					"$exists": false,
				},
			},
		},
	}).Decode(img)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Distinct(db *database.Database, storeId primitive.ObjectID) (
	keys []string, err error) {

	coll := db.Images()
	keys = []string{}

	keysInf, err := coll.Distinct(db, "key", &bson.M{
		"storage": storeId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	for _, keyInf := range keysInf {
		if key, ok := keyInf.(string); ok {
			keys = append(keys, key)
		}
	}

	return
}

func ExistsOrg(db *database.Database, orgId, imgId primitive.ObjectID) (
	exists bool, err error) {

	coll := db.Images()

	n, err := coll.CountDocuments(db, &bson.M{
		"_id": imgId,
		"$or": []*bson.M{
			&bson.M{
				"organization": orgId,
			},
			&bson.M{
				"organization": &bson.M{
					"$exists": false,
				},
			},
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if n > 0 {
		exists = true
	}

	return
}

func GetAll(db *database.Database, query *bson.M, page, pageCount int64) (
	imgs []*Image, count int64, err error) {

	coll := db.Images()
	imgs = []*Image{}

	count, err = coll.CountDocuments(db, query)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	maxPage := count / pageCount
	if count == pageCount {
		maxPage = 0
	}
	page = utils.Min64(page, maxPage)
	skip := utils.Min64(page*pageCount, count)

	cursor, err := coll.Find(
		db,
		query,
		&options.FindOptions{
			Sort: &bson.D{
				{"name", 1},
			},
			Skip:  &skip,
			Limit: &pageCount,
		},
	)
	defer cursor.Close(db)

	for cursor.Next(db) {
		img := &Image{}
		err = cursor.Decode(img)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		imgs = append(imgs, img)
		img = &Image{}
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllNames(db *database.Database, query *bson.M) (
	images []*Image, err error) {

	coll := db.Images()
	images = []*Image{}

	cursor, err := coll.Find(
		db,
		query,
		&options.FindOptions{
			Sort: &bson.D{
				{"name", 1},
			},
			Projection: &bson.D{
				{"name", 1},
				{"key", 1},
				{"firmware", 1},
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		img := &Image{}
		err = cursor.Decode(img)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		images = append(images, img)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllKeys(db *database.Database) (keys set.Set, err error) {
	coll := db.Images()
	keys = set.NewSet()

	cursor, err := coll.Find(
		db,
		&bson.M{},
		&options.FindOptions{
			Sort: &bson.D{
				{"name", 1},
			},
			Projection: &bson.D{
				{"_id", 1},
				{"etag", 1},
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		img := &Image{}
		err = cursor.Decode(img)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		keys.Add(fmt.Sprintf("%s-%s", img.Id.Hex(), img.Etag))
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, imgId primitive.ObjectID) (err error) {
	coll := db.Images()

	_, err = coll.DeleteOne(db, &bson.M{
		"_id": imgId,
	})
	if err != nil {
		err = database.ParseError(err)
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	return
}

func RemoveKeys(db *database.Database, storeId primitive.ObjectID,
	keys []string) (err error) {
	coll := db.Images()

	_, err = coll.DeleteMany(db, &bson.M{
		"storage": storeId,
		"key": &bson.M{
			"$in": keys,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
