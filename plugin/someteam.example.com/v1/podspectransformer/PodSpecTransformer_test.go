// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	"sigs.k8s.io/kustomize/pkg/kusttest"
	"sigs.k8s.io/kustomize/plugin"
)

func TestPodSpecTransformer(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"someteam.example.com", "v1", "PodSpecTransformer")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	rm := th.LoadAndRunTransformer(`
apiVersion: someteam.example.com/v1
kind: PodSpecTransformer
metadata:
  name: notImportantHere
labels:
  env: production
spec:
  nodeSelector:
    type: production
`, `
apiVersion: v1
kind: Pod
metadata:
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: nginx
    env: production
  name: nginx
spec:
  containers:
  - image: nginx
    name: nginx
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: Pod
metadata:
  name: elasticsearch
spec:
  containers:
  - image: elasticsearch
    name: elasticsearch
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: nginx
    env: production
  name: nginx
spec:
  containers:
  - image: nginx
    name: nginx
  nodeSelector:
    type: production
`)
}

func TestPodSpecTransformerDeployments(t *testing.T) {
	tc := plugin.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"someteam.example.com", "v1", "PodSpecTransformer")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	rm := th.LoadAndRunTransformer(`
apiVersion: someteam.example.com/v1
kind: PodSpecTransformer
metadata:
  name: notImportantHere
labels:
  env: production
spec:
  nodeSelector:
    type: production
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: elasticsearch
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    metadata:
      labels:
        app: nginx
        env: production
    spec:
      containers:
      - image: nginx
        name: nginx
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: elasticsearch
spec:
  template:
    spec:
      containers:
      - image: elasticsearch
        name: elasticsearch
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    metadata:
      labels:
        app: nginx
        env: production
    spec:
      containers:
      - image: nginx
        name: nginx
      nodeSelector:
        type: production
`)
}
