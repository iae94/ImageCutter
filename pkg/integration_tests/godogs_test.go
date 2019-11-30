package main

import (
	cfg "ImageCutter/pkg/config"
	logging "ImageCutter/pkg/logger"
	"ImageCutter/pkg/models"
	"ImageCutter/pkg/services/cutter"
	"fmt"
	"github.com/DATA-DOG/godog"
	"go.uber.org/zap"
	"log"
	"net/http"
)



type testCutterService struct {
	Config *cfg.CutterConfig
	Logger *zap.Logger
	Cutter *cutter.CutterService
	RemoteServer string
	CutterServer string
	tempInstance *models.Image
	//testUrls map[string]ScenarioData
	Scenario *ScenarioData
}

type ScenarioData struct {
	Name string
	Url string
	Code int
	Next *ScenarioData
}

func NewTestCutter() *testCutterService{

	// Read config
	config, err := cfg.ReadConfig()
	if err != nil {
		log.Fatalf("Reading cutter config give error: %v\n", err)
	}

	// Create logger
	logger, err := logging.CreateLogger(&config.Cutter.Logger)
	if err != nil {
		log.Fatalf("Creating cutter logger give error: %v\n", err)
	}

	// Create cutter instance
	cutterService, err := cutter.NewCutterService(logger, config)
	if err != nil {
		log.Fatalf("Creating cutter service instance give error: %v\n", err)
	}

	remoteServer := "http://nginx:80/static/"
	config.Cutter.Port = 5006
	cutterServer := fmt.Sprintf("http://cutter:%v", config.Cutter.Port)

	//testUrls:= map[string]ScenarioData{
	//	"Image is found in cache": ScenarioData{Url: fmt.Sprintf("%v/crop/300/400/%v/1.jpg", cutterService, remoteServer), Code: 200},
	//	"Remote server does not exist": ScenarioData{Url: fmt.Sprintf("%v/crop/300/400/%v/1.jpg", cutterService, "http://localMost:80/static/"), Code: 503},
	//	"Image is not found on remote server": ScenarioData{Url: fmt.Sprintf("%v/crop/300/400/%v/not_existing_image.jpg", cutterService, remoteServer), Code: 404},
	//	"Image has unsupported extension": ScenarioData{Url: fmt.Sprintf("%v/crop/300/400/%v/not_existing_image.jpg", cutterService, "some.exe"), Code: 422},
	//	"Remote server return error":  ScenarioData{Url: fmt.Sprintf("%v/crop/300/400/%v/not_existing_image.jpg", cutterService, "http://localhost:80/error/"), Code: 500},
	//}

	scenarios := &ScenarioData{
		Name: "", // for first s.BeforeScenario()
		Url:  "",
		Code: 0,
		Next: &ScenarioData{
			Name: "Image is found in cache",
			Url:  fmt.Sprintf("%v/crop/300/400/%v/1.jpg", cutterServer, remoteServer),
			Code: 200,
			Next: &ScenarioData{
				Name: "Remote server does not exist",
				Url:  fmt.Sprintf("%v/crop/300/400/%v/1.jpg", cutterServer, "http://somehost:80/static/"),
				Code: 503,
				Next: &ScenarioData{
					Name: "Image is not found on remote server",
					Url:  fmt.Sprintf("%v/crop/300/400/%v/not_existing_image.jpg", cutterServer, remoteServer),
					Code: 404,
					Next: &ScenarioData{
						Name: "Image has unsupported extension",
						Url:  fmt.Sprintf("%v/crop/300/400/%v/%v", cutterServer, remoteServer, "1.docx"),
						Code: 422,
						Next: &ScenarioData{
							Name: "Remote server return error",
							Url:  fmt.Sprintf("%v/crop/300/400/%v", cutterServer, "http://nginx:80/error/"),
							Code: 500,
							Next: nil,
						},
					},
				},
			},
		},
	}

	testCutter := &testCutterService{
		Logger: logger,
		Config: config,
		Cutter: cutterService,
		RemoteServer: remoteServer,
		CutterServer: cutterServer,
		Scenario: scenarios,
	}
	return testCutter
}

func (tc *testCutterService) moveScenario(interface{}) {
	fmt.Printf("Code was: %v\n", tc.Scenario.Code)
	tc.Scenario = tc.Scenario.Next
	fmt.Printf("Check scenario: %v\n", tc.Scenario.Name)
	fmt.Printf("This scenario url: %v\n", tc.Scenario.Url)
	fmt.Printf("This scenario code: %v\n", tc.Scenario.Code)

}



func (tc *testCutterService) clientMakeCorrectRequestsToImageOnRemoteServer() error {
	testUrl := tc.Scenario.Url

	//img, _, err := tc.Cutter.FetchImage(testUrl)
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("fetching test url: %v give error: %v", testUrl, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("fetching test url: %v return non 200 code: %v", testUrl, resp.StatusCode)
	}

	return nil
}

func (tc *testCutterService) cutterServiceShouldPutImageInToCache() error {
	//_, err := tc.Cutter.Cache.GetImageByUrl(tc.tempInstance.Url)
	testUrl := fmt.Sprintf("%v/cache/%v/%v", tc.CutterServer, tc.RemoteServer, "1.jpg")
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("fetching test url: %v give error: %v", testUrl, err)
	}
	if resp.StatusCode == 404 {
		return fmt.Errorf("image: %v not in cache after fetching", tc.tempInstance.Url)
	}
	return nil
}


func (tc *testCutterService) clientMakeRequestToNonExistServer() error {
	testUrl := tc.Scenario.Url
	//_, _, err := tc.Cutter.FetchImage(testUrl)
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("fetching test url: %v give error: %v", testUrl, err)
	}
	tc.Scenario.Code = resp.StatusCode
	return nil
}

func (tc *testCutterService) cutterServiceShouldReturnHttpCode(code int) error {
	if code != tc.Scenario.Code{
		return fmt.Errorf("return http code for scenario: %v must be %v but given: %v", tc.Scenario.Name, tc.Scenario.Code, code)
	}
	return nil
}

func (tc *testCutterService) clientMakeRequestToNonExistFile() error {
	testUrl := tc.Scenario.Url
	//_, _, err := tc.Cutter.FetchImage(testUrl)
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("fetching test url: %v give error: %v", testUrl, err)
	}
	tc.Scenario.Code = resp.StatusCode
	return nil
}

func (tc *testCutterService) clientMakeRequestToNonImageFile() error {
	testUrl := tc.Scenario.Url
	//_, _, err := tc.Cutter.FetchImage(testUrl)
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("fetching test url: %v give error: %v", testUrl, err)
	}
	tc.Scenario.Code = resp.StatusCode
	return nil
}

func (tc *testCutterService) remoteServerReturnErrorForCorrectRequest() error {
	testUrl := tc.Scenario.Url
	//_, _, err := tc.Cutter.FetchImage(testUrl)
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("fetching test url: %v give error: %v", testUrl, err)
	}
	tc.Scenario.Code = resp.StatusCode
	return nil
}


func FeatureContext(s *godog.Suite) {

	testCutter := NewTestCutter()

	s.BeforeScenario(testCutter.moveScenario)
	//s.AfterScenario(testCutter.moveScenario)

	s.Step(`^Client make correct requests to image on remote server$`, testCutter.clientMakeCorrectRequestsToImageOnRemoteServer)
	s.Step(`^Cutter service should put image in to cache$`, testCutter.cutterServiceShouldPutImageInToCache)
	s.Step(`^Client make request to non exist server$`, testCutter.clientMakeRequestToNonExistServer)
	s.Step(`^Cutter service should return (\d+) http code$`, testCutter.cutterServiceShouldReturnHttpCode)
	s.Step(`^Client make request to non exist file$`, testCutter.clientMakeRequestToNonExistFile)
	s.Step(`^Client make request to non image file$`, testCutter.clientMakeRequestToNonImageFile)
	s.Step(`^Remote server return error for correct request$`, testCutter.remoteServerReturnErrorForCorrectRequest)

}
