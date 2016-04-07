/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file is generated by client-gen with arguments: --clientset-name=iot --input=[api/v1]

package v1

import (
	api "k8s.io/kubernetes/pkg/api"
	v1 "k8s.io/kubernetes/pkg/api/v1"
	watch "k8s.io/kubernetes/pkg/watch"
)

// SensorAccessesGetter has a method to return a SensorAccessInterface.
// A group's client should implement this interface.
type SensorAccessesGetter interface {
	SensorAccesses(namespace string) SensorAccessInterface
}

// SensorAccessInterface has methods to work with SensorAccess resources.
type SensorAccessInterface interface {
	Create(*v1.SensorAccess) (*v1.SensorAccess, error)
	Update(*v1.SensorAccess) (*v1.SensorAccess, error)
	Delete(name string, options *api.DeleteOptions) error
	DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error
	Get(name string) (*v1.SensorAccess, error)
	List(opts api.ListOptions) (*v1.SensorAccessList, error)
	Watch(opts api.ListOptions) (watch.Interface, error)
	SensorAccessExpansion
}

// sensorAccesses implements SensorAccessInterface
type sensorAccesses struct {
	client *CoreClient
	ns     string
}

// newSensorAccesses returns a SensorAccesses
func newSensorAccesses(c *CoreClient, namespace string) *sensorAccesses {
	return &sensorAccesses{
		client: c,
		ns:     namespace,
	}
}

// Create takes the representation of a sensorAccess and creates it.  Returns the server's representation of the sensorAccess, and an error, if there is any.
func (c *sensorAccesses) Create(sensorAccess *v1.SensorAccess) (result *v1.SensorAccess, err error) {
	result = &v1.SensorAccess{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sensoraccesses").
		Body(sensorAccess).
		Do().
		Into(result)
	return
}

// Update takes the representation of a sensorAccess and updates it. Returns the server's representation of the sensorAccess, and an error, if there is any.
func (c *sensorAccesses) Update(sensorAccess *v1.SensorAccess) (result *v1.SensorAccess, err error) {
	result = &v1.SensorAccess{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sensoraccesses").
		Name(sensorAccess.Name).
		Body(sensorAccess).
		Do().
		Into(result)
	return
}

// Delete takes name of the sensorAccess and deletes it. Returns an error if one occurs.
func (c *sensorAccesses) Delete(name string, options *api.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sensoraccesses").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sensorAccesses) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sensoraccesses").
		VersionedParams(&listOptions, api.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the sensorAccess, and returns the corresponding sensorAccess object, and an error if there is any.
func (c *sensorAccesses) Get(name string) (result *v1.SensorAccess, err error) {
	result = &v1.SensorAccess{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sensoraccesses").
		Name(name).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SensorAccesses that match those selectors.
func (c *sensorAccesses) List(opts api.ListOptions) (result *v1.SensorAccessList, err error) {
	result = &v1.SensorAccessList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sensoraccesses").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sensorAccesses.
func (c *sensorAccesses) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("sensoraccesses").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}
