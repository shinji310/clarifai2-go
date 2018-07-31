package clarifai2

import (
	"fmt"
	"encoding/json"
	"encoding/base64"
	"io/ioutil"
)

/*
Data Structure Definition for JSON
*/

type Bytes []byte

type Input struct {
	Id string `json:"id,omitempty"`
	Data *Data`json:"data,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
	Base string `json:"base64,omitempty"`
	Crop []float64 `json:"crop,omitempty"`
}

type Video struct {
	URL string `json:"url,omitempty"`
	Base string `json:"base64,omitempty"`
}

type OutputConfig struct {
	SelectConcepts []Concept `json:"select_concepts,omitempty"`
	MaxConcepts int `json:"max_concepts,omitempty"`
	MinValue int `json:"min_value",omitempty"`
	ConceptsMutuallyExclusive bool `json:"concepts_mutually_exclusive,omitempty"`
	ClosedEnvironment bool `json:"closed_environment,omitempty"`
}

type OutputInfo struct {
	Message string `json:"message,omitempty"`
	Type string `json:"type,omitempty"`
	TypeExt string `json:"type_ext,omitempty"`
	OutputConfig *OutputConfig `json:"output_config,omitempty"`
	Data Data `json:"data,omitempty"`
}

type Model struct {
	Name string `json:"name,omitempty"`
	Id string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	AppId string `json:"app_id,omitempty"`
	OutputInfo *OutputInfo `json:"output_info"`
	ModelVersion *ModelVersion `json:"model_version,omitempty"`
}

type ModelVersion struct {
	Id string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Status Status `json:"status,omitempty"`
}

type Concept struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	AppId string `json:"app_id,omitempty"`
	Value float32 `json:"value,omitempty"`
}

type Status struct {
	code string `json:"code"`
	description string `json:"description"`
}

type GeoLimit struct {
	Type string `json:"type"`
	Value int `json:"value"`
}

type Geo struct {
	GeoPoint struct {
		Longitude int `json:"longitude"`
		Latitude int `json:"latitude"`
	} `json:"geo_point"`
	GeoLimit *GeoLimit `json:"geo_limit,omitempty"`
}

type Data struct {
	Image *Image `json:"image,omitempty"`
	Video *Video `json:"video,omitempty"`
	Concepts []Concept `json:"concepts,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Geo *Geo `json:"geo,omitempty"`
}

type Output struct {
	Data *Data `json:"data,omitempty"`
	Input *Input `json:"input,omitempty"`
}

type Ands struct {
	Input *Input `json:"input,omitempty"`
	Output *Output `json:"output,omitempty"`
}

type SearchReq struct {
	Query struct {
		Ands []Ands `json:"ands"`
	} `json:"query"`
}

type InputReq struct {
	Input []Input `json:"inputs,omitempty"`
	Action string `json:"action,omitempty"`
	DeleteAll bool `json:"delete_all,omitempty"`
}

type ModelReq struct {
	Models []Model `json:"models,omitempty"`
	Concepts []Concept `json:"concepts,omitempty"`
	Action string `json:"action,omitempty"`
}
		
type PredictReq struct {
	Input []Input `json:"inputs"`
	Model *Model `json:"model,omitempty"`
}

type Response struct {
	Status Status `json:"status"`
	Output []struct {
		Id string `json:"id"`
		Data Data `json:"data"`
		CreatedAt string `json:"created_at"`
		Model Model `json:"model"`
		Input Input `json:"input"`
		Status Status `json:"status"`
	} `json:"outputs"`
	Input []struct {
		Id string `json:"id"`
		Data Data `json:"data"`
		CreatedAt string`json:"created_at"`
		ModifiedAt string `json:"modified_at"`
		Status Status `json:"status"`
	} `json:"inputs"`
}

/*
Predict - analyze images to find tags
*/

func (client *Client) Predict(req PredictReq, model string) (*Response, error) {
	// build endpoint
	endpoint := "models/" + model + "/outputs"
	
	res, err := client.commonHTTPRequest(req, endpoint, "POST", false)
	if err != nil {
		return nil, err
	}
	prediction := new(Response)
	err = json.Unmarshal(res, prediction)

	return prediction, err
}

// Send images by specifying external URLs

func (client *Client) PredictByUrls(url []string, model string) (*Response, error) {
	req := PredictReq{Input: make([]Input, len(url))}
	for i := range url {
		req.Input[i].Data = new(Data)
		req.Input[i].Data.Image = new(Image)
		req.Input[i].Data.Image.URL = url[i]
	}
	return client.Predict(req, model)
}

// Send images as binary data (i.e. []byte)

func (client *Client) PredictByBytes(data []Bytes, model string) (*Response, error) {
	// encode by base64
	req := PredictReq{Input: make([]Input, len(data))}
	for i := range data {
		req.Input[i].Data = new(Data)
		req.Input[i].Data.Image = new(Image)
		req.Input[i].Data.Image.Base = base64.StdEncoding.EncodeToString(data[i])
	}
	return client.Predict(req, model)
}

// Send image files stored locally

func (client *Client) PredictByFiles(path []string, model string) (*Response, error) {
	body := make([]Bytes, len(path))
	for i := range path {
		b, err := ioutil.ReadFile(path[i])
		body[i] = b
		if err != nil {
			return nil, err
		}
	}
	return client.PredictByBytes(body, model)
}

// Send videos by specifying external URLs

func (client *Client) PredictVideoByUrls(url []string, model string) (*Response, error) {
	req := PredictReq{Input: make([]Input, len(url))}
	for i := range url {
		req.Input[i].Data = new(Data)
		req.Input[i].Data.Video = new(Video)
		req.Input[i].Data.Video.URL = url[i]
	}
	return client.Predict(req, model)
}

// Send videos as binary data (i.e. [] byte)

func (client *Client) PredictVideoByBytes(data []Bytes, model string) (*Response, error) {
	// encode by base64
	req := PredictReq{Input: make([]Input, len(data))}
	for i := range data {
		req.Input[i].Data = new(Data)
		req.Input[i].Data.Video = new(Video)
		req.Input[i].Data.Video.Base = base64.StdEncoding.EncodeToString(data[i])
	}
	return client.Predict(req, model)
}

// Send video files stored locally

func (client *Client) PredictVideoByFiles(path []string, model string) (*Response, error) {
	body := make([]Bytes, len(path))
	for i := range path {
		b, err := ioutil.ReadFile(path[i])
		body[i] = b
		if err != nil {
			return nil, err
		}
	}
	return client.PredictVideoByBytes(body, model)
}


/*
Input - Store images to Clarifai server
*/

func (client *Client) Input(req InputReq) (*Response, error) {
	// build endpoint
	endpoint := "inputs"
	
	res, err := client.commonHTTPRequest(req, endpoint, "POST", false)
	if err != nil {
		return nil, err
	}
	response := new(Response)
	err = json.Unmarshal(res, response)

	return response, err
}

// Send images specified by URLs to Clarifai

func (client *Client) InputByUrls(url []string) (*Response, error) {
	req := InputReq{Input: make([]Input, len(url))}
	for i := range url {
		req.Input[i].Data = new(Data)
		req.Input[i].Data.Image = new(Image)
		req.Input[i].Data.Image.URL = url[i]
	}
	return client.Input(req)
}

// Send images stored locally to Clarifai

func (client *Client) InputByBytes(data []Bytes) (*Response, error) {
	// encode by base64
	req := InputReq{Input: make([]Input, len(data))}
	for i := range data {
		req.Input[i].Data = new(Data)
		req.Input[i].Data.Image = new(Image)
		req.Input[i].Data.Image.Base = base64.StdEncoding.EncodeToString(data[i])
	}
	return client.Input(req)
}

// Send local files onto Clarifai server

func (client *Client) InputByFiles(path []string) (*Response, error) {
	body := make([]Bytes, len(path))
	for i := range path {
		b, err := ioutil.ReadFile(path[i])
		body[i] = b
		if err != nil {
			return nil, err
		}
	}
	return client.InputByBytes(body)
}

// Delete all the images on Clarifai server

func (client *Client) DeleteAllInputs() (*Response, error) {
	req := InputReq{DeleteAll:true}

	endpoint := "inputs"
	
	res, err := client.commonHTTPRequest(req, endpoint, "DELETE", false)
	if err != nil {
		return nil, err
	}
	response := new(Response)
	err = json.Unmarshal(res, response)

	return response, err
}


/*
Search - find an image that matches a criteria
*/

// need more work...


/* 
Utility functions
*/

// Obtain top five tags

func (client *Client) TopFive(prediction *Response) (tags []string) {
	tags = make([]string, 5)
	
	for i := 0; i < 5; i++ {
		tags[i] = prediction.Output[0].Data.Concepts[i].Name
		fmt.Println(tags[i])
	}
	return tags
}

