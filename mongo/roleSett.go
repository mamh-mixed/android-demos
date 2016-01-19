package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

var RoleSettCol = roleSettCollection{name: "roleSett"}

// roleSettCollection 根绝清算角色清算
type roleSettCollection struct {
	name string
}

// Upsert
func (r *roleSettCollection) Upsert(rs *model.RoleSett) error {
	rs.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	find := bson.M{
		"reportType": rs.ReportType,
		"settRole":   rs.SettRole,
		"settDate":   rs.SettDate,
	}
	log.Debugf("save ... %+v", rs)
	_, err := database.C(r.name).Upsert(find, rs)
	return err
}

// PaginationFind 分页查找清算数据
func (r *roleSettCollection) PaginationFind(role, date string, reportType, size, page int) (results []model.RoleSett, total int, err error) {
	results = make([]model.RoleSett, 0)

	match := bson.M{}
	if role != "" {
		match["settRole"] = role
	}

	if date != "" {
		match["settDate"] = date
	}

	if reportType != 0 {
		match["reportType"] = reportType
	}

	total, err = database.C(r.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	cond := []bson.M{
		{"$match": match},
	}

	sort := bson.M{"$sort": bson.M{"settDate": -1}}

	skip := bson.M{"$skip": (page - 1) * size}

	limit := bson.M{"$limit": size}

	cond = append(cond, sort, skip, limit)

	err = database.C(r.name).Pipe(cond).All(&results)

	return results, total, err
}
