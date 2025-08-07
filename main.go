package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tercnem/jsonserver/config"

	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
)

func init() {
	// Initialize the app responder and validator.

	// Initialize the Fiber configuration.
	config.InitFiberConfig()
}

var responseEndpoint = &ResponseEndpoint{}

func main() {
	name := flag.String("config", "", "api.json")

	flag.Parse()
	// This is the main entry point of the application.
	// You can add your application logic here.
	bData, err := os.ReadFile(*name)
	if err != nil {
		panic(err)
	}
	config.ApiConfig = &config.APIConfig{}
	err = json.Unmarshal(bData, config.ApiConfig)
	if err != nil {
		panic(err)
	}

	responseEndpoint.LoadEndpoints(*name)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					fmt.Println("Load config json")
					responseEndpoint.LoadEndpoints(*name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(*name)
	if err != nil {
		panic(err)
	}

	appfiber := fiber.New(config.FiberConfig)
	//appfiber.Use(recover.New(recover.Config{EnableStackTrace: true}))
	appfiber.Use(ResponseHandler)

	go func() {
		if err := appfiber.Listen(":" + fmt.Sprint(config.ApiConfig.Port)); err != nil {
			log.Panic(err)
		}
	}()
	c := make(chan os.Signal, 1) // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = appfiber.Shutdown()

}

func ResponseHandler(c *fiber.Ctx) error {
	endpoint, err := responseEndpoint.ReadEndpoints(c.Path(), c.Method())
	if err != nil {
		return err
	}
	for key, value := range endpoint.HeaderResponse {
		c.Set(key, fmt.Sprintf("%v", value))
	}
	c.Status(endpoint.Status)
	c.JSON(endpoint.JsonResponse)
	// This is a sample handler function.
	// You can add your logic here.
	return nil
}

type Endpoint struct {
	Method         string         `json:"method"`
	Status         int            `json:"status"`
	Path           string         `json:"path"`
	HeaderResponse map[string]any `json:"headerResponse,omitempty"`
	JsonResponse   any            `json:"JsonResponse"`
}
type DataEndPoint struct {
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

type ResponseEndpoint struct {
	Endpoints map[string]map[string]*Endpoint
	sync.Mutex
}

func (e *ResponseEndpoint) LoadEndpoints(apiPath string) {

	e.Lock()
	defer e.Unlock()
	bData, err := os.ReadFile(apiPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	dataEndPoint := &DataEndPoint{}
	err = json.Unmarshal(bData, dataEndPoint)
	if err != nil {
		fmt.Println("Error parsing json file:", err)
	}
	if e.Endpoints == nil {
		e.Endpoints = make(map[string]map[string]*Endpoint)
	}
	for _, endpoint := range dataEndPoint.Endpoints {
		if _, ok := e.Endpoints[endpoint.Path]; !ok {
			e.Endpoints[endpoint.Path] = make(map[string]*Endpoint)
		}
		e.Endpoints[endpoint.Path][endpoint.Method] = &endpoint
	}
}

func (e *ResponseEndpoint) ReadEndpoints(path string, method string) (*Endpoint, error) {

	e.Lock()
	defer e.Unlock()

	if res, ok := e.Endpoints[path][method]; ok {
		return res, nil
	}

	return nil, errors.New("endpoint not found")

}

// "port": 3000,
//   "logLevel": "info",
//   "logFormat": "text",
//   "endpoints": [
