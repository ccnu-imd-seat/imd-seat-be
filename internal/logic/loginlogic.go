package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/PuerkitoBio/goquery"
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

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	if string(req.Username[4]) == "2" {
		return l.LoginBK(req)
	} else {
		return l.LoginPGraduate(req)
	}
}

// 3.完整的
func (l *LoginLogic) LoginPGraduate(req *types.LoginReq) (resp *types.LoginRes, err error) {
	client := l.client()
	//从ccnu主页获取相关参数
	params, err := l.makePGAccountPreflightRequest(client)
	if err != nil {
		return nil, err
	}

	client, err = l.loginPGClient(l.ctx, client, req.Username, req.Password, params)
	if err != nil {
		return nil, err
	}

	username, err := l.GetNameFromPGXK(client, params)
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
func (l *LoginLogic) GetNameFromPGXK(client *http.Client, params *PGaccountRequestParams) (string, error) {
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
func (l *LoginLogic) loginPGClient(ctx context.Context, client *http.Client, studentId string, password string, params *PGaccountRequestParams) (*http.Client, error) {
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

type PGaccountRequestParams struct {
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
func (l *LoginLogic) makePGAccountPreflightRequest(client *http.Client) (*PGaccountRequestParams, error) {
	var csrftoken string
	var JSESSIONID string
	var route string
	var RSAParams RSAParams

	params := &PGaccountRequestParams{}

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

	return &PGaccountRequestParams{
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

// 以下是本科生的登录获取

// 3.完整的
func (l *LoginLogic) LoginBK(req *types.LoginReq) (resp *types.LoginRes, err error) {
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

	client, err = l.xkLoginClient(client)
	if err != nil {
		return nil, err
	}

	username, err := l.GetNameFromXK(l.ctx, client)
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

// 2.教务系统获取名字
func (l *LoginLogic) GetNameFromXK(ctx context.Context, client *http.Client) (string, error) {
	requestUrl := "https://xk.ccnu.edu.cn/jwglxt/xtgl/index_cxYhxxIndex.html?xt=jw&localeKey=zh_CN&gnmkdm=index"

	resp, err := client.Get(requestUrl)
	if err != nil {
		return "", errorx.WrapError(errorx.CCNUSERVER_ERROR, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errorx.WrapError(errorx.CCNUSERVER_ERROR, fmt.Errorf("HTML解析失败: %v", err))
	}

	// 精准定位目标元素
	nameSection := doc.Find("div.media-body h4.media-heading").First()
	if nameSection.Length() == 0 {
		return "", errorx.WrapError(errorx.CCNUSERVER_ERROR, errors.New("未找到姓名元素"))
	}

	// 处理HTML实体和特殊空格
	rawText := html.UnescapeString(nameSection.Text())

	// 强化文本清洗逻辑
	cleanedText := strings.NewReplacer(
		"\u00A0", " ", // Unicode非断行空格
		"&nbsp;", " ", // HTML实体空格
		"\u3000", " ", // 全角空格
	).Replace(rawText)

	// 使用正则精确匹配
	re := regexp.MustCompile(`^([^\s]+)`) // 匹配首个非空字段
	matches := re.FindStringSubmatch(cleanedText)

	if len(matches) < 2 {
		return "", errorx.WrapError(errorx.CCNUSERVER_ERROR, errors.New("姓名格式异常"))
	}

	return matches[1], nil
}

// 2.xkLoginClient 教务系统模拟登录
func (l *LoginLogic) xkLoginClient(client *http.Client) (*http.Client, error) {

	//华师本科生院登陆
	request, err := http.NewRequest("GET", "https://account.ccnu.edu.cn/cas/login?service=http%3A%2F%2Fxk.ccnu.edu.cn%2Fsso%2Fpziotlogin", nil)
	if err != nil {
		return nil, errorx.WrapError(errorx.CCNUSERVER_ERROR, err)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, errorx.WrapError(errorx.CCNUSERVER_ERROR, err)
	}
	defer resp.Body.Close()

	return client, nil
}

// 1.登陆ccnu通行证
func (l *LoginLogic) loginClient(ctx context.Context, client *http.Client, studentId string, password string, params *accountRequestParams) (*http.Client, error) {

	v := url.Values{}
	v.Set("username", studentId)
	v.Set("password", password)
	v.Set("lt", params.lt)
	v.Set("execution", params.execution)
	v.Set("_eventId", params._eventId)
	v.Set("submit", params.submit)

	request, _ := http.NewRequest("POST", "https://account.ccnu.edu.cn/cas/login;jsessionid="+params.JSESSIONID, strings.NewReader(v.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

	if len(resp.Header.Get("Set-Cookie")) == 0 {
		return nil, errorx.WrapError(errorx.PasswordErr, errors.New("学号或密码错误"))
	}

	return client, nil
}

type accountRequestParams struct {
	lt         string
	execution  string
	_eventId   string
	submit     string
	JSESSIONID string
}

// 0.前置请求,从html中提取相关参数
func (l *LoginLogic) makeAccountPreflightRequest(client *http.Client) (*accountRequestParams, error) {
	var JSESSIONID string
	var lt string
	var execution string
	var _eventId string

	params := &accountRequestParams{}

	// 初始化 http request
	request, err := http.NewRequest("GET", "https://account.ccnu.edu.cn/cas/login", nil)
	if err != nil {
		return params, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	// 发起请求
	resp, err := client.Do(request)
	if err != nil {
		return params, err
	}
	defer resp.Body.Close()

	// 读取 MsgContent
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return params, err
	}

	// 获取 Cookie 中的 JSESSIONID
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			JSESSIONID = cookie.Value
		}
	}

	if JSESSIONID == "" {
		return params, errorx.WrapError(errorx.SYSTEM_ERROR, errors.New("can not get JSESSIONID"))
	}

	// 正则匹配 HTML 返回的表单字段
	ltReg := regexp.MustCompile("name=\"lt\".+value=\"(.+)\"")
	executionReg := regexp.MustCompile("name=\"execution\".+value=\"(.+)\"")
	_eventIdReg := regexp.MustCompile("name=\"_eventId\".+value=\"(.+)\"")

	bodyStr := string(body)

	ltArr := ltReg.FindStringSubmatch(bodyStr)
	if len(ltArr) != 2 {
		return params, errorx.WrapError(errorx.SYSTEM_ERROR, errors.New("can not get lt"))
	}
	lt = ltArr[1]

	execArr := executionReg.FindStringSubmatch(bodyStr)
	if len(execArr) != 2 {
		return params, errorx.WrapError(errorx.SYSTEM_ERROR, errors.New("can not get execution"))
	}

	execution = execArr[1]

	_eventIdArr := _eventIdReg.FindStringSubmatch(bodyStr)
	if len(_eventIdArr) != 2 {
		return params, errorx.WrapError(errorx.SYSTEM_ERROR, errors.New("can not get _eventId"))
	}
	_eventId = _eventIdArr[1]

	params.lt = lt
	params.execution = execution
	params._eventId = _eventId
	params.submit = "LOGIN"
	params.JSESSIONID = JSESSIONID

	return params, nil
}
