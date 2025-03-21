package cloutility

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Node struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Href        string `json:"href"`
	TsmName     string `json:"tsmName"`
	TsmPassword string `json:"tsmPassword"`
}

type NodeData struct {
	OperatingSystem NodeOperatingSystem
	Type            NodeType
	Domain          Domain
	ClientOptionSet ClientOptionSet
	Contact         string
	CpuCount        int
}

func (c *AuthenticatedClient) CreateNode(bUnitID, consumerID, osType, clientType, domain, clientOptionSet int, contact string) (Node, error) {
	var (
		newNode  Node
		nodedata NodeData
	)

	// validate the base url to create the endpoint
	endpoint := "/v1/bunits/" + fmt.Sprintf("%d", bUnitID) + "/consumers/" + fmt.Sprintf("%d", consumerID) + "/node"

	// Assign data
	nodedata.Contact = contact
	nodedata.Domain.ID = domain
	nodedata.OperatingSystem.ID = osType
	nodedata.Type.ID = clientType
	nodedata.CpuCount = 1
	nodedata.ClientOptionSet.ID = clientOptionSet

	payload, err := json.Marshal(nodedata)
	if err != nil {
		return Node{}, fmt.Errorf("failed to encode json payload: %s", err)
	}

	resp, err := c.apiRequest(endpoint, http.MethodPost, payload)
	if err != nil {
		return Node{}, fmt.Errorf("failed to create node: %s", err)
	}

	if err := json.Unmarshal([]byte(resp), &newNode); err != nil {
		return Node{}, fmt.Errorf("failed to decode nodedata: %s", err)
	}

	return newNode, nil
}

func (c *AuthenticatedClient) DeleteNode(bUnitID, consumerID int) (Node, error) {
	var node Node

	endpoint := "/v1/bunits/" + fmt.Sprintf("%d", bUnitID) + "/consumers/" + fmt.Sprintf("%d", consumerID) + "/node?deleteAssociations=True"
	resp, err := c.apiRequest(endpoint, http.MethodDelete, nil)
	if err != nil {
		return Node{}, fmt.Errorf("error requesting nodedata: %s", err)
	}
	if err := json.Unmarshal([]byte(resp), &node); err != nil {
		return Node{}, fmt.Errorf("failed to decode nodedata: %s", err)
	}
	return node, nil
}

func (c *AuthenticatedClient) GetNode(bUnitID, consumerID int) (Node, error) {
	var node Node

	endpoint := "/v1/bunits/" + fmt.Sprintf("%d", bUnitID) + "/consumers/" + fmt.Sprintf("%d", consumerID) + "/node"
	resp, err := c.apiRequest(endpoint, http.MethodGet, nil)
	if err != nil {
		return Node{}, fmt.Errorf("error requesting nodedata: %s", err)
	}
	if err := json.Unmarshal([]byte(resp), &node); err != nil {
		return Node{}, fmt.Errorf("failed to decode nodedata: %s", err)
	}
	return node, nil
}

func (c *AuthenticatedClient) ActivateNode(bUnitID, consumerID int) (Node, error) {
	var (
		node Node
		err  error
	)

	endpoint := "/v1/bunits/" + fmt.Sprintf("%d", bUnitID) + "/consumers/" + fmt.Sprintf("%d", consumerID) + "/node/spname"
	node.TsmName, err = c.apiRequest(endpoint, http.MethodGet, nil)
	if err != nil {
		return Node{}, fmt.Errorf("error retrieving nodename: %s", err)
	}

	endpoint = "/v1/bunits/" + fmt.Sprintf("%d", bUnitID) + "/consumers/" + fmt.Sprintf("%d", consumerID) + "/node/activate?tsmName=" + node.TsmName
	_, err = c.apiRequest(endpoint, http.MethodGet, nil)
	if err != nil {
		return Node{}, fmt.Errorf("error activating node: %s", err)
	}

	node, err = c.GetNode(bUnitID, consumerID)
	if err != nil {
		return Node{}, fmt.Errorf("error retrieving node: %s", err)
	}

	return node, nil
}
