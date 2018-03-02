#!/bin/bash

name=$RANDOM
url='http://localhost:9093/api/v1/alerts'

echo "firing up alert $name"

curl -XPOST $url -d '[{
"labels": {
"alertname": "'$name'",
"service": "my-service",
"severity": "warning",
"instance": "'$name.example.net'"
},
"annotations": {
"deployment": "prometheus-kubernetes-autoscaler-go",
"action": "down"}
}]'
