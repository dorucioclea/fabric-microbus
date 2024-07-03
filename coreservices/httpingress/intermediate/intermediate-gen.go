/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by Microbus. DO NOT EDIT.

/*
Package intermediate serves as the foundation of the http.ingress.sys microservice.

The HTTP ingress microservice relays incoming HTTP requests to the NATS bus.
*/
package intermediate

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/microbus-io/fabric/cfg"
	"github.com/microbus-io/fabric/connector"
	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/frame"
	"github.com/microbus-io/fabric/httpx"
	"github.com/microbus-io/fabric/log"
	"github.com/microbus-io/fabric/openapi"
	"github.com/microbus-io/fabric/service"
	"github.com/microbus-io/fabric/sub"

	"gopkg.in/yaml.v3"

	"github.com/microbus-io/fabric/coreservices/httpingress/resources"
	"github.com/microbus-io/fabric/coreservices/httpingress/httpingressapi"
)

var (
	_ context.Context
	_ *embed.FS
	_ *json.Decoder
	_ fmt.Stringer
	_ *http.Request
	_ filepath.WalkFunc
	_ strconv.NumError
	_ strings.Reader
	_ time.Duration
	_ cfg.Option
	_ *errors.TracedError
	_ frame.Frame
	_ *httpx.ResponseRecorder
	_ *log.Field
	_ *openapi.Service
	_ service.Service
	_ sub.Option
	_ yaml.Encoder
	_ httpingressapi.Client
)

// ToDo defines the interface that the microservice must implement.
// The intermediate delegates handling to this interface.
type ToDo interface {
	OnStartup(ctx context.Context) (err error)
	OnShutdown(ctx context.Context) (err error)
	OnChangedPorts(ctx context.Context) (err error)
	OnChangedAllowedOrigins(ctx context.Context) (err error)
	OnChangedPortMappings(ctx context.Context) (err error)
	OnChangedReadTimeout(ctx context.Context) (err error)
	OnChangedWriteTimeout(ctx context.Context) (err error)
	OnChangedReadHeaderTimeout(ctx context.Context) (err error)
	OnChangedServerLanguages(ctx context.Context) (err error)
	OnChangedBlockedPaths(ctx context.Context) (err error)
}

// Intermediate extends and customizes the generic base connector.
// Code generated microservices then extend the intermediate.
type Intermediate struct {
	*connector.Connector
	impl ToDo
}

// NewService creates a new intermediate service.
func NewService(impl ToDo, version int) *Intermediate {
	svc := &Intermediate{
		Connector: connector.New("http.ingress.sys"),
		impl: impl,
	}
	svc.SetVersion(version)
	svc.SetDescription(`The HTTP ingress microservice relays incoming HTTP requests to the NATS bus.`)
	
	// Lifecycle
	svc.SetOnStartup(svc.impl.OnStartup)
	svc.SetOnShutdown(svc.impl.OnShutdown)

	// Configs
	svc.SetOnConfigChanged(svc.doOnConfigChanged)
	svc.DefineConfig(
		"TimeBudget",
		cfg.Description(`TimeBudget specifies the timeout for handling a request, after it has been read.
A value of 0 or less indicates no time budget.`),
		cfg.Validation(`dur [0s,]`),
		cfg.DefaultValue(`20s`),
	)
	svc.DefineConfig(
		"Ports",
		cfg.Description(`Ports is a comma-separated list of HTTP ports on which to listen for requests.`),
		cfg.DefaultValue(`8080`),
	)
	svc.DefineConfig(
		"RequestMemoryLimit",
		cfg.Description(`RequestMemoryLimit is the memory capacity used to hold pending requests, in megabytes.`),
		cfg.Validation(`int [1,]`),
		cfg.DefaultValue(`4096`),
	)
	svc.DefineConfig(
		"AllowedOrigins",
		cfg.Description(`AllowedOrigins is a comma-separated list of CORS origins to allow requests from.
The * origin can be used to allow CORS request from all origins.`),
		cfg.DefaultValue(`*`),
	)
	svc.DefineConfig(
		"PortMappings",
		cfg.Description(`PortMappings is a comma-separated list of mappings in the form x:y->z where x is the inbound
HTTP port, y is the requested NATS port, and z is the port to serve.
An HTTP request https://ingresshost:x/servicehost:y/path is mapped to internal NATS
request https://servicehost:z/path .
Both x and y can be * to indicate all ports. Setting z to * indicates to serve the requested
port y without change. Specific rules take precedence over * rules.
The default mapping grants access to all internal ports via HTTP port 8080 but restricts
HTTP ports 443 and 80 to only internal port 443.`),
		cfg.DefaultValue(`8080:*->*, 443:*->443, 80:*->443`),
	)
	svc.DefineConfig(
		"ReadTimeout",
		cfg.Description(`ReadTimeout specifies the timeout for fully reading a request.`),
		cfg.Validation(`dur [1s,]`),
		cfg.DefaultValue(`5m`),
	)
	svc.DefineConfig(
		"WriteTimeout",
		cfg.Description(`WriteTimeout specifies the timeout for fully writing the response to a request.`),
		cfg.Validation(`dur [1s,]`),
		cfg.DefaultValue(`5m`),
	)
	svc.DefineConfig(
		"ReadHeaderTimeout",
		cfg.Description(`ReadHeaderTimeout specifies the timeout for fully reading the header of a request.`),
		cfg.Validation(`dur [1s,]`),
		cfg.DefaultValue(`20s`),
	)
	svc.DefineConfig(
		"ServerLanguages",
		cfg.Description(`ServerLanguages is a comma-separated list of languages that the server supports.
This list is matched against the Accept-Language header of the request.`),
	)
	svc.DefineConfig(
		"BlockedPaths",
		cfg.Description(`A newline-separated list of paths or extensions to block with a 404.
Paths should not include any arguments and are matched exactly.
Extensions are specified with "*.ext" and are matched against the extension of the path only.`),
		cfg.DefaultValue(`/geoserver
/console/
/.env
/.amazon_aws
/solr/admin/info/system
/remote/login
/Autodiscover/Autodiscover.xml
/autodiscover/autodiscover.json
/api/v2/static/not.found
/api/sonicos/tfa
/_ignition/execute-solution
/admin.html
/auth.html
/auth1.html
/readme.txt
/__Additional
/Portal0000.htm
/docs/cplugError.html/
/CSS/Miniweb.css
/scripts/WPnBr.dll
/.git/config
/cgi-bin/.%2e/.%2e/.%2e/.%2e/bin/sh
/cgi-bin/luci/;stok=/locale
/actuator/gateway/routes
/actuator/health
/Public/home/js/check.js
/mifs/.;/services/LogService
/dns-query
/ecp/Current/exporttool/microsoft.exchange.ediscovery.exporttool.application
/owa/auth/x.js
/static/admin/javascript/hetong.js
/.git/HEAD
/sslvpnLogin.html
/vpn/index.html
/wsman
/geoserver/web
/remote/logincheck
/epa/scripts/win/nsepa_setup.exe
/.well-known/security.txt
/cf_scripts/scripts/ajax/ckeditor/ckeditor.js
/Temporary_Listen_Addresses/
/manager/html
/logon/LogonPoint/custom.html
/logon/LogonPoint/index.html
/catalog-portal/ui/oauth/verify
/error_log/.git/HEAD
*.cfm
*.asp
*.aspx
*.cgi
*.jsa
*.jsp
*.shtml
*.php
*.jhtml
*.mwsl
*.dll
*.esp
*.exe`),
	)

	// OpenAPI
	svc.Subscribe("GET", `:0/openapi.json`, svc.doOpenAPI)

	// Resources file system
	svc.SetResFS(resources.FS)

	return svc
}

// doOpenAPI renders the OpenAPI document of the microservice.
func (svc *Intermediate) doOpenAPI(w http.ResponseWriter, r *http.Request) error {
	oapiSvc := openapi.Service{
		ServiceName: svc.Hostname(),
		Description: svc.Description(),
		Version:     svc.Version(),
		Endpoints:   []*openapi.Endpoint{},
		RemoteURI:   frame.Of(r).XForwardedFullURL(),
	}

	if len(oapiSvc.Endpoints) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(&oapiSvc)
	return errors.Trace(err)
}

// doOnConfigChanged is called when the config of the microservice changes.
func (svc *Intermediate) doOnConfigChanged(ctx context.Context, changed func(string) bool) (err error) {
	if changed("Ports") {
		err := svc.impl.OnChangedPorts(ctx)
		if err != nil {
			return err // No trace
		}
	}
	if changed("AllowedOrigins") {
		err := svc.impl.OnChangedAllowedOrigins(ctx)
		if err != nil {
			return err // No trace
		}
	}
	if changed("PortMappings") {
		err := svc.impl.OnChangedPortMappings(ctx)
		if err != nil {
			return err // No trace
		}
	}
	if changed("ReadTimeout") {
		err := svc.impl.OnChangedReadTimeout(ctx)
		if err != nil {
			return err // No trace
		}
	}
	if changed("WriteTimeout") {
		err := svc.impl.OnChangedWriteTimeout(ctx)
		if err != nil {
			return err // No trace
		}
	}
	if changed("ReadHeaderTimeout") {
		err := svc.impl.OnChangedReadHeaderTimeout(ctx)
		if err != nil {
			return err // No trace
		}
	}
	if changed("ServerLanguages") {
		err := svc.impl.OnChangedServerLanguages(ctx)
		if err != nil {
			return err // No trace
		}
	}
	if changed("BlockedPaths") {
		err := svc.impl.OnChangedBlockedPaths(ctx)
		if err != nil {
			return err // No trace
		}
	}
	return nil
}

/*
TimeBudget specifies the timeout for handling a request, after it has been read.
A value of 0 or less indicates no time budget.
*/
func (svc *Intermediate) TimeBudget() (budget time.Duration) {
	_val := svc.Config("TimeBudget")
	_dur, _ := time.ParseDuration(_val)
	return _dur
}

/*
SetTimeBudget sets the value of the configuration property.

TimeBudget specifies the timeout for handling a request, after it has been read.
A value of 0 or less indicates no time budget.
*/
func (svc *Intermediate) SetTimeBudget(budget time.Duration) error {
	return svc.SetConfig("TimeBudget", fmt.Sprintf("%v", budget))
}

/*
Ports is a comma-separated list of HTTP ports on which to listen for requests.
*/
func (svc *Intermediate) Ports() (port string) {
	_val := svc.Config("Ports")
	return _val
}

/*
SetPorts sets the value of the configuration property.

Ports is a comma-separated list of HTTP ports on which to listen for requests.
*/
func (svc *Intermediate) SetPorts(port string) error {
	return svc.SetConfig("Ports", fmt.Sprintf("%v", port))
}

/*
RequestMemoryLimit is the memory capacity used to hold pending requests, in megabytes.
*/
func (svc *Intermediate) RequestMemoryLimit() (megaBytes int) {
	_val := svc.Config("RequestMemoryLimit")
	_i, _ := strconv.ParseInt(_val, 10, 64)
	return int(_i)
}

/*
SetRequestMemoryLimit sets the value of the configuration property.

RequestMemoryLimit is the memory capacity used to hold pending requests, in megabytes.
*/
func (svc *Intermediate) SetRequestMemoryLimit(megaBytes int) error {
	return svc.SetConfig("RequestMemoryLimit", fmt.Sprintf("%v", megaBytes))
}

/*
AllowedOrigins is a comma-separated list of CORS origins to allow requests from.
The * origin can be used to allow CORS request from all origins.
*/
func (svc *Intermediate) AllowedOrigins() (origins string) {
	_val := svc.Config("AllowedOrigins")
	return _val
}

/*
SetAllowedOrigins sets the value of the configuration property.

AllowedOrigins is a comma-separated list of CORS origins to allow requests from.
The * origin can be used to allow CORS request from all origins.
*/
func (svc *Intermediate) SetAllowedOrigins(origins string) error {
	return svc.SetConfig("AllowedOrigins", fmt.Sprintf("%v", origins))
}

/*
PortMappings is a comma-separated list of mappings in the form x:y->z where x is the inbound
HTTP port, y is the requested NATS port, and z is the port to serve.
An HTTP request https://ingresshost:x/servicehost:y/path is mapped to internal NATS
request https://servicehost:z/path .
Both x and y can be * to indicate all ports. Setting z to * indicates to serve the requested
port y without change. Specific rules take precedence over * rules.
The default mapping grants access to all internal ports via HTTP port 8080 but restricts
HTTP ports 443 and 80 to only internal port 443.
*/
func (svc *Intermediate) PortMappings() (mappings string) {
	_val := svc.Config("PortMappings")
	return _val
}

/*
SetPortMappings sets the value of the configuration property.

PortMappings is a comma-separated list of mappings in the form x:y->z where x is the inbound
HTTP port, y is the requested NATS port, and z is the port to serve.
An HTTP request https://ingresshost:x/servicehost:y/path is mapped to internal NATS
request https://servicehost:z/path .
Both x and y can be * to indicate all ports. Setting z to * indicates to serve the requested
port y without change. Specific rules take precedence over * rules.
The default mapping grants access to all internal ports via HTTP port 8080 but restricts
HTTP ports 443 and 80 to only internal port 443.
*/
func (svc *Intermediate) SetPortMappings(mappings string) error {
	return svc.SetConfig("PortMappings", fmt.Sprintf("%v", mappings))
}

/*
ReadTimeout specifies the timeout for fully reading a request.
*/
func (svc *Intermediate) ReadTimeout() (timeout time.Duration) {
	_val := svc.Config("ReadTimeout")
	_dur, _ := time.ParseDuration(_val)
	return _dur
}

/*
SetReadTimeout sets the value of the configuration property.

ReadTimeout specifies the timeout for fully reading a request.
*/
func (svc *Intermediate) SetReadTimeout(timeout time.Duration) error {
	return svc.SetConfig("ReadTimeout", fmt.Sprintf("%v", timeout))
}

/*
WriteTimeout specifies the timeout for fully writing the response to a request.
*/
func (svc *Intermediate) WriteTimeout() (timeout time.Duration) {
	_val := svc.Config("WriteTimeout")
	_dur, _ := time.ParseDuration(_val)
	return _dur
}

/*
SetWriteTimeout sets the value of the configuration property.

WriteTimeout specifies the timeout for fully writing the response to a request.
*/
func (svc *Intermediate) SetWriteTimeout(timeout time.Duration) error {
	return svc.SetConfig("WriteTimeout", fmt.Sprintf("%v", timeout))
}

/*
ReadHeaderTimeout specifies the timeout for fully reading the header of a request.
*/
func (svc *Intermediate) ReadHeaderTimeout() (timeout time.Duration) {
	_val := svc.Config("ReadHeaderTimeout")
	_dur, _ := time.ParseDuration(_val)
	return _dur
}

/*
SetReadHeaderTimeout sets the value of the configuration property.

ReadHeaderTimeout specifies the timeout for fully reading the header of a request.
*/
func (svc *Intermediate) SetReadHeaderTimeout(timeout time.Duration) error {
	return svc.SetConfig("ReadHeaderTimeout", fmt.Sprintf("%v", timeout))
}

/*
ServerLanguages is a comma-separated list of languages that the server supports.
This list is matched against the Accept-Language header of the request.
*/
func (svc *Intermediate) ServerLanguages() (languages string) {
	_val := svc.Config("ServerLanguages")
	return _val
}

/*
SetServerLanguages sets the value of the configuration property.

ServerLanguages is a comma-separated list of languages that the server supports.
This list is matched against the Accept-Language header of the request.
*/
func (svc *Intermediate) SetServerLanguages(languages string) error {
	return svc.SetConfig("ServerLanguages", fmt.Sprintf("%v", languages))
}

/*
A newline-separated list of paths or extensions to block with a 404.
Paths should not include any arguments and are matched exactly.
Extensions are specified with "*.ext" and are matched against the extension of the path only.
*/
func (svc *Intermediate) BlockedPaths() (blockedPaths string) {
	_val := svc.Config("BlockedPaths")
	return _val
}

/*
SetBlockedPaths sets the value of the configuration property.

A newline-separated list of paths or extensions to block with a 404.
Paths should not include any arguments and are matched exactly.
Extensions are specified with "*.ext" and are matched against the extension of the path only.
*/
func (svc *Intermediate) SetBlockedPaths(blockedPaths string) error {
	return svc.SetConfig("BlockedPaths", fmt.Sprintf("%v", blockedPaths))
}
