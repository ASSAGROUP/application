package resource

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"reflect"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/application/pkg/utils"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// Resourcer interface
type Resourcer interface {
	GetResource() *Resource
	GetMetas([]string) []Metaor
	CallFindMany(interface{}, *appsvr.Context) error
	CallFindOne(interface{}, *MetaValues, *appsvr.Context) error
	CallSave(interface{}, *appsvr.Context) error
	CallDelete(interface{}, *appsvr.Context) error
	NewSlice() interface{}
	NewStruct() interface{}
}

// ConfigureResourceBeforeInitializeInterface if a struct implemented this interface, it will be called before everything when create a resource with the struct
type ConfigureResourceBeforeInitializeInterface interface {
	ConfigureQorResourceBeforeInitialize(Resourcer)
}

// ConfigureResourceInterface if a struct implemented this interface, it will be called after configured by user
type ConfigureResourceInterface interface {
	ConfigureQorResource(Resourcer)
}

// Resource is a struct that including basic definition of Bhojpur Application resource
type Resource struct {
	Name            string
	Value           interface{}
	PrimaryFields   []*orm.StructField
	FindManyHandler func(interface{}, *appsvr.Context) error
	FindOneHandler  func(interface{}, *MetaValues, *appsvr.Context) error
	SaveHandler     func(interface{}, *appsvr.Context) error
	DeleteHandler   func(interface{}, *appsvr.Context) error
	Permission      *roles.Permission
	Validators      []*Validator
	Processors      []*Processor
	primaryField    *orm.Field
}

// New initialize Bhojpur Application resource
func New(value interface{}) *Resource {
	var (
		name = utils.HumanizeString(utils.ModelType(value).Name())
		res  = &Resource{Value: value, Name: name}
	)

	res.FindOneHandler = res.findOneHandler
	res.FindManyHandler = res.findManyHandler
	res.SaveHandler = res.saveHandler
	res.DeleteHandler = res.deleteHandler
	res.SetPrimaryFields()
	return res
}

// GetResource return itself to match interface `Resourcer`
func (res *Resource) GetResource() *Resource {
	return res
}

// SetPrimaryFields set primary fields
func (res *Resource) SetPrimaryFields(fields ...string) error {
	scope := orm.Scope{Value: res.Value}
	res.PrimaryFields = nil

	if len(fields) > 0 {
		for _, fieldName := range fields {
			if field, ok := scope.FieldByName(fieldName); ok {
				res.PrimaryFields = append(res.PrimaryFields, field.StructField)
			} else {
				return fmt.Errorf("%v is not a valid field for resource %v", fieldName, res.Name)
			}
		}
		return nil
	}

	if primaryField := scope.PrimaryField(); primaryField != nil {
		res.PrimaryFields = []*orm.StructField{primaryField.StructField}
		return nil
	}

	return fmt.Errorf("no valid primary field for resource %v", res.Name)
}

// Validator validator struct
type Validator struct {
	Name    string
	Handler func(interface{}, *MetaValues, *appsvr.Context) error
}

// AddValidator add validator to resource, it will invoked when creating, updating, and will rollback the change if validator return any error
func (res *Resource) AddValidator(validator *Validator) {
	for idx, v := range res.Validators {
		if v.Name == validator.Name {
			res.Validators[idx] = validator
			return
		}
	}

	res.Validators = append(res.Validators, validator)
}

// Processor processor struct
type Processor struct {
	Name    string
	Handler func(interface{}, *MetaValues, *appsvr.Context) error
}

// AddProcessor add processor to resource, it is used to process data before creating, updating, will rollback the change if it return any error
func (res *Resource) AddProcessor(processor *Processor) {
	for idx, p := range res.Processors {
		if p.Name == processor.Name {
			res.Processors[idx] = processor
			return
		}
	}
	res.Processors = append(res.Processors, processor)
}

// NewStruct initialize a struct for the Resource
func (res *Resource) NewStruct() interface{} {
	if res.Value == nil {
		return nil
	}
	return reflect.New(utils.Indirect(reflect.ValueOf(res.Value)).Type()).Interface()
}

// NewSlice initialize a slice of struct for the Resource
func (res *Resource) NewSlice() interface{} {
	if res.Value == nil {
		return nil
	}
	sliceType := reflect.SliceOf(reflect.TypeOf(res.Value))
	slice := reflect.MakeSlice(sliceType, 0, 0)
	slicePtr := reflect.New(sliceType)
	slicePtr.Elem().Set(slice)
	return slicePtr.Interface()
}

// GetMetas get defined metas, to match interface `Resourcer`
func (res *Resource) GetMetas([]string) []Metaor {
	panic("not defined")
}

// HasPermission check permission of resource
func (res *Resource) HasPermission(mode roles.PermissionMode, context *appsvr.Context) bool {
	if res == nil || res.Permission == nil {
		return true
	}

	var roles = []interface{}{}
	for _, role := range context.Roles {
		roles = append(roles, role)
	}
	return res.Permission.HasPermission(mode, roles...)
}
