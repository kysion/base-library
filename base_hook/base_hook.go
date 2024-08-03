package base_hook

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
	"github.com/kysion/base-library/base_model"
	"github.com/kysion/base-library/base_model/base_enum"
	"reflect"
	"strings"
	"time"
)

var wsArr = gmap.Map{}

type BaseHookModel struct {
	hookArr garray.Array
}

func (s *BaseHookModel) GetHookArr() *garray.Array {
	return &s.hookArr
}

type BaseHook[T any, F any] struct {
	BaseHookModel
}

type Option struct {
	Data any // 需要广播的数据
}

func (s *BaseHook[T, F]) GetBusinessType() base_enum.HookBusinessType {
	var t F
	return base_enum.Hook.BusinessType.New(reflect.TypeOf(t).String())
}

// InstallHook 安装Hook
func (s *BaseHook[T, F]) InstallHook(filter T, hookFunc F) {
	item := base_model.KeyValueT[T, F]{Key: filter, Value: hookFunc}

	s.hookArr.Append(item)
}

// UnInstallHook 卸载Hook
func (s *BaseHook[T, F]) UnInstallHook(filter T, f ...func(filter T, key T) bool) {
	newFuncArr := garray.NewArray()
	s.hookArr.Iterator(func(key int, value interface{}) bool {
		item := value.(base_model.KeyValueT[T, F])

		if len(f) > 0 && f[0](filter, item.Key) == false {
			newFuncArr.Append(value)
		}

		return true
	})
	s.hookArr = *newFuncArr
}

// ClearAllHook 清除Hook
func (s *BaseHook[T, F]) ClearAllHook() {
	s.hookArr.Clear()
}

// Iterator 遍历Hook
func (s *BaseHook[T, F]) Iterator(f func(key T, value F), options ...Option) {
	s.hookArr.Iterator(func(key int, value interface{}) bool {
		item := value.(base_model.KeyValueT[T, F])
		f(item.Key, item.Value)
		return true
	})

	if len(options) <= 0 || options[0].Data == nil {
		return
	}

	_ = g.Try(context.Background(), func(ctx context.Context) {
		var businessType F
		s2 := reflect.TypeOf(businessType).String()
		fmt.Println(s2)

		s.publish(options[0], base_enum.Hook.BusinessType.New(s2))
	})
}

func (s *BaseHook[T, F]) Where(filter T, f func(filter T, key T) bool) []F {
	result := make([]F, 0)

	s.hookArr.Iterator(func(key int, value interface{}) bool {
		item := value.(base_model.KeyValueT[T, F])

		if f(filter, item.Key) {
			result = append(result, item.Value)
		}

		return true
	})

	return result
}

// HookDistribution 开启一个websocket服务，用于接收广播消息
func HookDistribution(r *ghttp.Request) {
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(r.Context(), err)
		r.Exit()
	}

	if Gateway().HasHookMessage() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}

			fmt.Println(string(msg))

			data := base_model.HookModel{}
			// 解析订阅数据包
			err = gjson.DecodeTo(msg, &data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			addr := strings.Split(r.RemoteAddr, ":")
			data.Ctx = r.Context()
			data.Source = &base_model.HookHostInfo{
				Host: addr[0],
				Port: gconv.Int(addr[1]),
			}

			Gateway().BroadcastMessage(data)
		}
	}
}

func (s *BaseHook[T, F]) publish(dataInfo interface{}, businessType base_enum.HookBusinessType) {
	/*
			解决跨进城Hook订阅不了问题的方案：
				1、获取配置的服务注册表
				2、按照配置服务，通过ws协议连接到对应的服务 (TODO 重连检测机制)
				3、发送消息给对应的服务
		       	4、读取服务响应的消息 （可选）
	*/

	// 1、获取配置的服务注册表
	serviceArr := garray.NewStrArrayFrom(g.Cfg().MustGet(context.Background(), "service.hostAddressArr").Strings())

	data := base_model.HookModel{
		BusinessTypeStr: businessType.Code(),
		Data:            dataInfo, // 发送的数据,
	}

	// 2、按照配置服务，通过ws协议连接到对应的服务
	serviceArr.Iterator(func(k int, v string) bool {
		// ws链接服务
		urlStr := "ws://" + v + "/ws"
		fmt.Println(urlStr)

		var conn *websocket.Conn
		var err error
		wsConn := wsArr.Get(urlStr)
		if wsConn == nil {
			client := gclient.NewWebSocket()
			client.HandshakeTimeout = time.Second
			client.TLSClientConfig = &tls.Config{} // 设置 tls 配置
			conn, _, err = client.Dial(urlStr, nil)
			if err != nil {
				// 如果连接失败，则调用重链接方法
				return true
			}
			// 增加连接断开等消息检测，如果连接断开则将conn对象从wsArr中删除
			go func() {
				for {
					_, _, err = conn.ReadMessage()
					if err != nil {
						wsArr.Remove(urlStr)
						break
					}
				}
			}()

			wsArr.Set(urlStr, conn)
		} else {
			conn = wsConn.(*websocket.Conn)
		}

		//defer conn.Close()

		// 4、异步读取服务响应的消息 （可选）
		//go func() {
		//	_, msg, _ := conn.ReadMessage()
		//	fmt.Println(string(msg))
		//	conn.Close()
		//}()

		// 3、发送消息给对应的服务
		err = conn.WriteJSON(data)
		if err != nil {
			panic(err)
		}

		return false
	})

}
