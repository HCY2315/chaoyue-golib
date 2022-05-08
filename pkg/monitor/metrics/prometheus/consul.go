/*
 * @File   : consulsd
 * @Author : huangbin
 *
 * @Created on 2020/8/23 9:52 上午
 * @Project : prometheus
 * @Software: GoLand
 * @Description  : 一个实例表示注册一个服务
 */

package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	CSTConsulRegisterUrl     = "/v1/agent/service/register"
	CSTConsulDeRegisterUrl   = "/v1/agent/service/deregister/"
	CSTConsulQueryServiceUrl = "/v1/catalog/service/"
)

type consulRegistry struct {
	node Node
	addr string
}

/*
 * @Method   : NewConsulRegistry/新建consul句柄
 *
 * @param    : Node 需要注册的节点信息
 * @param    : consulAddr string consul地址【http://10.155.19.154:8500】
 * @Return   : *consulRegistry consulserver句柄
 * @Description :
 */
func NewConsulRegistry(node Node, consulAddr string) *consulRegistry {
	return &consulRegistry{
		node: node,
		addr: consulAddr,
	}
}

/*
 * @Method   : RegisterService/注册服务
 *
 * @param    : none
 * @Return   : []byte 服务返回信息
 *             error  错误
 * @Description :
 */
func (c *consulRegistry) RegisterService() ([]byte, error) {
	addr := c.addr + CSTConsulRegisterUrl
	bs, err := json.Marshal(c.node)
	if err != nil {
		return nil, fmt.Errorf("[consulsd] json.Marshal(Node): %s", err.Error())
	}

	_, resp, err := httpRequest(http.MethodPut, addr, bs, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("[consulsd] http put register service[%#v] : %v", c.node, err)
	}

	return resp, nil
}

/*
 * @Method   : DelRegisterService/删除服务
 *
 * @param    : none
 * @Return   : []byte 服务返回信息
 *             error  错误
 * @Description :
 */
func (c *consulRegistry) DelRegisterService() ([]byte, error) {
	addr := c.addr + CSTConsulDeRegisterUrl + c.node.Id
	_, resp, err := httpRequest(http.MethodPut, addr, nil, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("[consulsd] http put delete service[%v] : %v", c.node.Id, err)
	}

	return resp, nil
}

/*
 * @Method   : QueryService/查询服务
 *
 * @param    : serviceID string 查询的服务id
 * @param    : consulAddr string 需要查询的consul地址【http://x.x.x.x:8500】
 * @Return   : []byte 服务返回信息
 *             error  错误
 * @Description :
 */
func (c *consulRegistry) QueryService(serviceID, consulAddr string) ([]byte, error) {
	addr := consulAddr + CSTConsulQueryServiceUrl + serviceID
	_, resp, err := httpRequest(http.MethodPut, addr, nil, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("[consulsd] http put query service[%v] : %v", serviceID, err)
	}

	return resp, nil
}

/*
 * @Method   : Service/获取注册的服务信息
 *
 * @param    : none
 * @Return   : Node 服务节点信息
 * @Description :
 */
func (c *consulRegistry) Service() Node {
	return c.node
}

func httpRequest(method, url string, data []byte, headers map[string]string) (int, []byte, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return -1, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("ioutil.ReadAll(resp.Body): %v", err)
	}

	return resp.StatusCode, respBytes, nil
}
