package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
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

func main() {
	httpManager := &HttpManager{}
	HTTPServer_PORT := "3124"
	handler := http.NewServeMux()

	clusterStatusMap = make(map[string]map[string]string)

	handler.HandleFunc("/status", httpManager.statusHandler)
	handler.HandleFunc("/master", httpManager.redirectMasterHandler)
	handler.HandleFunc("/cluster1", httpManager.redirectCluster1Handler)
	handler.HandleFunc("/cluster2", httpManager.redirectCluster2Handler)
	// handler.HandleFunc("/scheduling", httpManager.schedulingHandler)

	// director := func(req *http.Request) {

	// 	IP, PORT := getSchedulingIP()

	// 	origin, _ := url.Parse("http://" + IP + ":" + PORT)
	// 	fmt.Println(origin.Host)
	// 	req.Header.Add("X-Forwarded-Host", req.Host)
	// 	req.Header.Add("X-Origin-Host", origin.Host)
	// 	req.URL.Scheme = "http"
	// 	req.URL.Host = origin.Host
	// 	// req.URL.Path = "/"

	// }

	// proxy := &httputil.ReverseProxy{Director: director}

	// handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("sctest")
	// 	proxy.ServeHTTP(w, r)

	// })

	server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}
	log.Fatal(server.ListenAndServe())

}
