package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"fmt"
	"strconv"
	"runtime"
	"github.com/op/go-logging"
	"github.com/gin-gonic/gin"
)

var log = logging.MustGetLogger("main")

type HttpResp struct {
	StatusCode int
	Context    string
}

func main() {
	fmt.Println("start api")
	router := gin.Default()

	// curl  http://localhost:80/
	router.GET("/", test)
	router.GET("/getJson", getContainersHandler)
	router.Run(":80")
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello",})
}

var caFile = flag.String("CA", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt", "A PEM eoncoded CA's certificate file.")

func getContainersHandler(c *gin.Context) {

	log.Info("===================")
	log.Info("GET All")
	
	PrintMemUsage()

	// Set CA Cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Errorf("caCert error", err)
		return
	}

	// Get Pod info
	log.Info(" Get Pod")
	urlStrPod := "/api/v1/namespaces/dnn/pods"
	
	PrintMemUsage()

	log.Info(urlStrPod)
	httpRespPod := kubeApiGet(urlStrPod, caCert, c)
	log.Info(strconv.Itoa(httpRespPod.StatusCode))

	// Get Service info
	log.Infof(" Get Service")
	urlStrService := "/api/v1/namespaces/dnn/services"
	
	PrintMemUsage()
	
	log.Info(urlStrService)
	httpRespSvc := kubeApiGet(urlStrService, caCert, c)
	log.Info(strconv.Itoa(httpRespSvc.StatusCode))

	if httpRespPod.StatusCode == 200 && httpRespSvc.StatusCode == 200 {
		respStr := "[" + httpRespPod.Context + "," + httpRespSvc.Context + "]"
		c.String(http.StatusOK, respStr)
		return
	} else if httpRespPod.StatusCode != 200 && httpRespSvc.StatusCode != 200 {
		log.Errorf("Can't find Pod & Service ")
		return
	} else if httpRespPod.StatusCode != 200 {
		c.String(http.StatusOK, httpRespSvc.Context)
		return
	} else if httpRespSvc.StatusCode != 200 {
		c.String(http.StatusOK, httpRespPod.Context)
		return
	}
}

func kubeApiGet(apiUrl string, caCert []byte, c *gin.Context) (httpResp HttpResp) {
	log.Info("start kubeApiGet")
	log.Info("apiUrl:" + apiUrl)
	runtime.GC()

	PrintMemUsage()
	
	// Create cert pool
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Get Token
	log.Info("Token")
	PrintMemUsage()

	buf, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Info("Token error")
	}
	tokenStr := string(buf)

	// Setup HTTPS client
	log.Info(" Setup HTTPS client")
	PrintMemUsage()

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	//transport := &http.Transport{TLSClientConfig: tlsConfig}
	//client := &http.Client{Transport: transport}

	keepAliveTimeout:= 600 * time.Second
	timeout:= 2 * time.Second
	defaultTransport := &http.Transport{
		Dial: (&net.Dialer{KeepAlive: keepAliveTimeout,}).Dial,
		MaxIdleConns: 100,
		MaxIdleConnsPerHost: 100,
		TLSClientConfig: tlsConfig,
	}
	client:= &http.Client{
		Transport: defaultTransport,
		Timeout:   timeout,
	}

	PrintMemUsage()

	req, err := http.NewRequest("GET", "https://kubernetes.default"+apiUrl, nil)
	if err != nil {
		log.Info("http.NewRequest error", err)
		return HttpResp{http.StatusUnauthorized, "call k8s api fail"}
	}
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	PrintMemUsage()
	// Send Reqest
	resp, err := client.Do(req)
	if err != nil {
		log.Info("client.Do error", err)
		return HttpResp{http.StatusUnauthorized, "call k8s api fail"}
	}
	defer resp.Body.Close()

	PrintMemUsage()

	context := convertBody2Str(resp)
	
	return HttpResp{resp.StatusCode, context}
}

func convertBody2Str(resp *http.Response) (context string) {
	log.Info("start convertBody2Str")
	PrintMemUsage()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info("Error ioutil.ReadAll")
		log.Info(string(data))
		return
	}
	PrintMemUsage()
	return string(data)
}


func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Info("Alloc = ", bToMb(m.Alloc))
	log.Info("TotalAlloc = ", bToMb(m.TotalAlloc))
	log.Info("NumGC = ", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b 
}