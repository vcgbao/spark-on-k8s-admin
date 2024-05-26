/*
Copyright https://github.com/datapunchorg

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiConfig struct {
	SparkApplicationNamespace string
	SparkUIServiceUrl         string
	ModifyRedirectUrl         bool
}

var sparkUIAppNameURLRegex = regexp.MustCompile("{{\\s*[$]appName\\s*}}")
var sparkUIAppNamespaceURLRegex = regexp.MustCompile("{{\\s*[$]appNamespace\\s*}}")

func getSparkUIServiceUrl(sparkUIServiceUrlFormat string, appName string, appNamespace string) string {
	return sparkUIAppNamespaceURLRegex.ReplaceAllString(sparkUIAppNameURLRegex.ReplaceAllString(sparkUIServiceUrlFormat, appName), appNamespace)
}

func ServeSparkUI(c *gin.Context, config *ApiConfig, uiRootPath string) {
	path := c.Param("path")
	// remove / prefix if there is any
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	// get application name
	appName := ""
	index := strings.Index(path, "/")
	if index <= 0 {
		appName = path
		path = ""
	} else {
		appName = path[0:index]
		path = path[index+1:]
	}
	// get url for the underlying Spark UI Kubernetes service, which is created by spark-on-k8s-operator
	sparkUIServiceUrl := getSparkUIServiceUrl(config.SparkUIServiceUrl, appName, config.SparkApplicationNamespace)
	proxyBasePath := ""
	if config.ModifyRedirectUrl {
		proxyBasePath = fmt.Sprintf("%s/%s", uiRootPath, appName)
	}
	proxy, err := newReverseProxy(sparkUIServiceUrl, path, proxyBasePath)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create reverse proxy for application %s: %s", appName, err.Error()))
		return
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func newReverseProxy(sparkUIServiceUrl string, targetPath string, proxyBasePath string) (*httputil.ReverseProxy, error) {
	log.Printf("Creating revers proxy for Spark UI service url %s", sparkUIServiceUrl)
	targetUrl := sparkUIServiceUrl
	if targetPath != "" {
		if !strings.HasPrefix(targetPath, "/") {
			targetPath = "/" + targetPath
		}
		targetUrl += targetPath
	}
	url, err := url.Parse(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target Spark UI url %s: %s", targetUrl, err.Error())
	}
	director := func(req *http.Request) {
		url.RawQuery = req.URL.RawQuery
		url.RawFragment = req.URL.RawFragment
		log.Printf("Reverse proxy: serving backend url %s for originally requested url %s", url, req.URL)
		req.URL = url
	}
	modifyResponse := func(resp *http.Response) error {
		if proxyBasePath != "" && resp.StatusCode == http.StatusFound {
			// Append the proxy base path before the redirect path.
			// Also modify redirect url to only contain path and not contain host name,
			// so redirect will retain the original requested host name.
			headerName := "Location"
			locationHeaderValues := resp.Header[headerName]
			if len(locationHeaderValues) > 0 {
				newValues := make([]string, 0, len(locationHeaderValues))
				for _, oldHeaderValue := range locationHeaderValues {
					parsedUrl, err := url.Parse(oldHeaderValue)
					if err != nil {
						log.Printf("Reverse proxy: invalid response header value %s: %s (backend url %s): %s", headerName, oldHeaderValue, url, err.Error())
						newValues = append(newValues, oldHeaderValue)
					} else {
						parsedUrl.Scheme = ""
						parsedUrl.Host = ""
						newPath := parsedUrl.Path
						if !strings.HasPrefix(newPath, "/") {
							newPath = "/" + newPath
						}
						parsedUrl.Path = proxyBasePath + newPath
						newHeaderValue := parsedUrl.String()
						log.Printf("Reverse proxy: modifying response header %s from %s to %s (backend url %s)", headerName, oldHeaderValue, newHeaderValue, url)
						newValues = append(newValues, newHeaderValue)
					}
				}
				resp.Header[headerName] = newValues
			}
		}
		return nil
	}
	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}, nil
}
