package element

import (
	"fmt"
	"github.com/lukasnoeller/smugmug-backupper/internal/auth"
	"github.com/lukasnoeller/smugmug-backupper/internal/misc"
	"time"
)

type ChildNodesResponse struct {
	Code     int      `json:"Code"`
	Message  string   `json:"Message"`
	Options  Options  `json:"Options"`
	Response NodeData `json:"Response"`
}

type NodeData struct {
	Uri          string     `json:"Uri"`
	Locator      string     `json:"Locator"`
	LocatorType  string     `json:"LocatorType"`
	EndpointType string     `json:"EndpointType"`
	Node         []NodeItem `json:"Node"` // Array of child nodes!
	Pages        Pages      `json:"Pages"`
}

type NodeItem struct {
	NodeID          string          `json:"NodeID"`
	Name            string          `json:"Name"`
	Type            string          `json:"Type"` // "Folder", "Album", or "Page"
	UrlName         string          `json:"UrlName"`
	UrlPath         string          `json:"UrlPath"`
	WebUri          string          `json:"WebUri"`
	Uri             string          `json:"Uri"`
	HasChildren     bool            `json:"HasChildren"`
	IsRoot          bool            `json:"IsRoot"`
	Privacy         string          `json:"Privacy"`
	SecurityType    string          `json:"SecurityType"`
	SortMethod      string          `json:"SortMethod"`
	SortDirection   string          `json:"SortDirection"`
	SortIndex       int             `json:"SortIndex"`
	DateAdded       time.Time       `json:"DateAdded"`
	DateModified    time.Time       `json:"DateModified"`
	FormattedValues FormattedValues `json:"FormattedValues"`
	Uris            NodeUris        `json:"Uris"`
}

type FormattedValues struct {
	Name struct {
		HTML string `json:"html"`
	} `json:"Name"`
	Description struct {
		HTML string `json:"html"`
		Text string `json:"text"`
	} `json:"Description"`
}

type NodeUris struct {
	ChildNodes     EndpointReference `json:"ChildNodes"`
	ParentNode     EndpointReference `json:"ParentNode"`
	ParentNodes    EndpointReference `json:"ParentNodes"`
	FolderByID     EndpointReference `json:"FolderByID"`
	User           EndpointReference `json:"User"`
	NodeCoverImage EndpointReference `json:"NodeCoverImage"`
	HighlightImage EndpointReference `json:"HighlightImage"`
	NodeComments   EndpointReference `json:"NodeComments"`
	NodeGrants     EndpointReference `json:"NodeGrants"`
	MoveNodes      EndpointReference `json:"MoveNodes"`
}

type EndpointReference struct {
	Uri            string `json:"Uri"`
	Locator        string `json:"Locator,omitempty"`
	LocatorType    string `json:"LocatorType,omitempty"`
	UriDescription string `json:"UriDescription,omitempty"`
	EndpointType   string `json:"EndpointType,omitempty"`
}

type Pages struct {
	Total          int    `json:"Total"`
	Start          int    `json:"Start"`
	Count          int    `json:"Count"`
	RequestedCount int    `json:"RequestedCount"`
	FirstPage      string `json:"FirstPage"`
	LastPage       string `json:"LastPage"`
	NextPage       string `json:"NextPage"`
}

type Options struct {
	Methods []string `json:"Methods"`
}

func GetChildNodes(nodeUri string) ChildNodesResponse {
	nodeRequestUrl := fmt.Sprintf("https://api.smugmug.com/api/v2/node/%s!children?count=100", nodeUri)
	body, err := auth.AuthenticatedRequest("GET", nodeRequestUrl)
	if err != nil {
		fmt.Println("error ocurred during request: ", err)
	}
	var childNodesResponse ChildNodesResponse

	err = misc.Parse(body, &childNodesResponse)
	if err != nil {
		fmt.Println("error ocurred during unmashalling of child nodes ", err)
	}
	return childNodesResponse

}

func GetElements(cnr ChildNodesResponse) {
	for _, element := range cnr.Response.Node {
		switch element.Type {

		case "Folder":
			fmt.Println("Folder: ", element.Name)
		case "Album":
			fmt.Println("Album: ", element.Name)
			fmt.Println("uri: ", element.Uri)
		case "Page":
			fmt.Println("Page: ", element.Name)

		}
	}
}
