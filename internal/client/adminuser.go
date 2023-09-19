package client

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	stderr "errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/1219796395/myProject2/internal/biz/bo"
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type AdminUserClient struct {
	log     *log.Helper
	bizConf *conf.Biz
}

// NewEnvManageRepo .
func NewAdminUserClient(logger log.Logger, bc *conf.Bootstrap) *AdminUserClient {
	return &AdminUserClient{
		log:     log.NewHelper(logger),
		bizConf: bc.Biz,
	}
}

func (c *AdminUserClient) VerifySSOCode(ctx context.Context, code string, redirectUrl string) (token string, status int, err error) {
	// assemble form data
	formData := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     c.bizConf.Admin.Sso.ClientId,
		"client_secret": c.bizConf.Admin.Sso.ClientSecrete,
		"code":          code,
		"redirect_uri":  redirectUrl,
	}

	// send request
	resp, err := Client.R().
		SetContext(ctx).
		SetResult(&bo.SSOTokenRes{}).
		SetFormData(formData).
		Post(c.bizConf.Admin.Sso.Domain + bo.SSO_GET_TOKEN)
	if err != nil {
		return "", 0, err
	}

	// get token from response
	// parse status code
	statusCode := resp.StatusCode()
	if statusCode == http.StatusOK { // 成功
		resObj, _ := resp.Result().(*bo.SSOTokenRes)
		if resObj == nil {
			return "", 0, stderr.New("result from SSO is empty")
		}
		if resObj.Code != 0 || resObj.Data == nil {
			resJson, _ := json.MarshalToString(resObj)
			c.log.WithContext(ctx).Errorf("sso token validation error: %s", resJson)
			return "", bo.SSOLogInFailed, nil
		}
		return resObj.Data.Token, bo.SSOLogInSuccess, nil
	} else {
		return "", 0, stderr.New("http error: " + string(resp.Body()))
	}
}

func (c *AdminUserClient) GetUserInfoBySSOToken(ctx context.Context, token string) (info *bo.SSOInfoData, status int, err error) {
	resp, err := Client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&bo.SSOInfoRes{}).
		Get(c.bizConf.Admin.Sso.Domain + bo.SSO_GET_INFO_BY_TOKEN)
	if err != nil {
		return nil, 0, err
	}

	// get token from response
	// parse status code
	statusCode := resp.StatusCode()
	if statusCode == http.StatusOK { // 成功
		resObj, _ := resp.Result().(*bo.SSOInfoRes)
		if resObj == nil {
			return nil, 0, stderr.New("result from SSO is empty")
		}
		if resObj.Code != 0 || resObj.Data == nil {
			return nil, int(resObj.Code), nil
		}
		return resObj.Data, 0, nil
	}
	return nil, 0, stderr.New("http error: " + string(resp.Body()))
}

// SearchUser
// search list of user infos by either name/nickname/email
func (c *AdminUserClient) SearchUsersFromMdm(ctx context.Context, key string) ([]*bo.MdmUserInfo, error) {
	// sign
	time := fmt.Sprintf("%d", time.Now().UnixMilli())
	sign := c.sign(time)

	// set query params
	params := url.Values{}
	// sign params
	params.Set("appId", c.bizConf.Admin.Mdm.AppId)
	params.Set("time", time)
	params.Set("sign", sign)
	// bussiness params
	params.Set("cond", key)

	// send request
	resp, err := Client.R().SetContext(ctx).
		SetResult(&bo.MdmRes{}).
		SetQueryParamsFromValues(params).
		Get(c.bizConf.Admin.Mdm.Domain + bo.MDM_SEARCH_USERS)
	if err != nil {
		return nil, err
	}

	// parse status code
	statusCode := resp.StatusCode()
	if statusCode == http.StatusOK { // 成功
		resObj, _ := resp.Result().(*bo.MdmRes)
		if resObj == nil {
			return nil, stderr.New("result from mdm is empty")
		}
		if resObj.Data == nil {
			return nil, stderr.New(resObj.Message)
		}
		return resObj.Data, nil
	} else {
		return nil, stderr.New("http error: " + string(resp.Body()))
	}
}

// GetUsers
// get a batch of users by their hgids
func (c *AdminUserClient) QueryUserInfosFromMdm(ctx context.Context, ids []string) ([]*bo.MdmUserInfo, error) {
	// sign
	time := fmt.Sprintf("%d", time.Now().UnixMilli())
	sign := c.sign(time)

	// set query params
	params := url.Values{}
	// sign params
	params.Set("appId", c.bizConf.Admin.Mdm.AppId)
	params.Set("time", time)
	params.Set("sign", sign)
	// bussiness params
	hgIds := strings.Join(ids, ",")
	params.Set("hgIds", hgIds)

	// send request
	resp, err := Client.R().SetContext(ctx).
		SetResult(&bo.MdmRes{}).
		SetQueryParamsFromValues(params).
		Get(c.bizConf.Admin.Mdm.Domain + bo.MDM_QUERY_USERS)
	if err != nil {
		return nil, err
	}

	// parse status code
	statusCode := resp.StatusCode()
	if statusCode == http.StatusOK { // 成功
		resObj, _ := resp.Result().(*bo.MdmRes)
		if resObj == nil {
			return nil, stderr.New("result from mdm is empty")
		}
		if resObj.Data == nil {
			return nil, stderr.New(resObj.Message)
		}
		return resObj.Data, nil
	} else {
		paramJson, _ := json.MarshalToString(params)
		c.log.WithContext(ctx).Errorf("[AdminUserClient.QueryUserInfosFromMdm] mdm request: %s, url: %s",
			paramJson, c.bizConf.Admin.Mdm.Domain+bo.MDM_QUERY_USERS)
		c.log.WithContext(ctx).Errorf("[AdminUserClient.QueryUserInfosFromMdm] mdm response: %s", resp.Body())
		return nil, stderr.New("http error: " + string(resp.Body()))
	}
}

// UpperCase(MD5（appSecret + appId + time + appSecret）)
func (c *AdminUserClient) sign(time string) string {
	// assemble sign string
	signStr := c.bizConf.Admin.Mdm.Secrete + c.bizConf.Admin.Mdm.AppId + time + c.bizConf.Admin.Mdm.Secrete

	// do md5 encoding
	h := md5.New()
	h.Write([]byte(signStr))
	sign := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(sign)
}
