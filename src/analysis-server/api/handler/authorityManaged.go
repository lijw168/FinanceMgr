package handler

import (
	//"analysis-server/api/service"
	//"common/log"
	// "context"
	"errors"
	// "net/http"
	// "strings"
	// "sync"
	// "time"
)

//该authorityManaged是对操作接口的权限的管理。
type AuthorityManaged struct {
	// loginCheckMu      sync.RWMutex
	// tokenToOptIDMap   map[string]int
	// expirationCheckMu sync.RWMutex
	// tokenToTimeMap    map[string]int64
	// quitCheckCh       chan bool
	// authService       *service.AuthenService
	//logger *log.Logger
}

func NewAuthorityManaged() *AuthorityManaged {
	authManaged := AuthorityManaged{}
	return &authManaged
}

//api接口的权限鉴别
func (am *AuthorityManaged) InterfaceAuthorityCheck(action, accessToken string) (bool, error) {
	bIsPass := true
	var err error
	switch action {
	case "CreateCompany":
		fallthrough
	case "DeleteCompany":
		fallthrough
	case "UpdateCompany":
		fallthrough
	case "AssociatedCompanyGroup":
		fallthrough
	case "ListCompany":
		fallthrough
	case "GenerateAccSubTemplate":
		if !GAccessTokenH.isRootToken(accessToken) {
			bIsPass = false
			err = errors.New("No authority,access the function")
		}
		break
	case "ListLoginInfo":
		//区分root和普通的管理员的操作，通过是否有operatorId这个参数作为选择条件
		fallthrough
	case "CreateOperator":
		//区分创建管理员和普通的操作员，是通过给role赋值来实现的。
		fallthrough
	case "InitResourceInfo":
		fallthrough
	case "DeleteOperator":
		if !GAccessTokenH.isAdminToken(accessToken) {
			bIsPass = false
			err = errors.New("No authority,access the function")
		}
		break
	default:
		break
	}
	return bIsPass, err
}
