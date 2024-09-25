package uhandlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-cloud/database"
	"github.com/pritunl/pritunl-cloud/demo"
	"github.com/pritunl/pritunl-cloud/errortypes"
	"github.com/pritunl/pritunl-cloud/event"
	"github.com/pritunl/pritunl-cloud/service"
	"github.com/pritunl/pritunl-cloud/utils"
)

type serviceData struct {
	Id               primitive.ObjectID `json:"id"`
	Name             string             `json:"name"`
	Comment          string             `json:"comment"`
	Organization     primitive.ObjectID `json:"organization"`
	DeleteProtection bool               `json:"delete_protection"`
	Units            []*service.Unit    `json:"units"`
}

type servicesData struct {
	Services []*service.Service `json:"services"`
	Count    int64              `json:"count"`
}

func servicePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	userOrg := c.MustGet("organization").(primitive.ObjectID)
	data := &serviceData{}

	serviceId, ok := utils.ParseObjectId(c.Param("service_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	pd, err := service.GetOrg(db, userOrg, serviceId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	pd.Name = data.Name
	pd.Comment = data.Comment
	pd.DeleteProtection = data.DeleteProtection

	fields := set.NewSet(
		"id",
		"name",
		"comment",
		"delete_protection",
	)

	errData, err := pd.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = pd.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "service.change")

	c.JSON(200, pd)
}

func servicePost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	userOrg := c.MustGet("organization").(primitive.ObjectID)
	data := &serviceData{
		Name: "New Service",
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	pd := &service.Service{
		Name:             data.Name,
		Comment:          data.Comment,
		Organization:     userOrg,
		DeleteProtection: data.DeleteProtection,
		Units:            data.Units,
	}

	errData, err := pd.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = pd.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "service.change")

	c.JSON(200, pd)
}

func serviceDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	userOrg := c.MustGet("organization").(primitive.ObjectID)

	serviceId, ok := utils.ParseObjectId(c.Param("service_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := service.RemoveOrg(db, userOrg, serviceId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "service.change")

	c.JSON(200, nil)
}

func servicesDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	userOrg := c.MustGet("organization").(primitive.ObjectID)
	data := []primitive.ObjectID{}

	err := c.Bind(&data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	err = service.RemoveMultiOrg(db, userOrg, data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "service.change")

	c.JSON(200, nil)
}

func serviceGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	userOrg := c.MustGet("organization").(primitive.ObjectID)

	serviceId, ok := utils.ParseObjectId(c.Param("service_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	pd, err := service.GetOrg(db, userOrg, serviceId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, pd)
}

func servicesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	userOrg := c.MustGet("organization").(primitive.ObjectID)

	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

	query := bson.M{
		"organization": userOrg,
	}

	serviceId, ok := utils.ParseObjectId(c.Query("id"))
	if ok {
		query["_id"] = serviceId
	}

	name := strings.TrimSpace(c.Query("name"))
	if name != "" {
		query["name"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", regexp.QuoteMeta(name)),
			"$options": "i",
		}
	}

	comment := strings.TrimSpace(c.Query("comment"))
	if comment != "" {
		query["comment"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", comment),
			"$options": "i",
		}
	}

	role := strings.TrimSpace(c.Query("role"))
	if role != "" {
		query["role"] = role
	}

	services, count, err := service.GetAllPaged(db, &query, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &servicesData{
		Services: services,
		Count:    count,
	}

	c.JSON(200, data)
}
