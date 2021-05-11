package cweb

import (
	"cx-mp-gateway/crisk"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/cerr"
	"github.com/waittttting/cRPC-common/tcp"
	"net/http"
)


// 风控
func (ws *WebServer) RiskCheck(c *gin.Context) {
	if !crisk.GRS.CheckRisk(c) {
		// 风控校验失败
		c.JSON(http.StatusOK, cerr.ErrBusy)
	}
}

// token 检查
func (ws *WebServer) AccessCheck(c *gin.Context) {

	header, err := getHeaderFromCtx(c)
	if err != nil {
		c.JSON(http.StatusOK, cerr.ErrHttpHeaderErr)
		return
	}
	switch header.Caller {
	case CallerTypeGeneralUsers: // 普通用户调用，检查 token 和 uid
		 if !CheckGeneralUserAccess(header) {
			c.JSON(http.StatusOK, cerr.ErrHttpHeaderErr)
			return
		 }
		break
	case CallerTypeServices: // 服务调用检查
		// todo: 相关检查
		break
	}
	c.Set("requestHeader", header)
}

// 转发
func (ws *WebServer) Transfer(c *gin.Context) {

	tmpHeader, _ := c.Get("requestHeader")
	header := tmpHeader.(*RequestHeader)
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, cerr.ErrHttpHeaderErr)
		logrus.Errorf("get payload err: [%v]", err)
		return
	}
	msg := &tcp.Message{
		Header: tcp.Header{
			Uid: header.Uid,
			MsgCodeType: tcp.MsgCodeTypeBus,
			ServerName: header.ServerName,
			ServerMethod: header.MethodName,
			ServerVersion: header.CommandVersion,
			PayloadLen: uint16(len(payload)),
		},
		Payload: payload,
	}
	response, err := ws.rpcClient.Send(msg)
	if err != nil {
		c.JSON(http.StatusOK, cerr.ErrHttpHeaderErr)
		logrus.Errorf("send to [%s:%s] err: [%v]", msg.Header.ServerName, msg.Header.ServerMethod, err)
		return
	}
	// todo: handle response
	println(response)
}

var unnecessaryService = map[string]bool {
	"User.register1" : true, 	// 注册
	"User.login" : true, 		// 登录
	"User.refreshToken" : true, // 刷新 accessToken
}

const (
	secretKey = "243223ffslsfsldfl412fdsfsdf" //私钥
)

//自定义Claims
type CustomClaims struct {

	UserId string
	jwt.StandardClaims
}

/**
 * @Description: 检查 token 是否合法
 * @param header
 * @return bool true 合法 false 不合法
 */
func CheckGeneralUserAccess(header *RequestHeader) bool {
	// 判断是否需要 check token
	curServiceAndCmd := fmt.Sprintf("%s.%s", header.ServerName, header.MethodName)
	if _, find := unnecessaryService[curServiceAndCmd]; find {
		return true
	}

	// todo: user 服务
	////生成token
	//customClaims := &CustomClaims{
	//	UserId : header.Uid, //用户id
	//	StandardClaims : jwt.StandardClaims{
	//		ExpiresAt : time.Now().Add(60 * 24 * 30 * time.Second).Unix(), // 过期时间，必须设置
	//	},
	//}
	//
	////采用HMAC SHA256加密算法
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	//tokenString, err := token.SignedString([]byte(secretKey))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("token: %v\n", tokenString)

	ret, err := ParseToken(header.AccessToken)
	if err != nil {
		fmt.Println(err)
	}
	// todo: 校验逻辑
	if ret.UserId == "0" {
		return true
	}
	return false
}

/**
 * @Description: 解析 token
 * @param tokenString
 * @return *CustomClaims
 * @return error
 */
func ParseToken(tokenString string) (*CustomClaims, error)  {

	println(len(tokenString))
	println(len([]byte(tokenString)))
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("err %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

