package main

import (
	"encoding/json"
	"fmt"
	"github.com/niemeyer/qml"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()
	servers := &Servers{}
	engine.Context().SetVar("servers", servers)
	component, err := engine.LoadFile("jenkinsmonitor.qml")
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	window.Show()
	servers.Add(&Server{Address: "http://jenkins.qa.ubuntu.com"})
	window.Wait()
	return nil
}

/// Servers: A model for all the servers.
type Servers struct {
	List []*Server
	Len  int
}

func (servers *Servers) Add(svr *Server) {
	servers.List = append(servers.List, svr)
	servers.Len = len(servers.List)
	qml.Changed(servers, &servers.Len)
	go svr.GetJobDetails()
}

func (servers *Servers) Server(index int) *Server {
	return servers.List[index]
}

/// Server - an object representing a single jenkins server.
type Server struct {
	Address string
	Port    int
	Jobs    Jobs
}

func (server *Server) GetJobDetails() {
	full_url := strings.TrimSuffix(server.Address, "/")
	if server.Port != 80 && server.Port != 0 {
		full_url += ":" + strconv.Itoa(server.Port)
	}
	full_url += "/api/json"
	// fmt.Printf("%s ", full_url)
	resp, err := http.Get(full_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// decode json data into ServerData
	type ServerData struct {
		Jobs []Job
	}
	var s ServerData
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatal(err)
	}
	// add found jobs to server model
	for i := 0; i < len(s.Jobs); i++ {
		// TODO: Why does adding 'go' here make this faster? surely I'm already
		// asynchronous?
		go server.Jobs.Add(s.Jobs[i])
	}
}

type Jobs struct {
	list []Job
	Len  int
}

func (jobs *Jobs) Add(job Job) {
	// fmt.Print(job.Color)
	jobs.list = append(jobs.list, job)
	jobs.Len = len(jobs.list)
	// TODO: Do we need this here?
	qml.Changed(jobs, &jobs.Len)
}

func (jobs *Jobs) Job(index int) Job {
	return jobs.list[index]
}

type Job struct {
	Name  string
	Url   string
	Color string
}

func (job Job) RenderColor() color.RGBA {
	switch string(job.Color) {
	case "blue":
		return color.RGBA{108, 166, 210, 128}
	case "red":
		return color.RGBA{195, 0, 3, 128}
	case "yellow":
		return color.RGBA{223, 195, 0, 128}
	}
	return color.RGBA{128, 128, 128, 128}
}
