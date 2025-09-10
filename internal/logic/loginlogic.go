package logic

import (
	"context"
	"errors"
	"io"
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
	client := l.client()

	params, err := l.makeAccountPreflightRequest(client)
	if err != nil {
		return nil, errorx.WrapError(errorx.SYSTEM_ERROR, err)
	}

	_, err = l.loginClient(l.ctx, client, req.Username, req.Password, params)
	if err != nil {
		return nil, errorx.WrapError(errorx.SYSTEM_ERROR, err)
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
			Name:      "",
			StudentId: req.Username,
		},
	}

	return resp, nil
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
