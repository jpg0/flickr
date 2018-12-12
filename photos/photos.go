package photos

import (
	"github.com/jpg0/flickr"
	"time"
)

type PhotoInfo struct {
	Id           string `xml:"id,attr"`
	Secret       string `xml:"secret,attr"`
	Server       string `xml:"server,attr"`
	Farm         string `xml:"farm,attr"`
	DateUploaded string `xml:"dateuploaded,attr"`
	IsFavorite   bool   `xml:"isfavorite,attr"`
	License      string `xml:"license,attr"`
	// NOTE: one less than safety level set on upload (ie, here 0 = safe, 1 = moderate, 2 = restricted)
	//       while on upload, 1 = safe, 2 = moderate, 3 = restricted
	SafetyLevel    int    `xml:"safety_level,attr"`
	Rotation       int    `xml:"rotation,attr"`
	OriginalSecret string `xml:"originalsecret,attr"`
	OriginalFormat string `xml:"originalformat,attr"`
	Views          int    `xml:"views,attr"`
	Media          string `xml:"media,attr"`
	Title          string `xml:"title"`
	Description    string `xml:"description"`
	Visibility     struct {
		IsPublic bool `xml:"ispublic,attr"`
		IsFriend bool `xml:"isfriend,attr"`
		IsFamily bool `xml:"isfamily,attr"`
	} `xml:"visibility"`
	Dates struct {
		Posted           string `xml:"posted,attr"`
		Taken            string `xml:"taken,attr"`
		TakenGranularity string `xml:"takengranularity,attr"`
		TakenUnknown     string `xml:"takenunknown,attr"`
		LastUpdate       string `xml:"lastupdate,attr"`
	} `xml:"dates"`
	Permissions struct {
		PermComment string `xml:"permcomment,attr"`
		PermAdMeta  string `xml:"permadmeta,attr"`
	} `xml:"permissions"`
	Editability struct {
		CanComment string `xml:"cancomment,attr"`
		CanAddMeta string `xml:"canaddmeta,attr"`
	} `xml:"editability"`
	PublicEditability struct {
		CanComment string `xml:"cancomment,attr"`
		CanAddMeta string `xml:"canaddmeta,attr"`
	} `xml:"publiceditability"`
	Usage struct {
		CanDownload string `xml:"candownload,attr"`
		CanBlog     string `xml:"canblog,attr"`
		CanPrint    string `xml:"canprint,attr"`
		CanShare    string `xml:"canshare,attr"`
	} `xml:"usage"`
	Comments int `xml:"comments"`
	// Notes XXX: not handled yet
	// People XXX: not handled yet
	// Tags XXX: not handled yet
	// Urls XXX: not handled yet
}

type PhotoInfoResponse struct {
	flickr.BasicResponse
	Photo PhotoInfo `xml:"photo"`
}

type PhotoSearchResponse struct {
	flickr.BasicResponse
	PhotoList struct {
		Page    int         `xml:"page,attr"`
		Pages   int         `xml:"pages,attr"`
		PerPage int         `xml:"perpage,attr"`
		Total   int         `xml:"total,attr"`
		Photos  []PhotoInfo `xml:"photo"`
	} `xml:"photos"`
}

type PhotoAllContextsResponse struct {
	flickr.BasicResponse
	PhotoAllContexts
}

type PhotoAllContexts struct {
	Sets []struct {
		ID    int    `xml:"id,attr"`
		Title string `xml:"title,attr"`
	} `xml:"set" json:"set,omitempty"`
}

const (
	mysql_layout = "2006-01-02 15:04:05"
)

// Delete a photo from Flickr
// This method requires authentication with 'delete' permission.
func Delete(client *flickr.FlickrClient, id string) (*flickr.BasicResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.delete")
	client.Args.Set("photo_id", id)
	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Search for photos on Flickr
func Search(client *flickr.FlickrClient, authenticate bool, user_id string, min_upload_date time.Time, max_upload_date time.Time) (*PhotoSearchResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.search")
	client.Args.Set("user_id", user_id)
	zeroTime := &time.Time{}

	if !min_upload_date.Equal(*zeroTime) {
		client.Args.Set("min_upload_date", min_upload_date.Format(mysql_layout))
	}
	if !max_upload_date.Equal(*zeroTime) {
		client.Args.Set("max_upload_date", max_upload_date.Format(mysql_layout))
	}
	// sign the client for authentication and authorization
	if authenticate {
		client.OAuthSign()
	} else {
		client.ApiSign()
	}

	response := &PhotoSearchResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Get information about a Flickr photo
func GetInfo(client *flickr.FlickrClient, id string, secret string) (*PhotoInfoResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.getInfo")
	client.Args.Set("photo_id", id)
	if secret != "" {
		client.Args.Set("secret", secret)
	}
	client.OAuthSign()

	response := &PhotoInfoResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Get information about a Flickr photo
func GetAllContexts(client *flickr.FlickrClient, id string, secret string) (*PhotoAllContextsResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.getAllContexts")
	client.Args.Set("photo_id", id)
	if secret != "" {
		client.Args.Set("secret", secret)
	}
	client.OAuthSign()

	response := &PhotoAllContextsResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}

// Set date posted and date taken on a Flickr photo
// datePosted and dateTaken are optional and may be set to ""
func SetDates(client *flickr.FlickrClient, id string, datePosted string, dateTaken string) (*flickr.BasicResponse, error) {
	client.Init()
	client.EndpointUrl = flickr.API_ENDPOINT
	client.HTTPVerb = "POST"
	client.Args.Set("method", "flickr.photos.setDates")
	client.Args.Set("photo_id", id)
	if datePosted != "" {
		client.Args.Set("date_posted", datePosted)
	}
	if dateTaken != "" {
		client.Args.Set("date_taken", dateTaken)
	}
	client.OAuthSign()

	response := &flickr.BasicResponse{}
	err := flickr.DoPost(client, response)
	return response, err
}
