// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package httpconv provides OpenTelemetry HTTP semantic conventions for
// tracing telemetry.
package httpconv // import "go.opentelemetry.io/otel/semconv/v1.24.0/httpconv"

import (
	"net/http"

	"github.com/nextmicro/gokit/trace/semconvutil"
	"github.com/nextmicro/gokit/trace/semconvutil/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	nc = &v4.NetConv{
		NetHostNameKey:     semconv.NetHostNameKey,
		NetHostPortKey:     semconv.NetHostPortKey,
		NetPeerNameKey:     semconv.NetPeerNameKey,
		NetPeerPortKey:     semconv.NetPeerPortKey,
		NetSockPeerAddrKey: semconv.NetSockPeerAddrKey,
		NetSockPeerPortKey: semconv.NetSockPeerPortKey,
		NetTransportOther:  semconv.NetTransportOther,
		NetTransportTCP:    semconv.NetTransportTCP,
		NetTransportUDP:    semconv.NetTransportUDP,
		NetTransportInProc: semconv.NetTransportInProc,
	}

	hc = &v4.HTTPConv{
		NetConv: nc,

		EnduserIDKey:                 semconv.EnduserIDKey,
		HTTPClientIPKey:              semconv.ClientAddressKey,
		NetProtocolNameKey:           semconv.NetProtocolNameKey,
		NetProtocolVersionKey:        semconv.NetProtocolVersionKey,
		HTTPMethodKey:                semconv.HTTPMethodKey,
		HTTPRequestContentLengthKey:  semconv.HTTPRequestContentLengthKey,
		HTTPResponseContentLengthKey: semconv.HTTPResponseContentLengthKey,
		HTTPRouteKey:                 semconv.HTTPRouteKey,
		HTTPSchemeHTTP:               semconv.HTTPSchemeKey.String("http"),
		HTTPSchemeHTTPS:              semconv.HTTPSchemeKey.String("https"),
		HTTPStatusCodeKey:            semconv.HTTPStatusCodeKey,
		HTTPTargetKey:                semconv.HTTPTargetKey,
		HTTPURLKey:                   semconv.HTTPURLKey,
		UserAgentOriginalKey:         semconv.UserAgentOriginalKey,
	}
)

// ClientResponse returns trace attributes for an HTTP response received by a
// client from a server. It will return the following attributes if the related
// values are defined in resp: "http.status.code",
// "http.response_content_length".
//
// This does not add all OpenTelemetry required attributes for an HTTP event,
// it assumes ClientRequest was used to create the span with a complete set of
// attributes. If a complete set of attributes can be generated using the
// request contained in resp. For example:
//
//	append(ClientResponse(resp), ClientRequest(resp.Request)...)
func ClientResponse(resp *http.Response) []attribute.KeyValue {
	return hc.ClientResponse(resp)
}

// ClientRequest returns trace attributes for an HTTP request made by a client.
// The following attributes are always returned: "http.url",
// "net.protocol.(name|version)", "http.method", "net.peer.name".
// The following attributes are returned if the related values are defined
// in req: "net.peer.port", "http.user_agent", "http.request_content_length",
// "enduser.id".
func ClientRequest(req *http.Request) []attribute.KeyValue {
	return hc.ClientRequest(req)
}

// ClientStatus returns a span status code and message for an HTTP status code
// value received by a client.
func ClientStatus(code int) (codes.Code, string) {
	return hc.ClientStatus(code)
}

// ServerRequest returns trace attributes for an HTTP request received by a
// server.
//
// The server must be the primary server name if it is known. For example this
// would be the ServerName directive
// (https://httpd.apache.org/docs/2.4/mod/core.html#servername) for an Apache
// server, and the server_name directive
// (http://nginx.org/en/docs/http/ngx_http_core_module.html#server_name) for an
// nginx server. More generically, the primary server name would be the host
// header value that matches the default virtual host of an HTTP server. It
// should include the host identifier and if a port is used to route to the
// server that port identifier should be included as an appropriate port
// suffix.
//
// If the primary server name is not known, server should be an empty string.
// The req Host will be used to determine the server instead.
//
// The following attributes are always returned: "http.method", "http.scheme",
// ""net.protocol.(name|version)", "http.target", "net.host.name".
// The following attributes are returned if they related values are defined
// in req: "net.host.port", "net.sock.peer.addr", "net.sock.peer.port",
// "user_agent.original", "enduser.id", "http.client_ip".
func ServerRequest(server string, req *http.Request) []attribute.KeyValue {
	return hc.ServerRequest(server, req)
}

// ServerStatus returns a span status code and message for an HTTP status code
// value returned by a server. Status codes in the 400-499 range are not
// returned as errors.
func ServerStatus(code int) (codes.Code, string) {
	return hc.ServerStatus(code)
}

// RequestHeader returns the contents of h as attributes.
//
// Instrumentation should require an explicit configuration of which headers to
// captured and then prune what they pass here. Including all headers can be a
// security risk - explicit configuration helps avoid leaking sensitive
// information.
//
// The User-Agent header is already captured in the user_agent.original attribute
// from ClientRequest and ServerRequest. Instrumentation may provide an option
// to capture that header here even though it is not recommended. Otherwise,
// instrumentation should filter that out of what is passed.
func RequestHeader(h http.Header) []attribute.KeyValue {
	return hc.RequestHeader(h)
}

// ResponseHeader returns the contents of h as attributes.
//
// Instrumentation should require an explicit configuration of which headers to
// captured and then prune what they pass here. Including all headers can be a
// security risk - explicit configuration helps avoid leaking sensitive
// information.
//
// The User-Agent header is already captured in the user_agent.original attribute
// from ClientRequest and ServerRequest. Instrumentation may provide an option
// to capture that header here even though it is not recommended. Otherwise,
// instrumentation should filter that out of what is passed.
func ResponseHeader(h http.Header) []attribute.KeyValue {
	return hc.ResponseHeader(h)
}

const (
	// The IP address of the original client behind all proxies, if known (e.g. from
	// [X-Forwarded-For](https://developer.mozilla.org/en-
	// US/docs/Web/HTTP/Headers/X-Forwarded-For)).
	//
	// Type: string
	// Required: No
	// Stability: stable
	// Examples: '83.164.160.102'
	// Note: This is not necessarily the same as `net.peer.ip`, which would
	// identify the network-level peer, which may be a proxy.

	// This attribute should be set when a source of information different
	// from the one used for `net.peer.ip`, is available even if that other
	// source just confirms the same value as `net.peer.ip`.
	// Rationale: For `net.peer.ip`, one typically does not know if it
	// comes from a proxy, reverse proxy, or the actual client. Setting
	// `http.client_ip` when it's the same as `net.peer.ip` means that
	// one is at least somewhat confident that the address is not that of
	// the closest proxy.
	HTTPClientIPKey = attribute.Key("http.client_ip")
	// Kind of HTTP protocol used.
	//
	// Type: Enum
	// Required: No
	// Stability: stable
	// Note: If `net.transport` is not specified, it can be assumed to be `IP.TCP`
	// except if `http.flavor` is `QUIC`, in which case `IP.UDP` is assumed.
	HTTPFlavorKey = attribute.Key("http.flavor")
	// The value of the [HTTP host
	// header](https://tools.ietf.org/html/rfc7230#section-5.4). An empty Host header
	// should also be reported, see note.
	//
	// Type: string
	// Required: No
	// Stability: stable
	// Examples: 'www.example.org'
	// Note: When the header is present but empty the attribute SHOULD be set to the
	// empty string. Note that this is a valid situation that is expected in certain
	// cases, according the aforementioned [section of RFC
	// 7230](https://tools.ietf.org/html/rfc7230#section-5.4). When the header is not
	// set the attribute MUST NOT be set.
	HTTPHostKey = attribute.Key("http.host")
	// The primary server name of the matched virtual host. This should be obtained
	// via configuration. If no such configuration can be obtained, this attribute
	// MUST NOT be set ( `net.host.name` should be used instead).
	//
	// Type: string
	// Required: No
	// Stability: stable
	// Examples: 'example.com'
	// Note: `http.url` is usually not readily available on the server side but would
	// have to be assembled in a cumbersome and sometimes lossy process from other
	// information (see e.g. open-telemetry/opentelemetry-python/pull/148). It is thus
	// preferred to supply the raw data that is available.
	HTTPServerNameKey = attribute.Key("http.server_name")
	// Value of the [HTTP User-
	// Agent](https://tools.ietf.org/html/rfc7231#section-5.5.3) header sent by the
	// client.
	//
	// Type: string
	// Required: No
	// Stability: stable
	// Examples: 'CERN-LineMode/2.15 libwww/2.17b3'
	HTTPUserAgentKey = attribute.Key("http.user_agent")
	// Like `net.peer.ip` but for the host IP. Useful in case of a multi-IP host.
	//
	// Type: string
	// Required: No
	// Stability: stable
	// Examples: '192.168.0.1'
	NetHostIPKey = attribute.Key("net.host.ip")
	// Remote address of the peer (dotted decimal for IPv4 or
	// [RFC5952](https://tools.ietf.org/html/rfc5952) for IPv6)
	//
	// Type: string
	// Required: No
	// Stability: stable
	// Examples: '127.0.0.1'
	NetPeerIPKey = attribute.Key("net.peer.ip")
)

// HTTP scheme attributes.
var (
	HTTPSchemeHTTP  = semconv.HTTPSchemeKey.String("http")
	HTTPSchemeHTTPS = semconv.HTTPSchemeKey.String("https")
	// Another IP-based protocol
	NetTransportIP = semconv.NetTransportKey.String("ip")
	// Unix Domain socket. See below
	NetTransportUnix = semconv.NetTransportKey.String("unix")
)

var sc = &semconvutil.SemanticConventions{
	EnduserIDKey:                semconv.EnduserIDKey,
	HTTPClientIPKey:             HTTPClientIPKey,
	HTTPFlavorKey:               HTTPFlavorKey,
	HTTPHostKey:                 HTTPHostKey,
	HTTPMethodKey:               semconv.HTTPMethodKey,
	HTTPRequestContentLengthKey: semconv.HTTPRequestContentLengthKey,
	HTTPRouteKey:                semconv.HTTPRouteKey,
	HTTPSchemeHTTP:              HTTPSchemeHTTP,
	HTTPSchemeHTTPS:             HTTPSchemeHTTPS,
	HTTPServerNameKey:           HTTPServerNameKey,
	HTTPStatusCodeKey:           semconv.HTTPStatusCodeKey,
	HTTPTargetKey:               semconv.HTTPTargetKey,
	HTTPURLKey:                  semconv.HTTPURLKey,
	HTTPUserAgentKey:            HTTPUserAgentKey,
	NetHostIPKey:                NetHostIPKey,
	NetHostNameKey:              semconv.NetHostNameKey,
	NetHostPortKey:              semconv.NetHostPortKey,
	NetPeerIPKey:                NetPeerIPKey,
	NetPeerNameKey:              semconv.NetPeerNameKey,
	NetPeerPortKey:              semconv.NetPeerPortKey,
	NetTransportIP:              NetTransportIP,
	NetTransportOther:           semconv.NetTransportOther,
	NetTransportTCP:             semconv.NetTransportTCP,
	NetTransportUDP:             semconv.NetTransportUDP,
	NetTransportUnix:            NetTransportUnix,
}

// NetAttributesFromHTTPRequest generates attributes of the net
// namespace as specified by the OpenTelemetry specification for a
// span.  The network parameter is a string that net.Dial function
// from standard library can understand.
func NetAttributesFromHTTPRequest(network string, request *http.Request) []attribute.KeyValue {
	return sc.NetAttributesFromHTTPRequest(network, request)
}

// EndUserAttributesFromHTTPRequest generates attributes of the
// enduser namespace as specified by the OpenTelemetry specification
// for a span.
func EndUserAttributesFromHTTPRequest(request *http.Request) []attribute.KeyValue {
	return sc.EndUserAttributesFromHTTPRequest(request)
}

// HTTPClientAttributesFromHTTPRequest generates attributes of the
// http namespace as specified by the OpenTelemetry specification for
// a span on the client side.
func HTTPClientAttributesFromHTTPRequest(request *http.Request) []attribute.KeyValue {
	return sc.HTTPClientAttributesFromHTTPRequest(request)
}

// HTTPServerMetricAttributesFromHTTPRequest generates low-cardinality attributes
// to be used with server-side HTTP metrics.
func HTTPServerMetricAttributesFromHTTPRequest(serverName string, request *http.Request) []attribute.KeyValue {
	return sc.HTTPServerMetricAttributesFromHTTPRequest(serverName, request)
}

// HTTPServerAttributesFromHTTPRequest generates attributes of the
// http namespace as specified by the OpenTelemetry specification for
// a span on the server side. Currently, only basic authentication is
// supported.
func HTTPServerAttributesFromHTTPRequest(serverName, route string, request *http.Request) []attribute.KeyValue {
	return sc.HTTPServerAttributesFromHTTPRequest(serverName, route, request)
}

// HTTPAttributesFromHTTPStatusCode generates attributes of the http
// namespace as specified by the OpenTelemetry specification for a
// span.
func HTTPAttributesFromHTTPStatusCode(code int) []attribute.KeyValue {
	return sc.HTTPAttributesFromHTTPStatusCode(code)
}

// SpanStatusFromHTTPStatusCode generates a status code and a message
// as specified by the OpenTelemetry specification for a span.
func SpanStatusFromHTTPStatusCode(code int) (codes.Code, string) {
	return semconvutil.SpanStatusFromHTTPStatusCode(code)
}

// SpanStatusFromHTTPStatusCodeAndSpanKind generates a status code and a message
// as specified by the OpenTelemetry specification for a span.
// Exclude 4xx for SERVER to set the appropriate status.
func SpanStatusFromHTTPStatusCodeAndSpanKind(code int, spanKind trace.SpanKind) (codes.Code, string) {
	return semconvutil.SpanStatusFromHTTPStatusCodeAndSpanKind(code, spanKind)
}
