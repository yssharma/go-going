package main

/* All useful imports */
import (
    "flag"
    "fmt"
    "net/http"
    "net"
    "os"
    "strings"
    "time"
    "math/rand"
    "encoding/json"
    "go-going/gone/messages"
)


var me messages.NodeInfo


/* The entry point for our System */
func main(){
    /* Parse the provided parameters on command line */
    makeMasterOnError := flag.Bool("makeMasterOnError", false, "make this node master if unable to connect to the cluster ip provided.")
    clusterip := flag.String("clusterip", "127.0.0.1:8001", "ip address of any node to connnect")
    myport := flag.String("myport", "8001", "ip address to run this node on. default is 8001.")
    flag.Parse()

    /* Generate id for myself */
    rand.Seed(time.Now().UTC().UnixNano())
    myid := rand.Intn(99999999)

    myIp,_ := net.InterfaceAddrs()
    me = messages.NodeInfo{ NodeId: myid, NodeIpAddr: myIp[0].String(),  Port: *myport}
    dest := messages.NodeInfo{ NodeId: -1, NodeIpAddr: strings.Split(*clusterip, ":")[0], Port: strings.Split(*clusterip, ":")[1]}
    fmt.Println("My details:", me.String())

    /* Try to connect to the cluster, and send request to cluster if able to connect */
    ableToConnect := connectToCluster(me, dest)

    /* 
     * Listen for other incoming requests form other nodes to join cluster
     * Note: We are not doing anything fancy right now to make this node as master. Not yet!
     */
    //isMaster := false

     // Start HTTP server on every node //
    // Todo : Use Nginx to load balance the http servers //
    // Also handle dns rebalancing with nginx when few nodes are down //
    startHttpServer()
    
    if !ableToConnect && *makeMasterOnError {
        fmt.Println("Will start this node as master.")
        //isMaster = true
        listenOnPort(me)
    } else if ableToConnect {
        fmt.Println("Will start this node as slave.")
        listenOnPort(me)
    } else {
        fmt.Println("Quitting system. Set makeMasterOnError flag to make the node master.", myid)
        os.Exit(-1)
    }
}


func connectToCluster(me messages.NodeInfo, dest messages.NodeInfo) (bool){
    /* connect to this socket details provided */
    connOut, err := net.DialTimeout("tcp", dest.NodeIpAddr + ":" + dest.Port, time.Duration(10) * time.Second)
    if err != nil {
        if _, ok := err.(net.Error); ok {
            fmt.Println("Couldn't connect to cluster.", me.NodeId)
            return false
        }
    } else {
        fmt.Println("Connected to cluster. Sending message to node.")
        text := "Hi nody.. please add me to the cluster.."
        requestMessage := messages.GetAddToClusterMessage(me, dest, text)
        json.NewEncoder(connOut).Encode(&requestMessage)

        decoder := json.NewDecoder(connOut)
        var responseMessage messages.AddToClusterMessage
        decoder.Decode(&responseMessage)
        fmt.Println("Got response:\n" + responseMessage.String())
        
        return true
    }
    return false
}


func listenOnPort(me messages.NodeInfo){
    /* Listen for incoming messages */
    ln, _ := net.Listen("tcp", fmt.Sprint(":" + me.Port))
    /* accept connection on port */
    /* not sure if looping infinetely on ln.Accept() is good idea */
    for{
        connIn, err := ln.Accept()
        if err != nil {
            if _, ok := err.(net.Error); ok {
                fmt.Println("Error received while listening.", me.NodeId)
            }
        } else {
            var requestMessage messages.AddToClusterMessage
            json.NewDecoder(connIn).Decode(&requestMessage)
            fmt.Println("Got request:\n" + requestMessage.String())

            text := "Sure buddy.. too easy.."
            responseMessage := messages.GetAddToClusterMessage(me, requestMessage.Source, text)
            json.NewEncoder(connIn).Encode(&responseMessage)
            connIn.Close()
        }
    }
}


func startHttpServer() {
    fmt.Println("Starting http server.")
    http.HandleFunc("/query", queryHandler)
    http.HandleFunc("/", homeHandler)
    http.ListenAndServe(":9001", nil)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var t messages.JsonRequest
    decoder.Decode(&t)
    fmt.Println("Got request string : ", t.JsonRequestString)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Me : ", me.String())
}
