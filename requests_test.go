package clarifai2

import "testing"

// Needs to put in your own API key
const (
	MyAPIKey = "Your_API_key_here"
	GeneralModel = "aaa03c23b3724a16a56b629203edc62c"
	ImageFile = "TEST.jpg"
	VideoFile = "TEST.mp4"
	ImageUrl = "https://samples.clarifai.com/metro-north.jpg"
	VideoUrl = "https://samples.clarifai.com/beer.mp4"
)

// Sanity check for available functions

func TestDeleteAllInputs(t *testing.T) {
	client := NewClient(MyAPIKey)

	_, err := client.DeleteAllInputs()

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPredictByUrls(t *testing.T) {
	client := NewClient(MyAPIKey)

	urls := []string{ImageUrl}
	_, err := client.PredictByUrls(urls, GeneralModel)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPredictByFiles(t *testing.T) {
	client := NewClient(MyAPIKey)

	files := []string{ImageFile}
	_, err := client.PredictByFiles(files, GeneralModel)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPredictVideoByUrls(t *testing.T) {
	client := NewClient(MyAPIKey)

	urls := []string{VideoUrl}
	_, err := client.PredictVideoByUrls(urls, GeneralModel)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPredictVideoByFiles(t *testing.T) {
	client := NewClient(MyAPIKey)

	files := []string{VideoFile}
	_, err := client.PredictVideoByFiles(files, GeneralModel)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestInputByUrls(t *testing.T) {
	client := NewClient(MyAPIKey)

	urls := []string{ImageUrl}
	_, err := client.InputByUrls(urls)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestInputByFiles(t *testing.T) {
	client := NewClient(MyAPIKey)

	files := []string{ImageFile}
	_, err := client.InputByFiles(files)

	if err != nil {
		t.Errorf(err.Error())
	}
}
