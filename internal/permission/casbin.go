package permission

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/galaxy-future/BridgX/config"
	"github.com/galaxy-future/BridgX/internal/logs"
)

var E *casbin.Enforcer

func Init() {
	E = initCasbin()
}

func initCasbin() *casbin.Enforcer {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.GlobalConfig.WriteDB.User,
		config.GlobalConfig.WriteDB.Password, config.GlobalConfig.WriteDB.Host, config.GlobalConfig.WriteDB.Port,
		config.GlobalConfig.WriteDB.Name)
	a, err := gormadapter.NewAdapter("mysql", dataSourceName, true)
	if err != nil {
		logs.Logger.Fatal(err.Error())
	}
	e, err := casbin.NewEnforcer("conf/model.conf", a)
	if err != nil {
		logs.Logger.Fatal(err.Error())
	}
	// Load the policy from DB.
	e.LoadPolicy()
	return e
}
