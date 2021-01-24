package disk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-cloud/database"
	"github.com/pritunl/pritunl-cloud/errortypes"
	"github.com/pritunl/pritunl-cloud/event"
	"github.com/pritunl/pritunl-cloud/paths"
	"github.com/pritunl/pritunl-cloud/utils"
	"github.com/sirupsen/logrus"
)

type Disk struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Comment          string             `bson:"comment" json:"comment"`
	State            string             `bson:"state" json:"state"`
	Node             primitive.ObjectID `bson:"node" json:"node"`
	Organization     primitive.ObjectID `bson:"organization,omitempty" json:"organization"`
	Instance         primitive.ObjectID `bson:"instance,omitempty" json:"instance"`
	SourceInstance   primitive.ObjectID `bson:"source_instance,omitempty" json:"source_instance"`
	DeleteProtection bool               `bson:"delete_protection" json:"delete_protection"`
	Image            primitive.ObjectID `bson:"image,omitempty" json:"image"`
	RestoreImage     primitive.ObjectID `bson:"restore_image,omitempty" json:"restore_image"`
	Backing          bool               `bson:"backing" json:"backing"`
	BackingImage     string             `bson:"backing_image" json:"backing_image"`
	Index            string             `bson:"index" json:"index"`
	Size             int                `bson:"size" json:"size"`
	NewSize          int                `bson:"new_size" json:"new_size"`
	Backup           bool               `bson:"backup" json:"backup"`
	LastBackup       time.Time          `bson:"last_backup" json:"last_backup"`
	curIndex         string             `bson:"-" json:"-"`
	curInstance      primitive.ObjectID `bson:"-" json:"-"`
}

func (d *Disk) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if !d.Instance.IsZero() && d.Index != "" {
		index, e := strconv.Atoi(d.Index)
		if e != nil {
			errData = &errortypes.ErrorData{
				Error:   "index_invalid",
				Message: "Disk index invalid",
			}
			return
		}

		if index < 0 || index > 10 {
			errData = &errortypes.ErrorData{
				Error:   "index_out_of_range",
				Message: "Disk index out of range",
			}
			return
		}

		d.Index = strconv.Itoa(index)
	}

	if d.Node.IsZero() {
		errData = &errortypes.ErrorData{
			Error:   "node_required",
			Message: "Missing required node",
		}
		return
	}

	if d.Backup && d.BackingImage != "" {
		errData = &errortypes.ErrorData{
			Error:   "backing_image_backup",
			Message: "Cannot enable backups with backing image",
		}
		return
	}

	if d.State == Restore && d.RestoreImage.IsZero() {
		errData = &errortypes.ErrorData{
			Error:   "restore_missing_image",
			Message: "Cannot restore without image set",
		}
		return
	}

	if d.Instance.IsZero() && !strings.HasPrefix(d.Index, "hold") {
		d.Index = fmt.Sprintf("hold_%s", primitive.NewObjectID().Hex())
	}

	if !d.Instance.IsZero() {
		disks, e := GetInstance(db, d.Instance)
		if e != nil {
			err = e
			return
		}

		for _, dsk := range disks {
			if dsk.Id != d.Id && dsk.Index == d.Index {
				errData = &errortypes.ErrorData{
					Error:   "disk_index_in_use",
					Message: "Disk index is already in use on instance",
				}
				return
			}
		}
	}

	if d.State == "" {
		d.State = Provision
	}

	if d.Size < 10 {
		d.Size = 10
	}

	if d.State == Expand {
		if d.NewSize == 0 {
			errData = &errortypes.ErrorData{
				Error:   "new_size_missing",
				Message: "Cannot expand without new size",
			}
			return
		}

		if d.NewSize < d.Size {
			errData = &errortypes.ErrorData{
				Error:   "new_size_invalid",
				Message: "New size cannot be less then current size",
			}
			return
		}
	} else {
		d.NewSize = 0
	}

	if d.DeleteProtection && d.curInstance != d.Instance {
		errData = &errortypes.ErrorData{
			Error:   "delete_protection_index",
			Message: "Cannot change instance with delete protection enabled",
		}
		return
	}

	if d.DeleteProtection && d.curIndex != d.Index {
		errData = &errortypes.ErrorData{
			Error:   "delete_protection_index",
			Message: "Cannot change index with delete protection enabled",
		}
		return
	}

	return
}

func (d *Disk) PreCommit() {
	d.curIndex = d.Index
	d.curInstance = d.Instance
}

func (d *Disk) Commit(db *database.Database) (err error) {
	coll := db.Disks()

	err = coll.Commit(d.Id, d)
	if err != nil {
		return
	}

	return
}

func (d *Disk) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Disks()

	err = coll.CommitFields(d.Id, d, fields)
	if err != nil {
		return
	}

	return
}

func (d *Disk) Insert(db *database.Database) (err error) {
	coll := db.Disks()

	_, err = coll.InsertOne(db, d)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (d *Disk) Destroy(db *database.Database) (err error) {
	dskPath := paths.GetDiskPath(d.Id)

	if d.DeleteProtection {
		logrus.WithFields(logrus.Fields{
			"disk_id": d.Id.Hex(),
		}).Info("disk: Delete protection ignore disk destroy")

		d.State = Available
		err = d.CommitFields(db, set.NewSet("state"))
		if err != nil {
			return
		}

		event.PublishDispatch(db, "disk.change")

		return
	}

	logrus.WithFields(logrus.Fields{
		"disk_id":   d.Id.Hex(),
		"disk_path": dskPath,
	}).Info("qemu: Destroying disk")

	err = utils.RemoveAll(dskPath)
	if err != nil {
		return
	}

	err = Remove(db, d.Id)
	if err != nil {
		return
	}

	return
}
