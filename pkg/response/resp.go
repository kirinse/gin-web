package response

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

// http请求响应封装
type Resp struct {
	Code int         `json:"code"` // 错误代码代码
	Data interface{} `json:"data"` // 数据内容
	Msg  string      `json:"msg"`  // 消息提示
}

// 分页封装
type PageInfo struct {
	PageNum  uint `json:"page_num" form:"page_num"`   // 当前页码
	PageSize uint `json:"page_size" form:"page_size"` // 每页显示条数
	Total    uint `json:"total"`                      // 数据总条数
}

// 带分页数据封装
type PageData struct {
	PageInfo
	List interface{} `json:"list"` // 数据列表
}

// 计算limit/offset, 如果需要用到返回的PageSize, PageNum, 务必保证Total值有效
func (s *PageInfo) GetLimit() (limit uint, offset uint) {
	// 传入参数可能不合法, 设置默认值
	// 每页显示条数不能小于1
	if s.PageSize < 1 {
		s.PageSize = 10
	}
	// 页码不能小于1
	if s.PageNum < 1 {
		s.PageNum = 1
	}

	limit = s.PageSize
	offset = limit * (s.PageNum - 1)

	if s.Total > 0 {
		// 如果偏移量比总条数还多
		if limit > s.Total {
			limit = s.Total
		}
		if offset > s.Total {
			offset = s.Total
		}
		if offset+limit > s.Total {
			if limit > s.Total {
				offset = limit - s.Total
			} else {
				offset = s.Total - limit
			}
		}
	}
	// 恢复真实的页码和大小
	if limit > 0 {
		s.PageSize = limit
		s.PageNum = uint(math.Ceil(float64(offset)/float64(limit))) + 1
	}
	return
}

const (
	SUCCESS   = 201
	FAIL      = 405
	EXCEPTION = 500
)

var errorMsg = map[int]string{
	SUCCESS:   "操作成功",
	FAIL:      "操作失败",
	EXCEPTION: "系统异常",
}

func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Success(c *gin.Context) {
	Result(c, SUCCESS, errorMsg[SUCCESS], map[string]interface{}{})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, errorMsg[SUCCESS], data)
}

func SuccessWithMsg(c *gin.Context, msg string) {
	Result(c, SUCCESS, msg, map[string]interface{}{})
}

func Fail(c *gin.Context) {
	Result(c, FAIL, errorMsg[FAIL], map[string]interface{}{})
}

func FailWithMsg(c *gin.Context, msg string) {
	Result(c, FAIL, msg, map[string]interface{}{})
}

func FailWithCode(c *gin.Context, code int) {
	Result(c, code, errorMsg[code], map[string]interface{}{})
}