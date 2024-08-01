package domain

import (
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-cloud/database"
	"github.com/pritunl/pritunl-cloud/errortypes"
)

type Record struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain    primitive.ObjectID `bson:"domain" json:"domain"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	SubDomain string             `bson:"sub_domain" json:"sub_domain"`
	Type      string             `bson:"type" json:"type"`
	Value     string             `bson:"value" json:"value"`
	Operation string             `bson:"-" json:"operation"`
}

func (r *Record) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if r.Domain.IsZero() {
		errData = &errortypes.ErrorData{
			Error:   "domain_required",
			Message: "Missing required domain",
		}
		return
	}

	switch r.Type {
	case A:
		break
	case AAAA:
		break
	default:
		err = &errortypes.UnknownError{
			errors.New("domain: Unknown record type"),
		}
		return
	}

	if r.Value == "" {
		errData = &errortypes.ErrorData{
			Error:   "value_required",
			Message: "Missing required value",
		}
		return
	}

	return
}

func (r *Record) Commit(db *database.Database) (err error) {
	coll := db.DomainsRecords()

	err = coll.Commit(r.Id, r)
	if err != nil {
		return
	}

	return
}

func (r *Record) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.DomainsRecords()

	err = coll.CommitFields(r.Id, r, fields)
	if err != nil {
		return
	}

	return
}

func (r *Record) Insert(db *database.Database) (err error) {
	coll := db.DomainsRecords()

	if !r.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("domain: Record already exists"),
		}
		return
	}

	r.Id = primitive.NewObjectID()

	_, err = coll.InsertOne(db, r)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
