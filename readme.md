# Prometheus Kubernetes Autoscaler
Scale Kubernetes deployments based on Prometheus alerts. 

## How does it work?
The autoscaler is registered as a webhook alertmanager receiver. When an alert is forwarded to this receiver the alertmanager will contact the autoscaler.  
The notifications send by the alertmanager need to contain these annotations: 

```yaml
alertmanager-config.yaml:
global:
    resolve_timeout: 5m

route:
    group_by: ['service']
    group_wait: 10s
    group_interval: 10s 
    repeat_interval: 1h
    receiver: 'webhook'
receivers:
- name: 'webhook'
    webhook_configs:
    - url: 'http://prometheus-kubernetes-autoscaler-go/alert
```

| name | description |
| --- | --- | --- |
| deployment|name of the deployment that should be scaled |
| action|scaling that should be performed. Either `up` or `down` |

## Installation
To run this project you need the following tools:
* [Helm](https://github.com/kubernetes/helm/)
* [Draft](https://github.com/Azure/draft/)
* [Helmfile](https://github.com/roboll/helmfile)


1. Run `helmfile sync` to install the alertmanager
2. Run `draft up` to install the autoscaler
3. Port-forward the alertmanager api to `localhost:9093`.
4. Run `generate-alert.sh` to generate an alert and watch the scaling happen

## License
MIT License

Copyright (c) 2018 Lukas Eichler

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.