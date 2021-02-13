package proxy

import (
	"comradequinn/hflow/proxy/intercept"
	"sync"
)

var (
	// TrafficSummaryWriter specifies an intercept that is responsible for writing traffic summary data
	TrafficSummaryWriter, SetTrafficSummaryWriter = func() (func() intercept.Intercept, func(intercept.Intercept)) {
		w := intercept.Unset
		mx := sync.Mutex{}
		return func() intercept.Intercept {
				mx.Lock()
				defer mx.Unlock()
				return w
			}, func(i intercept.Intercept) {
				mx.Lock()
				defer mx.Unlock()
				w = i
			}
	}()
	// CaptureFileWriter specifies an intercept that is responsible for writing traffic
	CaptureFileWriter, SetCaptureFileWriter = func() (func() intercept.Intercept, func(intercept.Intercept)) {
		w := intercept.Unset
		mx := sync.Mutex{}
		return func() intercept.Intercept {
				mx.Lock()
				defer mx.Unlock()
				return w
			}, func(i intercept.Intercept) {
				mx.Lock()
				defer mx.Unlock()
				w = i
			}
	}()
	// RequestRerouter specifies an intercept that is responsible for rerouting traffic
	RequestRerouter, SetRequestRerouter = func() (func() intercept.Intercept, func(intercept.Intercept)) {
		r := intercept.Unset
		mx := sync.Mutex{}
		return func() intercept.Intercept {
				mx.Lock()
				defer mx.Unlock()
				return r
			}, func(i intercept.Intercept) {
				mx.Lock()
				defer mx.Unlock()
				r = i
			}
	}()
	// TrafficEditor specifies an intercept that is responsible for editing traffic
	TrafficEditor, SetTrafficEditor = func() (func() intercept.Intercept, func(intercept.Intercept)) {
		e := intercept.Unset
		mx := sync.Mutex{}
		return func() intercept.Intercept {
				mx.Lock()
				defer mx.Unlock()
				return e
			}, func(i intercept.Intercept) {
				mx.Lock()
				defer mx.Unlock()
				e = i
			}
	}()
)
