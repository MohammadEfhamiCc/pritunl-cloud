package domain

import (
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-cloud/database"
	"github.com/pritunl/pritunl-cloud/utils"
)

func Get(db *database.Database, domnId primitive.ObjectID) (
	domn *Domain, err error) {

	coll := db.Domains()
	domn = &Domain{}

	err = coll.FindOneId(domnId, domn)
	if err != nil {
		return
	}

	return
}

func GetOrg(db *database.Database, orgId, domnId primitive.ObjectID) (
	domn *Domain, err error) {

	coll := db.Domains()
	domn = &Domain{}

	err = coll.FindOne(db, &bson.M{
		"_id":          domnId,
		"organization": orgId,
	}).Decode(domn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func ExistsOrg(db *database.Database, orgId, domnId primitive.ObjectID) (
	exists bool, err error) {

	coll := db.Domains()

	n, err := coll.CountDocuments(db, &bson.M{
		"_id":          domnId,
		"organization": orgId,
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

func GetAll(db *database.Database, query *bson.M) (
	domns []*Domain, err error) {

	coll := db.Domains()
	domns = []*Domain{}

	cursor, err := coll.Find(db, query)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		dmn := &Domain{}
		err = cursor.Decode(dmn)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		domns = append(domns, dmn)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllPaged(db *database.Database, query *bson.M,
	page, pageCount int64) (domns []*Domain, count int64, err error) {

	coll := db.Domains()
	domns = []*Domain{}

	if len(*query) == 0 {
		count, err = coll.EstimatedDocumentCount(db)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	} else {
		count, err = coll.CountDocuments(db, query)
		if err != nil {
			err = database.ParseError(err)
			return
		}
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
		dmn := &Domain{}
		err = cursor.Decode(dmn)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		domns = append(domns, dmn)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllName(db *database.Database, query *bson.M) (
	domns []*Domain, err error) {

	coll := db.Domains()
	domns = []*Domain{}

	cursor, err := coll.Find(
		db,
		query,
		&options.FindOptions{
			Projection: &bson.D{
				{"name", 1},
				{"organization", 1},
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		dmn := &Domain{}
		err = cursor.Decode(dmn)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		domns = append(domns, dmn)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, domnId primitive.ObjectID) (err error) {
	coll := db.Domains()

	_, err = coll.DeleteOne(db, &bson.M{
		"_id": domnId,
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

func RemoveOrg(db *database.Database, orgId, domnId primitive.ObjectID) (
	err error) {

	coll := db.Domains()

	_, err = coll.DeleteOne(db, &bson.M{
		"_id":          domnId,
		"organization": orgId,
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

func RemoveMulti(db *database.Database, domnIds []primitive.ObjectID) (err error) {
	coll := db.Domains()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": domnIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func RemoveMultiOrg(db *database.Database, orgId primitive.ObjectID,
	domnIds []primitive.ObjectID) (err error) {

	coll := db.Domains()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": domnIds,
		},
		"organization": orgId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetRecordAll(db *database.Database, query *bson.M) (
	recrds []*Record, err error) {

	coll := db.DomainsRecord()
	recrds = []*Record{}

	cursor, err := coll.Find(db, query)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		recrd := &Record{}
		err = cursor.Decode(recrd)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		recrds = append(recrds, recrd)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func RemoveRecord(db *database.Database, recrdId primitive.ObjectID) (
	err error) {

	coll := db.DomainsRecord()

	_, err = coll.DeleteOne(db, &bson.M{
		"_id": recrdId,
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
