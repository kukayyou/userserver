package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kataras/iris"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/token"
	"io/ioutil"
	"time"
)

//错误码
const (
	PARAMS_PARSE_ERROR      = 1001 + iota //请求参数解析错误
	TOKEN_CHECK_ERROR       = 1001 + iota //token验证错误
	USER_CHECK_ERROR        = 1001 + iota //用户验证错误，非本人操作
	USER_REGISTER_ERROR     = 1001 + iota //注册错误
	USER_LOGIN_ERROR        = 1001 + iota //登录错误
	USER_GET_INFOS_ERROR    = 1001 + iota //获取用户信息错误
	USER_UPDATE_INFOS_ERROR = 1001 + iota //更新用户信息错误
	DEMAND_CREATE_ERROR     = 1001 + iota //创建需求错误
	DEMAND_UPDATE_ERROR     = 1001 + iota //更新需求错误
	DEMAND_QUERY_ERROR      = 1001 + iota //查询需求错误
	DEMAND_DELETE_ERROR     = 1001 + iota //删除需求错误
	SKILL_CREATE_ERROR      = 1001 + iota //创建需求错误
	SKILL_UPDATE_ERROR      = 1001 + iota //更新需求错误
	SKILL_QUERY_ERROR       = 1001 + iota //查询需求错误
	SKILL_DELETE_ERROR      = 1001 + iota //删除需求错误
)

type BaseController struct {
	mylog.LogInfo
	ReqParams   []byte
	ServerToken string
	StartTime   time.Time
	Resp        Response
	ReqContext  context.Context
}

type Response struct {
	Code      int64       `json:"code"`      //错误码
	Msg       string      `json:"msg"`       //错误信息
	RequestID string      `json:"requestId"` //请求id
	CostTime  string      `json:"costTime"`  //请求耗时
	Data      interface{} `json:"data"`      //返回数据
}

func (bc *BaseController) Prepare(ctx *gin.Context) {
	bc.StartTime = time.Now()
	//设置requestid
	bc.SetRequestId()
	bc.ReqContext = context.WithValue(context.TODO(), "requestID", bc.GetRequestId())
	//设置请求url
	bc.SetRequestUrl(ctx.Request.RequestURI)
	bc.ReqContext = context.WithValue(bc.ReqContext, "requestUrl", bc.GetRequestUrl())
	//设置返回requestid
	bc.Resp.RequestID = bc.GetRequestId()
	//获取请求参数
	bc.ReqParams, _ = ioutil.ReadAll(ctx.Request.Body)

	mylog.Info("requestId:%s, requestUrl:%s, params : %s", bc.GetRequestId(), bc.GetRequestUrl(), string(bc.ReqParams))
}

func (bc *BaseController) PrepareIris(ctx iris.Context) {
	//执行开始时间
	bc.StartTime = time.Now()
	//设置requestid
	bc.SetRequestId()
	bc.ReqContext = context.WithValue(context.TODO(), "requestID", bc.GetRequestId())
	//设置请求url
	bc.SetRequestUrl(ctx.Request().RequestURI)
	bc.ReqContext = context.WithValue(bc.ReqContext, "requestUrl", bc.GetRequestUrl())
	//设置返回requestid
	bc.Resp.RequestID = bc.GetRequestId()
	//获取请求参数
	bc.ReqParams, _ = ioutil.ReadAll(ctx.Request().Body)

	mylog.SugarLogger.Info(fmt.Sprintf("requestId:%s, requestUrl:%s, params : %s", bc.GetRequestId(), bc.GetRequestUrl(), string(bc.ReqParams)))
}

func (bc *BaseController) FinishResponse(ctx *gin.Context) {
	//执行结束时间
	endTime := time.Now()
	bc.Resp.CostTime = fmt.Sprintf("%4v", endTime.Sub(bc.StartTime))
	if len(bc.Resp.Msg) <= 0 {
		bc.Resp.Msg = "success"
	}
	ctx.JSON(200,
		gin.H{
			"errcode":   bc.Resp.Code,
			"errmsg":    bc.Resp.Msg,
			"requestId": bc.Resp.RequestID,
			"costTime":  bc.Resp.CostTime,
			"data":      bc.Resp.Data,
		})
	r, _ := json.Marshal(bc.Resp)
	mylog.Info("requestUrl:%s, response data:%s", bc.GetRequestUrl(), string(r))
}

func (bc *BaseController) FinishResponseIris(ctx iris.Context) {
	//执行结束时间
	endTime := time.Now()
	bc.Resp.CostTime = fmt.Sprintf("%4v", endTime.Sub(bc.StartTime))
	if len(bc.Resp.Msg) <= 0 {
		bc.Resp.Msg = "success"
	}

	_, err := ctx.JSON(iris.Map{
		"errcode":   bc.Resp.Code,
		"errmsg":    bc.Resp.Msg,
		"requestId": bc.Resp.RequestID,
		"costTime":  bc.Resp.CostTime,
		"data":      bc.Resp.Data,
	})

	if err != nil {
		mylog.SugarLogger.Error(fmt.Sprintf("requestId:%s, requestUrl:%s, response data err:%s", bc.Resp.RequestID, bc.GetRequestUrl(), err.Error()))
	}

	r, _ := json.Marshal(bc.Resp)
	mylog.SugarLogger.Info(fmt.Sprintf("requestUrl:%s, response data:%s", bc.GetRequestUrl(), string(r)))
}

func (bc *BaseController) CheckToken(userID, tokenData string) (err error) {
	if len(bc.ServerToken) == 0 {
		err = bc.userCheck(userID, tokenData)
	} else {
		err = bc.serverCheck()
	}
	return
}

func (bc *BaseController) userCheck(userID, tokenData string) error {
	if claim, err := token.CheckUserToken(tokenData); err != nil {
		bc.Resp.Code = TOKEN_CHECK_ERROR
		bc.Resp.Msg = "token check failed!"
		return fmt.Errorf("token check failed!")
	} else if claim.UserData.UserID != userID {
		bc.Resp.Code = USER_CHECK_ERROR
		bc.Resp.Msg = "user is invilid!"
		return fmt.Errorf("user is invalid!")
	}
	return nil
}

func (bc *BaseController) serverCheck() error {
	if _, err := token.CheckServerToken(bc.ServerToken); err != nil {
		bc.Resp.Code = TOKEN_CHECK_ERROR
		bc.Resp.Msg = "token check failed!"
		return fmt.Errorf("token check failed!")
	}
	return nil
}
