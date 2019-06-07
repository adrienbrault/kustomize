// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	Labels map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec   map[string]interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
}

var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) (err error) {
	p.Labels = nil
	p.Spec = nil
	return yaml.Unmarshal(c, p)
}

func (p *plugin) Transform(m resmap.ResMap) error {
	for _, resource := range m.Resources() {
		resource.GetGvk()

		podSpecLabels := getResourcePodSpecLabels(resource)
		if !labelsMatch(podSpecLabels, p.Labels) {
			continue
		}

		updateResourcePodSpec(resource, p.Spec)
	}

	return nil
}

func labelsMatch(labels map[string]interface{}, checkLabels map[string]string) bool {
	for k, v := range checkLabels {
		if labels[k] != v {
			return false
		}
	}

	return true
}

func getResourcePodSpecLabels(resource *resource.Resource) map[string]interface{} {
	gvk := resource.GetGvk()

	if gvk.Kind == "Pod" {
		metadata, _ := resource.Kunstructured.Map()["metadata"].(map[string]interface{})
		labels, _ := metadata["labels"].(map[string]interface{})

		return labels
	}

	if gvk.Kind == "Deployment" {
		return getMapAtPath(resource.Kunstructured.Map(), []string{"spec", "template", "metadata", "labels"})
	}

	return map[string]interface{}{}
}

func updateResourcePodSpec(resource *resource.Resource, specToMerge map[string]interface{}) {
	gvk := resource.GetGvk()

	resourceMap := resource.Kunstructured.Map()

	if gvk.Kind == "Pod" {
		resourceMapSpec, _ := resourceMap["spec"].(map[string]interface{})
		for k, v := range specToMerge {
			resourceMapSpec[k] = v
		}
	}

	if gvk.Kind == "Deployment" {
		templateSpec := getMapAtPath(resource.Kunstructured.Map(), []string{"spec", "template", "spec"})

		for k, v := range specToMerge {
			templateSpec[k] = v
		}
	}

	resource.Kunstructured.SetMap(resourceMap)
}

func getMapAtPath(keyValues map[string]interface{}, path []string) map[string]interface{} {
	for _, key := range path {
		keyValues, _ = keyValues[key].(map[string]interface{})
	}

	return keyValues
}
