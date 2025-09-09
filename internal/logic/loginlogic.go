package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 3.完整的
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	client := l.client()
	//从ccnu主页获取相关参数
	params, err := l.makeAccountPreflightRequest(client)
	if err != nil {
		return nil, err
	}

	client, err = l.loginClient(l.ctx, client, req.Username, req.Password, params)
	if err != nil {
		return nil, err
	}

	username, err := l.GetNameFromXK(client, params)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserModel.FindScoreByID(l.ctx, req.Username)
	if err != nil {
		newUser := &model.User{
			StudentId: req.Username,
			Score:     100,
		}
		_, err = l.svcCtx.UserModel.Insert(l.ctx, newUser)
		if err != nil {
			return nil, errorx.WrapError(errorx.CreateErr, err)
		}
	}

	resp = &types.LoginRes{
		Base: response.Success(),
		Data: types.LoginData{
			Name:      username,
			StudentId: req.Username,
		},
	}

	return resp, nil
}

// 2.教务系统获取姓名
func (l *LoginLogic) GetNameFromXK(client *http.Client, params *accountRequestParams) (string, error) {
	type Data struct {
		XM string `json:"xm"`
	}
	type Resp struct {
		Data Data `json:"data"`
	}
	var resp1 Resp

	//华师本科生院登陆
	request, err := http.NewRequest("POST", "https://grd.ccnu.edu.cn/yjsxt/xtgl/index_cxUserInfo.html?gnmkdm=index", nil)
	if err != nil {
		return "", errorx.WrapError(errorx.CCNUSERVER_ERROR, err)
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")
	request.Header.Set("cookie", "JSESSIONID="+params.JSESSIONID+";route="+params.route)

	resp, err := client.Do(request)
	if err != nil {
		return "", errorx.WrapError(errorx.CCNUSERVER_ERROR, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	err = json.Unmarshal(body, &resp1)
	if err != nil {
		return "", errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	return resp1.Data.XM, nil
}

// 1.登录研究生教务系统
func (l *LoginLogic) loginClient(ctx context.Context, client *http.Client, studentId string, password string, params *accountRequestParams) (*http.Client, error) {
	encrypted, err := EncryptPassword(params.RSAParams.Modulus, params.RSAParams.Exponent, password)
	if err != nil {
		return nil, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	v := url.Values{}
	v.Set("csrftoken", params.csrftoken)
	v.Set("yhm", studentId)
	v.Set("mm", encrypted)
	v.Set("mm", encrypted)

	request, _ := http.NewRequest("POST", "https://grd.ccnu.edu.cn/yjsxt/xtgl/login_slogin.html", strings.NewReader(v.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("cookie", "JSESSIONID="+params.JSESSIONID+";route="+params.route)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")
	request = request.WithContext(ctx)
	//创建一个带jar的客户端
	j, _ := cookiejar.New(&cookiejar.Options{})
	client.Jar = j
	//发送请求
	resp, err := client.Do(request)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			return nil, errorx.WrapError(errorx.CCNUSERVER_ERROR, err)
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Request.URL.String() != "http://grd.ccnu.edu.cn/yjsxt/xtgl/index_initMenu.html" {
		return nil, errorx.WrapError(errorx.PasswordErr, errors.New("学号或密码错误"))
	}

	return client, nil
}

type accountRequestParams struct {
	csrftoken  string
	JSESSIONID string
	route      string
	RSAParams  RSAParams
}

type RSAParams struct {
	Exponent string `json:"exponent"`
	Modulus  string `json:"modulus"`
}

// -1.获取必要参数
func (l *LoginLogic) makeAccountPreflightRequest(client *http.Client) (*accountRequestParams, error) {
	var csrftoken string
	var JSESSIONID string
	var route string
	var RSAParams RSAParams

	params := &accountRequestParams{}

	request, err := http.NewRequest("GET", "https://grd.ccnu.edu.cn/yjsxt/xtgl/login_slogin.html", nil)
	if err != nil {
		return params, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")

	resp, err := client.Do(request)
	if err != nil {
		return params, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			JSESSIONID = cookie.Value
		} else if cookie.Name == "route" {
			route = cookie.Value
		}
	}

	csrftoken = ""

	// 获取rsa密码参数
	requestRSA, err := http.NewRequest("GET", "https://grd.ccnu.edu.cn/yjsxt/xtgl/login_getPublicKey.html", nil)
	if err != nil {
		return nil, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	requestRSA.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")

	respRSA, err := client.Do(requestRSA)
	if err != nil {
		return nil, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	bodyRSA, err := io.ReadAll(respRSA.Body)
	if err != nil {
		return nil, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	err = json.Unmarshal(bodyRSA, &RSAParams)
	if err != nil {
		return nil, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	return &accountRequestParams{
		csrftoken:  csrftoken,
		JSESSIONID: JSESSIONID,
		route:      route,
		RSAParams:  RSAParams,
	}, nil
}

// -1.前置工作,用于初始化client
func (l *LoginLogic) client() *http.Client {
	j, _ := cookiejar.New(&cookiejar.Options{})
	return &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Jar: j,
	}
}

// tools 登录过程所需密码加密参数
func EncryptPassword(modulusB64, exponentB64, password string) (string, error) {
	// base64 → hex
	modHex, err := b64tohex(modulusB64)
	if err != nil {
		return "", err
	}
	expHex, err := b64tohex(exponentB64)
	if err != nil {
		return "", err
	}

	// hex → 大整数
	n := new(big.Int)
	n.SetString(modHex, 16)

	e := new(big.Int)
	e.SetString(expHex, 16)

	pub := rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	// PKCS#1 v1.5 加密
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, &pub, []byte(password))
	if err != nil {
		return "", err
	}

	// hex2b64 → base64 输出
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func b64tohex(b64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}
