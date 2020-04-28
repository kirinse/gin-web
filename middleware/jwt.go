package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-shipment-api/models"
	"go-shipment-api/pkg/global"
	"go-shipment-api/pkg/request"
	"go-shipment-api/pkg/response"
	"go-shipment-api/pkg/service"
	"go-shipment-api/pkg/utils"
	"time"
)

var jwtSecret = "jwt-secret"

func InitAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test realm",
		Key:             []byte("secret key"),
		Timeout:         time.Hour * 24,
		MaxRefresh:      time.Hour * 24,
		IdentityKey:     jwtSecret,                                          // jwt密钥
		PayloadFunc:     payloadFunc,                                        // 有效载荷处理
		IdentityHandler: identityHandler,                                    // 解析Claims
		Authenticator:   login,                                              // 校验token的正确性, 处理登录逻辑
		Authorizator:    authorizator,                                       // 校验用户的正确性
		Unauthorized:    unauthorized,                                       // 校验失败处理
		LoginResponse:   loginResponse,                                      // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                     // 登出后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt", // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                           // header名称
		TimeFunc:        time.Now,
	})
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var user models.SysUser
		// 将用户json转为结构体
		utils.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.Id,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

func login(c *gin.Context) (interface{}, error) {
	var req request.RegisterAndLoginRequestStruct
	// 请求json绑定
	_ = c.ShouldBindJSON(&req)

	u := &models.SysUser{
		Username: req.Username,
		Password: req.Password,
	}

	// 密码校验
	user, err := service.LoginCheck(u)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	return map[string]interface{}{
		"user": utils.Struct2Json(user),
	}, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		var user models.SysUser
		// 将用户json转为结构体
		utils.JsonI2Struct(v["user"], &user)
		// 将用户保存到context, api调用时取数据方便
		c.Set("user", user)
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	global.Log.Debug(message)
	response.FailWithMsg(c, "jwt校验失败")
}

func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.SuccessWithData(c, map[string]interface{}{
		"token":   token,
		"expires": expires,
	})
}

func logoutResponse(c *gin.Context, code int) {
	response.Success(c)
}