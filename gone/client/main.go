package main

/* Al useful imports */
import (
    "flag"
    "fmt"
    "net"
    "strings"
    "time"
    "encoding/json"

    "go-going/gone/messages"
)


func main(){
    /* Parse the provided parameters on command line */
    clusterip := flag.String("clusterip", "127.0.0.1:8001", "ip address of any node to connnect")
    myport := flag.String("myport", "9001", "ip address to run client on. default is 9001.")
    flag.Parse()

    myIp,_ := net.InterfaceAddrs()
    me := messages.NodeInfo{ NodeId: -1, NodeIpAddr: myIp[0].String(), Port: *myport}
    remote := messages.NodeInfo{ NodeId: -1, NodeIpAddr: strings.Split(*clusterip, ":")[0], Port: strings.Split(*clusterip, ":")[1]}

    /* Try to connect to the cluster, and send request to cluster if able to connect */
    fmt.Println("Initiating client. Connecting to cluster.", remote)
    connectToCluster(me, remote)
}


func connectToCluster(me messages.NodeInfo, remote messages.NodeInfo) (bool){
    /* connect to this socket details provided */
    connOut, err := net.DialTimeout("tcp", remote.NodeIpAddr + ":" + remote.Port, time.Duration(10) * time.Second)
    if err != nil {
        if _, ok := err.(net.Error); ok {
            fmt.Println("Couldn't connect to cluster.", me.NodeId)
            return false
        }
    } else {
        fmt.Println("Connected to cluster. Sending message to node.")
        text := "Hi Cluster.. I am a client.."
        requestMessage := messages.GetAddToClusterMessage(me, remote, text)
        json.NewEncoder(connOut).Encode(&requestMessage)

        decoder := json.NewDecoder(connOut)
        var responseMessage messages.AddToClusterMessage
        decoder.Decode(&responseMessage)
        fmt.Println("Got response:\n" + responseMessage.String())
        
        return true
    }
    return false
}

