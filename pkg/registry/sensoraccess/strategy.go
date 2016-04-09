/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package sensoraccess

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

// sensoraccessStrategy implements behavior for SensorAccess objects
type sensoraccessStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

// Strategy is the default logic that applies when creating and updating SensorAccess
// objects via the REST API.
var Strategy = sensoraccessStrategy{api.Scheme, api.SimpleNameGenerator}

// NamespaceScoped is true for sensoraccesss.
func (sensoraccessStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (sensoraccessStrategy) PrepareForCreate(obj runtime.Object) {
	sensoraccess := obj.(*api.SensorAccess)
	sensoraccess.Access = 0
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (sensoraccessStrategy) PrepareForUpdate(obj, old runtime.Object) {
	newSensorAccess := obj.(*api.SensorAccess)
	oldSensorAccess := old.(*api.SensorAccess)
	newSensorAccess.Access = oldSensorAccess.Access
}

// Validate validates a new sensoraccess.
func (sensoraccessStrategy) Validate(ctx api.Context, obj runtime.Object) field.ErrorList {
	sensoraccess := obj.(*api.SensorAccess)
	return validation.ValidateSensorAccess(sensoraccess)
}

// Canonicalize normalizes the object after validation.
func (sensoraccessStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for sensoraccesss.
func (sensoraccessStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (sensoraccessStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateSensorAccess(obj.(*api.SensorAccess))
	return append(errorList, validation.ValidateSensorAccessUpdate(obj.(*api.SensorAccess), old.(*api.SensorAccess))...)
}

func (sensoraccessStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchSensorAccess returns a generic matcher for a given label and field selector.
func MatchSensorAccess(label labels.Selector, field fields.Selector) generic.Matcher {
	return generic.MatcherFunc(func(obj runtime.Object) (bool, error) {
		resourcequotaObj, ok := obj.(*api.SensorAccess)
		if !ok {
			return false, fmt.Errorf("not a resourcequota")
		}
		fields := SensorAccessToSelectableFields(resourcequotaObj)
		return label.Matches(labels.Set(resourcequotaObj.Labels)) && field.Matches(fields), nil
	})
}

// SensorAccessToSelectableFields returns a label set that represents the object
func SensorAccessToSelectableFields(sensoraccess *api.SensorAccess) labels.Set {
	return labels.Set(generic.ObjectMetaFieldsSet(sensoraccess.ObjectMeta, true))
}
