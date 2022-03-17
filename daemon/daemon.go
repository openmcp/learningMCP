package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/juju/fslock"
)

const INDEX_HTML = `
<!DOCTYPE html>
<html>
<head>
    <title>shell2http</title>
</head>
<body>
    <h1>shell2http</h1>
    <ul>
        %s
        <li><a href="/exit">/exit</a></li>
    </ul>
    Get from: <a href="https://github.com/msoap/shell2http">github.com/msoap/shell2http</a>
</body>
</html>
`

var HOME string
var STATUS_HOME string
var MANAGER_HOME string
var clusterNameMap map[string][]string
var openmcpCnt int
var waitTime time.Duration
var schedulerIP string
var schedulerPort string

func setConfig() {
	schedulerIP = "10.0.3.60"
	schedulerPort = "3124"

	HOME = "/home"
	STATUS_HOME = "/home/keti/learningMCP/status"
	MANAGER_HOME = "/home/keti/learningMCP/manager"

	waitTime = 3 * time.Second

	clusterNameMap = map[string][]string{
		"openmcp1": []string{"cluster01", "cluster02", "cluster03"},
		// "openmcp2": []string{"cluster04", "cluster05", "cluster06"},
		//"openmcp3": []string{"cluster07", "cluster08", "cluster09"},
		//"openmcp4": []string{"cluster10", "cluster11", "cluster12"},
		//"openmcp5": []string{"cluster13", "cluster14", "cluster15"},
		//"openmcp6": []string{"cluster16", "cluster17", "cluster18"},
	}

	openmcpCnt = len(clusterNameMap)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// stat = [created, creating, connected, terminated, needcreate, deleting, unknown]
// func resultStatus(status []string) string {
// 	temp := status[0]
// 	flag := true
// 	for i := 1; i < len(status); i++ {
// 		if temp != status[i] {
// 			flag = false
// 			break
// 		}
// 		if status[i] == "creating" {
// 			return "creating"
// 		}
// 	}
// 	if flag { // 모든 클러스터의 상태 일치
// 		if temp == "" {
// 			return "needcreate"
// 		}
// 		return temp
// 	} else {
// 		// 상태가 다름
// 		return "unknown"
// 	}

// }

type clusterStatus struct {
	OpenmcpName string
	ClusterName string
	Status      string
	ServerIP    string
}

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// os.Stdout.WriteString(ipnet.IP.String() + "\n")
				return ipnet.IP.String()
			}
		}
	}
	return ""

}
func doPostStatus(openmcpName, clusterName, sts string) {
	s := clusterStatus{
		OpenmcpName: openmcpName,
		ClusterName: clusterName,
		Status:      sts,
		ServerIP:    getMyIP(),
	}

	pbytes, _ := json.Marshal(s)
	buff := bytes.NewBuffer(pbytes)

	resp, err := http.Post("http://"+schedulerIP+":"+schedulerPort+"/status", "application/json", buff)
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			_ = data
			// fmt.Printf("%s\n", string(data))
		}
	}

}
func statusGoRoutine(wg *sync.WaitGroup) {
	defer wg.Done()
	for {

		for openmcpName, clusterNames := range clusterNameMap {
			// fmt.Println("--------Status Check: " + openmcpName + " ---------")

			status := []string{}

			for i, clusterName := range clusterNames {

				sts := "needcreate"
				if fileExists(path.Join(STATUS_HOME, clusterName)) {
					lock := fslock.New("/tmp/" + clusterName + ".lock")
					lockErr := lock.Lock()
					if lockErr != nil {
						fmt.Println("falied to acquire lock > " + lockErr.Error())
						break
					}
					bytes, err := ioutil.ReadFile(path.Join(STATUS_HOME, clusterName))

					lock.Unlock()

					if err != nil {
						log.Fatal(err)
					}

					if string(bytes[:len(bytes)-1]) != "" {
						sts = string(bytes[:len(bytes)-1])
					}
				}
				status = append(status, sts)

				// fmt.Println(clusterName + ": " + status[i])
				doPostStatus(openmcpName, clusterName, status[i])
				if status[i] == "terminated" || status[i] == "needcreate" || status[i] == "envCreated" || status[i] == "deleted" {
					go managerGoRoutine(status[i], clusterName, i)
				}
			}
			// result := resultStatus(status)
			// fmt.Println("=> Result: ", result)

			// if result == "terminated" || result == "needcreate" {
			// 	for i, clusterName := range clusterNames {
			// 		go managerGoRoutine(result, clusterName, i)
			// 	}

			// }

		}

		time.Sleep(waitTime)
	}
}

func managerGoRoutine(status, clusterName string, index int) {
	if status == "needcreate" {
		if index == 0 { // Master
			cmdExec2(path.Join(MANAGER_HOME, "1.create_master_env.sh")+" "+clusterName, clusterName)
		} else {
			cmdExec2(path.Join(MANAGER_HOME, "2.create_member_env.sh")+" "+clusterName, clusterName)
		}

	} else if status == "terminated" {
		cmdExec2(path.Join(MANAGER_HOME, "5.delete.sh")+" "+clusterName, clusterName)
	} else if status == "envCreated" || status == "deleted" {
		if index == 0 { // Master
			cmdExec2(path.Join(MANAGER_HOME, "3.create_openmcp_master_cluster.sh")+" "+clusterName, clusterName)
		} else {
			cmdExec2(path.Join(MANAGER_HOME, "4.create_openmcp_member_cluster.sh")+" "+clusterName, clusterName)
		}
	}

}
func cmdExec2(cmdStr, clusterName string) error {
	cmd := exec.Command("bash", "-c", cmdStr)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println("[" + clusterName + "] " + stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Println("[" + clusterName + "] " + stderrScanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}

	return nil
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	localPath := "/home/keti/learningMCP/web" + req.URL.Path
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	contentType := getContentType(localPath)
	w.Header().Add("Content-Type", contentType)
	w.Write(content)
}
func getContentType(localPath string) string {
	var contentType string
	ext := filepath.Ext(localPath)

	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	default:
		contentType = "text/plain"
	}

	return contentType
}

func webServerGoRoutine(wg *sync.WaitGroup) {
	defer wg.Done()

	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":5000", nil)
}
func SplitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}
func removeDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	fmt.Println("LearningMCP Daemon Start")
	setConfig()

	var wg sync.WaitGroup
	wg.Add(3)

	go statusGoRoutine(&wg)
	//go webServerGoRoutine(&wg)

	wg.Wait()

}
