package messages


import (
    "strconv"
)

/* Information/Metadata about node */
type NodeInfo struct {
    NodeId  int  `json:"nodeId"`
    NodeIpAddr  string  `json:"nodeIpAddr"`
    Port  string  `json:"port"`
}

/* A standard format for a Request/Response for adding node to cluster */
type AddToClusterMessage struct {
    Source NodeInfo  `json:"source"`
    Dest NodeInfo  `json:"dest"`
    Message string  `json:"message"`
}

/* Just for pretty printing the node info */
func (node NodeInfo) String() string {
    return "NodeInfo:{ nodeId:" + strconv.Itoa(node.NodeId) + ", nodeIpAddr:" + node.NodeIpAddr + ", port:" + node.Port + " }"
}

/* Just for pretty printing Request/Response info */
func (req AddToClusterMessage) String() string {
    return "AddToClusterMessage:{\n  source:" + req.Source.String() + ",\n  dest: " + req.Dest.String() + ",\n  message:" + req.Message + " }"
}


/* 
 * This is a useful utility to format the json packet to send requests
 * This tiny block is sort of important else you will end up sending blank messages.
 */
func GetAddToClusterMessage(source NodeInfo, dest NodeInfo, message string) (AddToClusterMessage){
    return AddToClusterMessage{
        Source: NodeInfo{
                NodeId: source.NodeId,
                NodeIpAddr: source.NodeIpAddr,
                Port: source.Port,
                },
        Dest: NodeInfo{
                NodeId: dest.NodeId,
                NodeIpAddr: dest.NodeIpAddr,
                Port: dest.Port,
                },
        Message: message,
    }
}