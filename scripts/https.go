package main

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"euphoria/scripts/utils"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Touch(uri, method, body, authorization, contentType, cookie, referer string, headers map[string]string) ([]byte, error) {
	transport := &http.Transport{
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
	}
	req, err := http.NewRequest(method, uri, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	/*
	   accept,
	   acceptCh,
	   acceptChLifetime,
	   acceptCharset,
	   acceptEncoding,
	   acceptLanguage,
	   acceptPatch,
	   acceptRanges,
	   accessControlAllowCredentials,
	   accessControlAllowHeaders,
	   accessControlAllowMethods,
	   accessControlAllowOrigin,
	   accessControlExposeHeaders,
	   accessControlMaxAge,
	   accessControlRequestHeaders,
	   accessControlRequestMethod,
	   age,
	   allow,
	   altSvc,
	   authorization,
	   cacheControl,
	   clearSiteData,
	   connection,
	   contentDisposition,
	   contentEncoding,
	   contentLanguage,
	   contentLength,
	   contentLocation,
	   contentRange,
	   contentSecurityPolicy,
	   contentSecurityPolicyReportOnly,
	   contentType,
	   cookie,
	   cookie2,
	   crossOriginResourcePolicy,
	   dnt,
	   dpr,
	   date,
	   deviceMemory,
	   digest,
	   etag,
	   earlyData,
	   expect,
	   expectCt,
	   expires,
	   featurePolicy,
	   forwarded,
	   from,
	   host,
	   ifMatch,
	   ifModifiedSince,
	   ifNoneMatch,
	   ifRange,
	   ifUnmodifiedSince,
	   index,
	   keepAlive,
	   largeAllocation,
	   lastModified,
	   link,
	   location,
	   origin,
	   pragma,
	   proxyAuthenticate,
	   proxyAuthorization,
	   publicKeyPins,
	   publicKeyPinsReportOnly,
	   range,
	   referer,
	   referrerPolicy,
	   retryAfter,
	   saveData,
	   secWebsocketAccept,
	   server,
	   serverTiming,
	   setCookie,
	   setCookie2,
	   sourcemap,
	   strictTransportSecurity,
	   te,
	   timingAllowOrigin,
	   tk,
	   trailer,
	   transferEncoding,
	   upgradeInsecureRequests,
	   userAgent,
	   vary,
	   via,
	   wwwAuthenticate,
	   wantDigest,
	   warning,
	   xContentTypeOptions,
	   xDnsPrefetchControl,
	   xForwardedFor,
	   xForwardedHost,
	   xForwardedProto,
	   xFrameOptions,
	   xXssProtection,

	*/
	//if accept != "" {req.Header.Set("Accept",accept)}
	//if acceptCh != "" {req.Header.Set("Accept-CH",acceptCh)}
	//if acceptChLifetime != "" {req.Header.Set("Accept-CH-Lifetime",acceptChLifetime)}
	//if acceptCharset != "" {req.Header.Set("Accept-Charset",acceptCharset)}
	//if acceptEncoding != "" {req.Header.Set("Accept-Encoding",acceptEncoding)}
	//if acceptLanguage != "" {req.Header.Set("Accept-Language",acceptLanguage)}
	//if acceptPatch != "" {req.Header.Set("Accept-Patch",acceptPatch)}
	//if acceptRanges != "" {req.Header.Set("Accept-Ranges",acceptRanges)}
	//if accessControlAllowCredentials != "" {req.Header.Set("Access-Control-Allow-Credentials",accessControlAllowCredentials)}
	//if accessControlAllowHeaders != "" {req.Header.Set("Access-Control-Allow-Headers",accessControlAllowHeaders)}
	//if accessControlAllowMethods != "" {req.Header.Set("Access-Control-Allow-Methods",accessControlAllowMethods)}
	//if accessControlAllowOrigin != "" {req.Header.Set("Access-Control-Allow-Origin",accessControlAllowOrigin)}
	//if accessControlExposeHeaders != "" {req.Header.Set("Access-Control-Expose-Headers",accessControlExposeHeaders)}
	//if accessControlMaxAge != "" {req.Header.Set("Access-Control-Max-Age",accessControlMaxAge)}
	//if accessControlRequestHeaders != "" {req.Header.Set("Access-Control-Request-Headers",accessControlRequestHeaders)}
	//if accessControlRequestMethod != "" {req.Header.Set("Access-Control-Request-Method",accessControlRequestMethod)}
	//if age != "" {req.Header.Set("Age",age)}
	//if allow != "" {req.Header.Set("Allow",allow)}
	//if altSvc != "" {req.Header.Set("Alt-Svc",altSvc)}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	//if cacheControl != "" {req.Header.Set("Cache-Control",cacheControl)}
	//if clearSiteData != "" {req.Header.Set("Clear-Site-Data",clearSiteData)}
	//if connection != "" {req.Header.Set("Connection",connection)}
	//if contentDisposition != "" {req.Header.Set("Content-Disposition",contentDisposition)}
	//if contentEncoding != "" {req.Header.Set("Content-Encoding",contentEncoding)}
	//if contentLanguage != "" {req.Header.Set("Content-Language",contentLanguage)}
	//if contentLength != "" {req.Header.Set("Content-Length",contentLength)}
	//if contentLocation != "" {req.Header.Set("Content-Location",contentLocation)}
	//if contentRange != "" {req.Header.Set("Content-Range",contentRange)}
	//if contentSecurityPolicy != "" {req.Header.Set("Content-Security-Policy",contentSecurityPolicy)}
	//if contentSecurityPolicyReportOnly != "" {req.Header.Set("Content-Security-Policy-Report-Only",contentSecurityPolicyReportOnly)}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	//if cookie2 != "" {req.Header.Set("Cookie2",cookie2)}
	//if crossOriginResourcePolicy != "" {req.Header.Set("Cross-Origin-Resource-Policy",crossOriginResourcePolicy)}
	//if dnt != "" {req.Header.Set("DNT",dnt)}
	//if dpr != "" {req.Header.Set("DPR",dpr)}
	//if date != "" {req.Header.Set("Date",date)}
	//if deviceMemory != "" {req.Header.Set("Device-Memory",deviceMemory)}
	//if digest != "" {req.Header.Set("Digest",digest)}
	//if etag != "" {req.Header.Set("ETag",etag)}
	//if earlyData != "" {req.Header.Set("Early-Data",earlyData)}
	//if expect != "" {req.Header.Set("Expect",expect)}
	//if expectCt != "" {req.Header.Set("Expect-CT",expectCt)}
	//if expires != "" {req.Header.Set("Expires",expires)}
	//if featurePolicy != "" {req.Header.Set("Feature-Policy",featurePolicy)}
	//if forwarded != "" {req.Header.Set("Forwarded",forwarded)}
	//if from != "" {req.Header.Set("From",from)}
	//if host != "" {req.Header.Set("Host",host)}
	//if ifMatch != "" {req.Header.Set("If-Match",ifMatch)}
	//if ifModifiedSince != "" {req.Header.Set("If-Modified-Since",ifModifiedSince)}
	//if ifNoneMatch != "" {req.Header.Set("If-None-Match",ifNoneMatch)}
	//if ifRange != "" {req.Header.Set("If-Range",ifRange)}
	//if ifUnmodifiedSince != "" {req.Header.Set("If-Unmodified-Since",ifUnmodifiedSince)}
	//if index != "" {req.Header.Set("Index",index)}
	//if keepAlive != "" {req.Header.Set("Keep-Alive",keepAlive)}
	//if largeAllocation != "" {req.Header.Set("Large-Allocation",largeAllocation)}
	//if lastModified != "" {req.Header.Set("Last-Modified",lastModified)}
	//if link != "" {req.Header.Set("Link",link)}
	//if location != "" {req.Header.Set("Location",location)}
	//if origin != "" {req.Header.Set("Origin",origin)}
	//if pragma != "" {req.Header.Set("Pragma",pragma)}
	//if proxyAuthenticate != "" {req.Header.Set("Proxy-Authenticate",proxyAuthenticate)}
	//if proxyAuthorization != "" {req.Header.Set("Proxy-Authorization",proxyAuthorization)}
	//if publicKeyPins != "" {req.Header.Set("Public-Key-Pins",publicKeyPins)}
	//if publicKeyPinsReportOnly != "" {req.Header.Set("Public-Key-Pins-Report-Only",publicKeyPinsReportOnly)}
	//if range != "" {req.Header.Set("Range",range)}
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	//if referrerPolicy != "" {req.Header.Set("Referrer-Policy",referrerPolicy)}
	//if retryAfter != "" {req.Header.Set("Retry-After",retryAfter)}
	//if saveData != "" {req.Header.Set("Save-Data",saveData)}
	//if secWebsocketAccept != "" {req.Header.Set("Sec-WebSocket-Accept",secWebsocketAccept)}
	//if server != "" {req.Header.Set("Server",server)}
	//if serverTiming != "" {req.Header.Set("Server-Timing",serverTiming)}
	//if setCookie != "" {req.Header.Set("Set-Cookie",setCookie)}
	//if setCookie2 != "" {req.Header.Set("Set-Cookie2",setCookie2)}
	//if sourcemap != "" {req.Header.Set("SourceMap",sourcemap)}
	//if strictTransportSecurity != "" {req.Header.Set("Strict-Transport-Security",strictTransportSecurity)}
	//if te != "" {req.Header.Set("TE",te)}
	//if timingAllowOrigin != "" {req.Header.Set("Timing-Allow-Origin",timingAllowOrigin)}
	//if tk != "" {req.Header.Set("Tk",tk)}
	//if trailer != "" {req.Header.Set("Trailer",trailer)}
	//if transferEncoding != "" {req.Header.Set("Transfer-Encoding",transferEncoding)}
	//if upgradeInsecureRequests != "" {req.Header.Set("Upgrade-Insecure-Requests",upgradeInsecureRequests)}
	//if userAgent != "" {req.Header.Set("User-Agent",userAgent)}
	//if vary != "" {req.Header.Set("Vary",vary)}
	//if via != "" {req.Header.Set("Via",via)}
	//if wwwAuthenticate != "" {req.Header.Set("WWW-Authenticate",wwwAuthenticate)}
	//if wantDigest != "" {req.Header.Set("Want-Digest",wantDigest)}
	//if warning != "" {req.Header.Set("Warning",warning)}
	//if xContentTypeOptions != "" {req.Header.Set("X-Content-Type-Options",xContentTypeOptions)}
	//if xDnsPrefetchControl != "" {req.Header.Set("X-DNS-Prefetch-Control",xDnsPrefetchControl)}
	//if xForwardedFor != "" {req.Header.Set("X-Forwarded-For",xForwardedFor)}
	//if xForwardedHost != "" {req.Header.Set("X-Forwarded-Host",xForwardedHost)}
	//if xForwardedProto != "" {req.Header.Set("X-Forwarded-Proto",xForwardedProto)}
	//if xFrameOptions != "" {req.Header.Set("X-Frame-Options",xFrameOptions)}
	//if xXssProtection != "" {req.Header.Set("X-XSS-Protection",xXssProtection)}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close()

	buf, err := ioutil.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf, nil
}
func main() {
	m := loadingSettings()
	accessToken := m["AccessToken"].(string)

	debug := true
	host := "http://localhost:5050"

	if len(os.Args) > 1 {
		for _, o := range os.Args[1:] {
			if o == "d" {
				debug = true
				break
			}
		}
	}
	if !debug {
		host = "https://halalla.cn"
	}
	if len(os.Args) > 1 {
		for _, o := range os.Args[1:] {
			if o == "s" {
				insertSearch(host, accessToken)
				break
			} else if o == "slide" {
				insertSlide(host, accessToken)
			} else if o == "category" {
				insertCategory(host, accessToken)
			} else if o == "i" {
				insertStore(host, accessToken)
			} else {

			}
		}
		return
	}
	s, err := utils.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(s, "http") {
		fetchStore(s, accessToken)
		return
	}
	f, err := os.Stat(s)
	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasSuffix(f.Name(), ".json") {
		return
	}
	postJson(s, accessToken)
}
func postJson(filename string, accessToken string) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var items map[string]interface{}
	err = json.Unmarshal(buf, &items)
	if err != nil {
		log.Fatal(err)
	}
	buf, err = Touch(items["Url"].(string), "POST", string(buf), "Bearer "+accessToken, "", "", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))

}
func insertSearch(host, accessToken string) {
	buf, err := ioutil.ReadFile("./json/all_search.json")
	if err != nil {
		log.Fatal(err)
	}
	uri := host + "/store/api/search?method=insert"
	buf, err = Touch(uri, "POST", string(buf), "Bearer "+accessToken, "", "", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}
func insertSlide(host, accessToken string) {
	body := "[" + `"TB1kkQMlLb2gK0jSZK9XXaEgFXa-1080-498.jpg","TB1LYdTkuT2gK0jSZFvXXXnFXXa-1080-498.jpg","TB1PRiGm7L0gK0jSZFxXXXWHVXa-1080-498.jpg"` + "]"
	uri := host + "/store/api/slide"
	buf, err := Touch(uri, "POST", body, "Bearer "+accessToken, "", "", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}
func insertStore(host, accessToken string) {
	buf, err := ioutil.ReadFile("./json/all.json")
	if err != nil {
		log.Fatal(err)
	}
	uri := host + "/store/api/store?method=insert"
	buf, err = Touch(uri, "POST", string(buf), "Bearer "+accessToken, "", "", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}

/*
"汽车用品",
"家居百货",
"家用电器",
"女装",
"配饰",
"内衣",
"童装",
"玩具",
"母婴用品",
"包装",
"美妆日化",
"食品生鲜",
"进口",
"家装建材",
"安防",
"五金工具",
"家纺家饰",
"宠物园艺",
"运动户外",
"鞋靴",
"办公文化",
"男装",
"箱包皮具",
"电工电气",
"数码电子",
"机械设备",
"化工冶金",
"纺织皮革",
"照明工业",
"橡塑",
*/
func insertCategory(host, accessToken string) {
	body := "[" + `"汽车用品",
"家居百货",
"家用电器",
"女装",
"配饰",
"内衣",
"童装",
"玩具",
"母婴用品",
"包装",
"美妆日化",
"食品生鲜",
"进口",
"家装建材",
"安防",
"五金工具",
"家纺家饰",
"宠物园艺",
"运动户外",
"鞋靴",
"办公文化",
"男装",
"箱包皮具",
"电工电气",
"数码电子",
"机械设备",
"化工冶金",
"纺织皮革",
"照明工业",
"橡塑"` + "]"
	uri := host + "/store/api/category"
	buf, err := Touch(uri, "POST", body, "Bearer "+accessToken, "", "", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}
func loadingSettings() map[string]interface{} {
	buf, err := ioutil.ReadFile("../settings/settings.json")
	if err != nil {
		log.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
func fetchStore(uri, accessToken string) {
	buf, err := Touch(uri, "GET", "", "Bearer "+accessToken, "", "", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	var m map[string]interface{}

	err = json.Unmarshal(buf, &m)

	uid := m["uid"].(string)

	filename := "./json/update_" + uid + ".json"
	ioutil.WriteFile(filename, buf, 0644)
}

// http://localhost:5050/store/api/store?method=details&uid=lpkxgi
