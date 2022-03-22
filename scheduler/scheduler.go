package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type HttpManager struct {
}
type ClusterStatus struct {
	OpenmcpName string
	ClusterName string //`json:"clusterName"`
	Status      string //`json:"status"`
	ServerIP    string
}
type sshInfo struct {
	OpenmcpName string
	IP          string
	Port        string
	UserName    string
	PW          string
}

var clusterStatusMap map[string]map[string]string
var openmcpIPMap map[string]string
var mu sync.Mutex

func (h *HttpManager) redirectMasterHandler(w http.ResponseWriter, r *http.Request) {
	openmcpName := getSchedulingReadyOpenMCP()
	IP, PORT := getSchedulingMasterIP(openmcpName)
	fmt.Println(openmcpName, "master", IP, PORT)
	http.Redirect(w, r, "http://"+IP+":"+PORT, 307)

}
func (h *HttpManager) redirectCluster1Handler(w http.ResponseWriter, r *http.Request) {
	openmcpName := getSchedulingReadyOpenMCP()

	IP, PORT := getSchedulingCluster1IP(openmcpName)
	fmt.Println(openmcpName, "cluster1", IP, PORT)
	http.Redirect(w, r, "http://"+IP+":"+PORT, 307)

}
func (h *HttpManager) redirectCluster2Handler(w http.ResponseWriter, r *http.Request) {
	openmcpName := getSchedulingReadyOpenMCP()
	IP, PORT := getSchedulingCluster2IP(openmcpName)
	fmt.Println(openmcpName, "cluster2", IP, PORT)
	http.Redirect(w, r, "http://"+IP+":"+PORT, 307)

}
func (h *HttpManager) statusHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.RemoteAddr)
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	if r.Body == nil {
		//http.Error(w, "Please send a request body", 400)
		errorResponse(w, "Please send a request body", 400)
		return
	}
	var s ClusterStatus
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&s)

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	//fmt.Println(s.OpenmcpName, s.ClusterName, s.Status)
	mu.Lock()
	if _, ok := clusterStatusMap[s.OpenmcpName]; ok {
		clusterStatusMap[s.OpenmcpName][s.ClusterName] = s.Status

	} else {
		tempMap := make(map[string]string)
		clusterStatusMap[s.OpenmcpName] = tempMap
		clusterStatusMap[s.OpenmcpName][s.ClusterName] = s.Status
	}
	if _, ok := openmcpIPMap[s.OpenmcpName]; ok {
		openmcpIPMap[s.OpenmcpName] = s.ServerIP
	} else {
		openmcpIPMap = make(map[string]string)
		openmcpIPMap[s.OpenmcpName] = s.ServerIP
	}

	for k, v := range clusterStatusMap {
		for k2, v2 := range v {
			fmt.Println("openmcp : ", k, ", cluster: ", k2, ", status: ", v2, " IP: ", s.ServerIP)
		}
	}
	mu.Unlock()

	errorResponse(w, "Success", http.StatusOK)
	return
}
func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// func (h *HttpManager) schedulingHandler(w http.ResponseWriter, r *http.Request) {
// 	enableCors(&w)

// 	s := sshInfo{
// 		OpenmcpName: "openmcp2",
// 		IP:          "10.0.3.60",
// 		Port:        "57575",
// 		UserName:    "cluster04",
// 		PW:          "1234",
// 	}

// 	bytesJson, _ := json.Marshal(s)
// 	var prettyJSON bytes.Buffer
// 	err := json.Indent(&prettyJSON, bytesJson, "", "\t")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	w.Write(prettyJSON.Bytes())
// }
func getSchedulingReadyOpenMCP() string {

	sortKeys := make([]string, 0, len(clusterStatusMap))

	for k := range clusterStatusMap {
		fmt.Println("k", k)
		sortKeys = append(sortKeys, k)
	}
	sort.Strings(sortKeys) //keys로 정렬을 함

	//정렬한 keys 값으로 데이터를 출력함
	for _, k := range sortKeys {
		fmt.Println(k, clusterStatusMap[k])
	}

	for _, key := range sortKeys {
		k := key
		v := clusterStatusMap[key]

		findReady := true

		for k2, v2 := range v {
			fmt.Println("openmcp : ", k, ", cluster: ", k2, ", status: ", v2)

			if v2 != "clusterCreated" {
				findReady = false
				break
			}

		}
		if findReady {
			return k
		}

	}

	return ""
}

func getSchedulingMasterIP(openmcpName string) (string, string) {

	min_num := -1

	for k, _ := range clusterStatusMap[openmcpName] {
		n := k[len(k)-2:]
		n_int, _ := strconv.Atoi(n)

		if min_num == -1 || n_int < min_num {
			min_num = n_int

		}
	}
	IP := openmcpIPMap[openmcpName]
	return IP, "575" + fmt.Sprintf("%02d", min_num)
}
func getSchedulingCluster1IP(openmcpName string) (string, string) {

	min_num := -1

	for k, _ := range clusterStatusMap[openmcpName] {
		n := k[len(k)-2:]
		n_int, _ := strconv.Atoi(n)

		if min_num == -1 || n_int < min_num {
			min_num = n_int

		}
	}
	IP := openmcpIPMap[openmcpName]
	return IP, "575" + fmt.Sprintf("%02d", min_num+1)
}
func getSchedulingCluster2IP(openmcpName string) (string, string) {

	min_num := -1

	for k, _ := range clusterStatusMap[openmcpName] {
		n := k[len(k)-2:]
		n_int, _ := strconv.Atoi(n)

		if min_num == -1 || n_int < min_num {
			min_num = n_int

		}
	}
	IP := openmcpIPMap[openmcpName]
	return IP, "575" + fmt.Sprintf("%02d", min_num+2)
}
func getSchedulingIP() (string, string) {

	for k, v := range clusterStatusMap {

		findReady := true

		min_num := -1

		for k2, v2 := range v {
			fmt.Println("openmcp : ", k, ", cluster: ", k2, ", status: ", v2)

			if v2 != "clusterCreated" {
				findReady = false
				break
			}
			n := k2[len(k2)-2:]
			n_int, _ := strconv.Atoi(n)

			if min_num == -1 || n_int < min_num {
				min_num = n_int

			}

		}
		IP := openmcpIPMap[k]

		if findReady {
			return IP, "575" + fmt.Sprintf("%02d", min_num)
		}

	}

	return "", ""
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}

func main() {
	httpManager := &HttpManager{}
	// HTTPServer_PORT := "3124"
	// handler := http.NewServeMux()

	clusterStatusMap = make(map[string]map[string]string)

	// handler.HandleFunc("/status", httpManager.statusHandler)
	// // handler.HandleFunc("/master", httpManager.redirectMasterHandler)
	// // handler.HandleFunc("/cluster1", httpManager.redirectCluster1Handler)
	// // handler.HandleFunc("/cluster2", httpManager.redirectCluster2Handler)

	// masterDirector := func(req *http.Request) {

	// 	openmcpName := getSchedulingReadyOpenMCP()
	// 	IP, PORT := getSchedulingMasterIP(openmcpName)
	// 	fmt.Println(openmcpName, "master", IP, PORT)

	// 	origin, _ := url.Parse("http://" + IP + ":" + PORT)

	// 	req.Header.Add("X-Forwarded-Host", req.Host)
	// 	req.Header.Add("X-Origin-Host", origin.Host)
	// 	req.URL.Scheme = "http"
	// 	req.URL.Host = origin.Host
	// 	req.Host = origin.Host
	// 	if req.URL.Path == "/master" {
	// 		req.URL.Path = ""
	// 	}
	// 	fmt.Println(req.RequestURI)

	// }

	// cluster1Director := func(req *http.Request) {

	// 	openmcpName := getSchedulingReadyOpenMCP()
	// 	IP, PORT := getSchedulingCluster1IP(openmcpName)
	// 	fmt.Println(openmcpName, "cluster1", IP, PORT)

	// 	origin, _ := url.Parse("http://" + IP + ":" + PORT)

	// 	req.Header.Add("X-Forwarded-Host", req.Host)
	// 	req.Header.Add("X-Origin-Host", origin.Host)
	// 	req.URL.Scheme = "http"
	// 	req.URL.Host = origin.Host
	// }

	// cluster2Director := func(req *http.Request) {

	// 	openmcpName := getSchedulingReadyOpenMCP()
	// 	IP, PORT := getSchedulingCluster2IP(openmcpName)
	// 	fmt.Println(openmcpName, "cluster2", IP, PORT)

	// 	origin, _ := url.Parse("http://" + IP + ":" + PORT)

	// 	req.Header.Add("X-Forwarded-Host", req.Host)
	// 	req.Header.Add("X-Origin-Host", origin.Host)
	// 	req.URL.Scheme = "http"
	// 	req.URL.Host = origin.Host
	// }

	// proxyMaster := &httputil.ReverseProxy{Director: masterDirector}
	// proxyCluster1 := &httputil.ReverseProxy{Director: cluster1Director}
	// proxyClsuter2 := &httputil.ReverseProxy{Director: cluster2Director}

	// handler.HandleFunc("/master", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("sctest")
	// 	proxyMaster.ServeHTTP(w, r)

	// })

	// handler.HandleFunc("/cluster1", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("sctest")
	// 	proxyCluster1.ServeHTTP(w, r)

	// })
	// handler.HandleFunc("/cluster2", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("sctest")
	// 	proxyClsuter2.ServeHTTP(w, r)

	// })

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env File")
	}

	// log.Printf("Forwarding Target : %s://%s", remote.Scheme, remote.Host)

	director := func(req *http.Request) {
		openmcpName := getSchedulingReadyOpenMCP()
		IP, PORT := getSchedulingMasterIP(openmcpName)
		fmt.Println(IP, PORT)

		if req.URL.Path == "/master" {
			IP, PORT = getSchedulingMasterIP(openmcpName)
		} else if req.URL.Path == "/cluster1" {
			IP, PORT = getSchedulingCluster1IP(openmcpName)
		} else if req.URL.Path == "/cluster2" {
			IP, PORT = getSchedulingCluster2IP(openmcpName)
		}
		fmt.Println(IP, PORT)

		//PORT = "57501"

		remote, err := url.Parse("http://" + IP + ":" + PORT)
		if err != nil {
			panic(err)
		}

		targetQuery := remote.RawQuery

		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = singleJoiningSlash(remote.Path, req.URL.Path)

		req.Header.Set("X-Forwarded-For-Host", req.URL.Host)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ModifyResponse = corsHeaderModify

	http.HandleFunc("/", handler(proxy))
	http.HandleFunc("/status", httpManager.statusHandler)

	err = http.ListenAndServe(":"+os.Getenv("LOCAL_PORT"), nil)
	if err != nil {
		panic(err)
	}

}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, r *http.Request) {
		//  If, current request is pre-flight
		if r.Method == "OPTIONS" {
			resp.Header().Set("Access-Control-Allow-Origin", os.Getenv("ACCESS_CONTROL_ALLOWS_ORIGIN"))
			resp.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, content-type")
			resp.Header().Set("Access-Control-Allow-Methods", "*")
			resp.Header().Set("Access-Control-Expose-Headers", "Set-Cookie, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Credential, Authorization")
			resp.Header().Set("Vary", "Origin")
			resp.Header().Set("Vary", "Access-Control-Request-Method")
			resp.Header().Set("Vary", "Access-Control-Request-Headers")
			resp.Header().Set("Access-Control-Allow-Credentials", "true")
			return
		} else {
			if r.URL.Path == "/master" || r.URL.Path == "/cluster1" || r.URL.Path == "/cluster2" {
				r.URL.Path = ""
			}
			log.Printf("%s %s -> Cros_Proxy -> %s", r.Method, r.Host+r.URL.RequestURI(), os.Getenv("TARGET_HOST")+r.URL.RequestURI())
			p.ServeHTTP(resp, r)
		}
	}
}

func corsHeaderModify(resp *http.Response) error {
	// Set Basic Cors related header
	resp.Header.Set("Access-Control-Allow-Origin", os.Getenv("ACCESS_CONTROL_ALLOWS_ORIGIN"))
	resp.Header.Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, content-type")
	resp.Header.Set("Access-Control-Allow-Methods", "*")
	resp.Header.Set("Access-Control-Expose-Headers", "Set-Cookie, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Credential, Authorization")
	resp.Header.Set("Vary", "Origin")
	resp.Header.Set("Vary", "Access-Control-Request-Method")
	resp.Header.Set("Vary", "Access-Control-Request-Headers")
	resp.Header.Set("Access-Control-Allow-Credentials", "true")

	// Parsing cookie in header
	for _, value := range strings.Split(resp.Header.Get("Set-Cookie"), ";") {
		// If remove the domain value, the client host information is automatically set to the domain value by the browser.
		if strings.Contains(value, "Domain=") {
			var newCookie = strings.Replace(resp.Header.Get("Set-Cookie"), value, "", 1)
			resp.Header.Set("Set-Cookie", newCookie)
		}
	}
	return nil
}
