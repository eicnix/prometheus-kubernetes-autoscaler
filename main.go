package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ericchiang/k8s"

	appsv1 "github.com/ericchiang/k8s/apis/apps/v1beta1"
)

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt,omitempty"`
	EndsAt       string            `json:"endsAt,omitempty"`
	GeneratorURL string            `json:"generatorURL,omitempty"`
}

type Notification struct {
	Receiver string  `json:"receiver"`
	Status   string  `json:"status"`
	Alerts   []Alert `json:"alerts"`

	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
	Version     string `json:"version"`
}

func main() {
	client, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			log.Printf("Warning: Received a non POST request: %v %v\n", r.Method, r.URL)
			w.WriteHeader(405)
			w.Write([]byte("Only POST requests are allowed."))

			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error: Requestbody cannot be read: %v\n", err)
			w.WriteHeader(400)
			return
		}

		notification := Notification{}
		if err = json.Unmarshal(body, &notification); err != nil {
			log.Printf("Error: Request body cannot be parsed to JSON: %v\n", err)
			w.WriteHeader(400)
			return
		}

		deployment := notification.CommonAnnotations["deployment"]
		action := notification.CommonAnnotations["action"]

		log.Printf("Received notification to scale %v %v", deployment, action)

		var k8sDeployment appsv1.Deployment
		if err := client.Get(context.Background(), "default", deployment, &k8sDeployment); err != nil {
			log.Fatal(err)
		}

		newReplicaCount := k8sDeployment.Spec.GetReplicas()
		switch action {
		case "up":
			newReplicaCount += int32(1)
			break
		case "down":
			newReplicaCount -= int32(1)
			break
		}
		k8sDeployment.Spec.Replicas = &newReplicaCount
		if err := client.Update(context.Background(), &k8sDeployment); err != nil {
			log.Fatal(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type appHandler func(http.ResponseWriter, *http.Request) *error

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(r.RequestURI)
}
