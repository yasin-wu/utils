package weed

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	js "github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"
)

const (
	defaultMaster = "http://127.0.0.1:9333"
)

const (
	assign = "/dir/assign"
)

type Client struct {
	master string
}

func New(master string) *Client {
	if master == "" {
		master = defaultMaster
	}
	if strings.HasSuffix(master, "/") {
		master = master[0 : len(master)-1]
	}
	return &Client{master: master}
}

func (this *Client) Upload(file *os.File) (*js.Json, error) {
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("%s%s", this.master, assign))
	if err != nil {
		return nil, err
	}
	respJson, err := js.NewJson(resp.Body())
	if err != nil {
		return nil, err
	}
	if respJson == nil {
		return nil, errors.New("response body is err")
	}

	fid := respJson.Get("fid").MustString()
	volumeId := strings.Split(fid, ",")[0]
	url := respJson.Get("url").MustString()
	targetUrl := fmt.Sprintf("http://%s/%s", url, fid)

	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	formFile, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}
	contentType := writer.FormDataContentType()
	writer.Close()

	response, err := http.Post(targetUrl, contentType, buffer)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyBuffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyJson, err := js.NewJson(bodyBuffer)
	if err != nil {
		return nil, err
	}
	if _, ok := bodyJson.CheckGet("error"); ok {
		return nil, errors.New(bodyJson.Get("error").MustString())
	}
	bodyJson.Set("fid", fid)
	bodyJson.Set("volume_id", volumeId)
	bodyJson.Set("url", targetUrl)
	bodyJson.Set("file_type", strings.Replace(path.Ext(path.Base(file.Name())), ".", "", -1))
	return bodyJson, nil
}

func (this *Client) Delete(fileUrl string) error {
	client := resty.New()
	resp, err := client.R().Delete(fileUrl)
	if err != nil {
		return err
	}
	respJson, err := js.NewJson(resp.Body())
	if err != nil {
		return err
	}
	if respJson == nil {
		return errors.New("response body is err")
	}
	return nil
}
