package nut

import (
	"reflect"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

type networkCheck struct {
}

func (*networkCheck) Check() error {
	return nil
}

type databaseCheck struct {
}

func (*databaseCheck) Check() error {
	db, err := orm.GetDB()
	if err != nil {
		return err
	}
	return db.Ping()
}

type cacheCheck struct {
}

func (p *cacheCheck) Check() error {
	return Cache().Put("hi", AuthorName, 10*time.Second)
}

// RegisterHealthChecker register health checker
func RegisterHealthChecker(args ...toolbox.HealthChecker) {
	for _, it := range args {
		toolbox.AddHealthCheck(reflect.TypeOf(it).String(), it)
	}
}

func init() {
	RegisterHealthChecker(&databaseCheck{}, &cacheCheck{}, &networkCheck{})
}
