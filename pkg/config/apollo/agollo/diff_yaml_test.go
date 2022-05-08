package agollo

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	yamlCfg1 = `log:
    file_name: ../log/apigateway.log
    console: false
    level: debug
    max_file_size: 512


hystrix: # 熔断配置
    ask_app_server: # 访问appserver的熔断参数
        enable: true
        timeout: 1000 # 超时时间，单位: ms
        max_concurrent_requests: 50 # hystrix同一时刻允许执行的任务数，qps*平均延迟
        request_volume_threshold: 100
        sleep_window: 100 # 熔断持续时间, 单位: ms
        error_percent_threshold: 25 # 触发熔断的失败率阈值`

	yamlCfg2 = `log:
    file_name: ../log/apigateway.log
    console: false
    level: debug
    max_file_size: 512`

	yamlCfg3 = `log:
    file_name: ../log/apigateway.log
    console: true
    level: debug
    max_file_size: 512


hystrix: # 熔断配置
    ask_app_server: # 访问appserver的熔断参数
        enable: false
        timeout: 1000 # 超时时间，单位: ms
        max_concurrent_requests: 50 # hystrix同一时刻允许执行的任务数，qps*平均延迟
        request_volume_threshold: 100
        sleep_window: 100 # 熔断持续时间, 单位: ms
        error_percent_threshold: 25 # 触发熔断的失败率阈值`
)

func TestProcessYamlDiffUnderCertainKeys(t *testing.T) {
	var observerCalled int32 = 0
	var wg *sync.WaitGroup
	type args struct {
		observersDelegate map[string]Observer
		namespace         string
		oldYamlStr        string
		newYamlStr        string
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		observerExpectCalled int32
		setup                func(t *testing.T) func(t *testing.T)
	}{
		{
			name: "same config",
			args: args{
				observersDelegate: map[string]Observer{
					"config.yaml:log": func(key string, oldValue interface{}, newValue interface{}) {
						assert.Fail(t, "should not run here")
						assert.Equal(t, "log", key)
						observerCalled++
						fmt.Printf("same config++\n")
						wg.Done()
					},
					"config.yaml:hystrix": func(key string, oldValue interface{}, newValue interface{}) {
						assert.Fail(t, "should not run here")
						observerCalled++
						fmt.Printf("same config++\n")
						wg.Done()
					},
				},
				namespace:  "config.yaml",
				oldYamlStr: yamlCfg1,
				newYamlStr: yamlCfg1,
			},
			wantErr: false,
			setup: func(t *testing.T) func(t *testing.T) {
				fmt.Printf("test for same config\n")
				observerCalled = 0
				wg = &sync.WaitGroup{}
				return func(t *testing.T) {}
			},
			observerExpectCalled: 0,
		},
		{
			name: "same config item",
			args: args{
				observersDelegate: map[string]Observer{
					"config.yaml:log.level": func(key string, oldValue interface{}, newValue interface{}) {
						assert.Fail(t, "should not run here")
						assert.Equal(t, "log.level", key)
						observerCalled++
						wg.Done()
					},
					"config.yaml:hystrix.max_concurrent_requests": func(key string, oldValue interface{}, newValue interface{}) {
						assert.Fail(t, "should not run here")
						assert.Equal(t, "hystrix.max_concurrent_requ", key)
						observerCalled++
						wg.Done()
					},
				},
				namespace:  "config.yaml",
				oldYamlStr: yamlCfg1,
				newYamlStr: yamlCfg3,
			},
			wantErr: false,
			setup: func(t *testing.T) func(t *testing.T) {
				observerCalled = 0
				fmt.Printf("test for same config item\n")
				wg = &sync.WaitGroup{}
				return func(t *testing.T) {}
			},
			observerExpectCalled: 0,
		},
		{
			name: "modify config",
			args: args{
				observersDelegate: map[string]Observer{
					"config.yaml:log": func(key string, oldValue interface{}, newValue interface{}) {
						assert.False(t, reflect.DeepEqual(oldValue, newValue))
						assert.Equal(t, "log", key)
						atomic.AddInt32(&observerCalled, 1)
						wg.Done()
					},
					"config.yaml:log.console": func(key string, oldValue interface{}, newValue interface{}) {
						assert.False(t, reflect.DeepEqual(oldValue, newValue))
						assert.Equal(t, "log.console", key)
						atomic.AddInt32(&observerCalled, 1)
						wg.Done()
					},
				},
				namespace:  "config.yaml",
				oldYamlStr: yamlCfg1,
				newYamlStr: yamlCfg3,
			},
			wantErr: false,
			setup: func(t *testing.T) func(t *testing.T) {
				observerCalled = 0
				wg = &sync.WaitGroup{}
				wg.Add(2)
				return func(t *testing.T) {}
			},
			observerExpectCalled: 2,
		},
		{
			name: "delete some config",
			args: args{
				observersDelegate: map[string]Observer{
					"config.yaml:hystrix": func(key string, oldValue interface{}, newValue interface{}) {
						assert.NotEqual(t, nil, oldValue)
						assert.Equal(t, "hystrix", key)
						assert.Equal(t, nil, newValue)
						observerCalled++
						wg.Done()
					},
				},
				namespace:  "config.yaml",
				oldYamlStr: yamlCfg1,
				newYamlStr: yamlCfg2,
			},
			wantErr: false,
			setup: func(t *testing.T) func(t *testing.T) {
				observerCalled = 0
				wg = &sync.WaitGroup{}
				wg.Add(1)
				return func(t *testing.T) {}
			},
			observerExpectCalled: 1,
		},
		{
			name: "add some config",
			args: args{
				observersDelegate: map[string]Observer{
					"config.yaml:hystrix": func(key string, oldValue interface{}, newValue interface{}) {
						assert.Equal(t, nil, oldValue)
						assert.Equal(t, "hystrix", key)
						assert.NotEqual(t, nil, newValue)
						observerCalled++
						wg.Done()
					},
				},
				namespace:  "config.yaml",
				oldYamlStr: yamlCfg2,
				newYamlStr: yamlCfg1,
			},
			wantErr: false,
			setup: func(t *testing.T) func(t *testing.T) {
				observerCalled = 0
				wg = &sync.WaitGroup{}
				wg.Add(1)
				return func(t *testing.T) {}
			},
			observerExpectCalled: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			err := ProcessYamlDiffUnderCertainKeys(tt.args.observersDelegate, tt.args.namespace, tt.args.oldYamlStr, tt.args.newYamlStr)
			assert.Equal(t, tt.wantErr, err != nil)
			wg.Wait()
			assert.Equal(t, tt.observerExpectCalled, observerCalled)
		})
	}
}
