package main

/* Al useful imports */
import (
    "bufio"
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "strings"
    //"time"
    "os"
    "go-going/gone/messages"
)

func main(){
    /* Parse the provided parameters on command line */
    clusterip := flag.String("clusterip", "127.0.0.1:9001", "ip address of any node to connnect")
    flag.Parse()

    myIp,_ := net.InterfaceAddrs()

    /* Try to connect to the cluster, and send request to cluster if able to connect */
    fmt.Println("Initiating client. Connecting to cluster.")
    connectToCluster(myIp[0].String(), *clusterip)
}


func connectToCluster(myIp string, clusterip string) {
    url := fmt.Sprintf("http://%s/query", clusterip)
    fmt.Println("URL: ", url)

    for {
        query, _ := bufio.NewReader(os.Stdin).ReadString('\n')
        //commandJson := fmt.Sprintf(`{"Query" : "%s"}`, strings.TrimSpace(query))

        commandJson := messages.JsonRequest{
            JsonRequestString : strings.TrimSpace(query),
        }
        fmt.Println("Json req:", commandJson)
        var buf []byte
        buf, _ = json.Marshal(commandJson)

        req, _ := http.NewRequest("POST", url, bytes.NewBuffer(buf))
        //req.Header.Set("X-Custom-Header", "myvalue")
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        resp, _ := client.Do(req)
        //defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Println(body)
    }
}

