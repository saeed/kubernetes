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

// This file is generated by client-gen with arguments: --clientset-name=internalclientset --input=[api/]

package fake

import (
	api "k8s.io/kubernetes/pkg/api"
	core "k8s.io/kubernetes/pkg/client/testing/core"
	labels "k8s.io/kubernetes/pkg/labels"
	watch "k8s.io/kubernetes/pkg/watch"
)

// FakeSensorAccesses implements SensorAccessInterface
type FakeSensorAccesses struct {
	Fake *FakeCore
	ns   string
}

func (c *FakeSensorAccesses) Create(sensorAccess *api.SensorAccess) (result *api.SensorAccess, err error) {
	obj, err := c.Fake.
		Invokes(core.NewCreateAction("sensoraccesses", c.ns, sensorAccess), &api.SensorAccess{})

	if obj == nil {
		return nil, err
	}
	return obj.(*api.SensorAccess), err
}

func (c *FakeSensorAccesses) Update(sensorAccess *api.SensorAccess) (result *api.SensorAccess, err error) {
	obj, err := c.Fake.
		Invokes(core.NewUpdateAction("sensoraccesses", c.ns, sensorAccess), &api.SensorAccess{})

	if obj == nil {
		return nil, err
	}
	return obj.(*api.SensorAccess), err
}

func (c *FakeSensorAccesses) Delete(name string, options *api.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(core.NewDeleteAction("sensoraccesses", c.ns, name), &api.SensorAccess{})

	return err
}

func (c *FakeSensorAccesses) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	action := core.NewDeleteCollectionAction("sensoraccesses", c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &api.SensorAccessList{})
	return err
}

func (c *FakeSensorAccesses) Get(name string) (result *api.SensorAccess, err error) {
	obj, err := c.Fake.
		Invokes(core.NewGetAction("sensoraccesses", c.ns, name), &api.SensorAccess{})

	if obj == nil {
		return nil, err
	}
	return obj.(*api.SensorAccess), err
}

func (c *FakeSensorAccesses) List(opts api.ListOptions) (result *api.SensorAccessList, err error) {
	obj, err := c.Fake.
		Invokes(core.NewListAction("sensoraccesses", c.ns, opts), &api.SensorAccessList{})

	if obj == nil {
		return nil, err
	}

	label := opts.LabelSelector
	if label == nil {
		label = labels.Everything()
	}
	list := &api.SensorAccessList{}
	for _, item := range obj.(*api.SensorAccessList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sensorAccesses.
func (c *FakeSensorAccesses) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(core.NewWatchAction("sensoraccesses", c.ns, opts))

}
