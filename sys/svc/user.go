package svc

import (
	"fmt"
	"time"

	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

type UserSvc struct {
	*svc.DataSvc[mdl.User]
	Perm      *casbin.SyncedEnforcer
	SignedKey []byte
}

func NewUserSvc(i *do.Injector) (*UserSvc, error) {
	s := &UserSvc{
		DataSvc:   svc.NewDataSvc[mdl.User]("sys:users"),
		SignedKey: []byte("kskj"),
	}
	model := model.NewModel()
	model.AddDef("r", "r", "sub, obj, act")                                                                                                                    // 请求（sub：用户编码，obj：请求路径，act：方法）
	model.AddDef("p", "p", "sub, obj, act, res")                                                                                                               // 角色资源（sub：角色编码，obj：请求路径，act：方法，id：资源编码，type：资源类型）
	model.AddDef("g", "g", "_, _")                                                                                                                             // 用户角色（_：用户编码，_：角色编码）
	model.AddDef("e", "e", "some(where (p.eft == allow))")                                                                                                     // 策略
	model.AddDef("m", "m", "r.sub == '300000000000001' || r.obj == '/mgt/sys/resources/menu' || g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act") // 匹配
	adapter, err := gormadapter.NewAdapterByDBUseTableName(s.Make(), "sys", "rule")
	if err != nil {
		logrus.Fatal(err)
	}
	perm, err := casbin.NewSyncedEnforcer(model, adapter)
	if err != nil {
		logrus.Fatal(err)
	}
	err = perm.LoadPolicy()
	if err != nil {
		logrus.Fatal(err)
	}
	s.Perm = perm
	return s, nil
}

func (s *UserSvc) ListUsers(conds map[string]any) ([]mdl.User, int64, error) {
	var users []mdl.User
	var count int64
	db := s.Make(conds).Joins("Dept").Find(&users)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	err := db.Error
	if err != nil {
		return users, count, err
	}
	return users, count, nil
}

func (s *UserSvc) ListUserRoles(id string) ([]int64, error) {
	roles := []int64{}
	rs, err := s.Perm.GetRolesForUser(id)
	if err != nil {
		return roles, err
	}
	for _, r := range rs {
		roles = append(roles, cast.ToInt64(r))
	}
	return roles, nil
}

func (s *UserSvc) Login(userName, password string) (*mdl.User, error) {
	var user *mdl.User
	err := s.Make().Select("id, user_name, password, frozen").First(&user, "user_name = ?", userName).Error
	if err != nil {
		return nil, fmt.Errorf("UserWrong")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("UserWrong")
	}
	if user.Frozen {
		return nil, fmt.Errorf("UserFrozen")
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":    user.Id,
			"userName":  user.UserName,
			"expiresAt": time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	)
	user.Password = ""
	user.Token, _ = token.SignedString(s.SignedKey)
	return user, nil
}

func (s *UserSvc) RemoveUsers(ids []string) error {
	return s.Remove(ids, "id != 300000000000001")
}

func (s *UserSvc) SaveUser(user mdl.User) error {
	if len(user.Password) != 60 {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(password)
	}
	return s.Save(user)
}

func (s *UserSvc) SaveUserRoles(id string, roles []int64) error {
	_, err := s.Perm.DeleteRolesForUser(id)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return nil
	}
	_, err = s.Perm.AddRolesForUser(id, cast.ToStringSlice(roles))
	return err
}
