package main

//MockRepositoryTemplate template mock repository
const MockRepositoryTemplate = `package mock

import (
	"context"
	"{{$.GoModules}}/internal/modules/{{$.module}}/domain"

	"github.com/mrapry/go-lib/golibshared"
	"github.com/stretchr/testify/mock"
)

// {{clean (upper $.module)}}Repository is an autogenerated mock type for the {{clean (upper $.module)}}Repository type
type {{clean (upper $.module)}}Repository struct {
	mock.Mock
}

// FindAll provides a mock function with given fields: ctx, filter
func (_m *{{clean (upper $.module)}}Repository) FindAll(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result {
	ret := _m.Called(ctx, filter)

	var r0 <-chan golibshared.Result
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Filter) <-chan golibshared.Result); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan golibshared.Result)
		}
	}

	return r0
}

// Find provides a mock function with given fields: ctx, obj
func (_m *{{clean (upper $.module)}}Repository) Find(ctx context.Context, obj domain.{{clean (upper $.module)}}) <-chan golibshared.Result {
	ret := _m.Called(ctx, obj)

	var r0 <-chan golibshared.Result
	if rf, ok := ret.Get(0).(func(context.Context, domain.{{clean (upper $.module)}}) <-chan golibshared.Result); ok {
		r0 = rf(ctx, obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan golibshared.Result)
		}
	}

	return r0
}

// FindByID provides a mock function with given fields: ctx, ID
func (_m *{{clean (upper $.module)}}Repository) FindByID(ctx context.Context, ID string) <-chan golibshared.Result {
	ret := _m.Called(ctx, ID)

	var r0 <-chan golibshared.Result
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan golibshared.Result); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan golibshared.Result)
		}
	}

	return r0
}

// Count provides a mock function with given fields: ctx, filter
func (_m *{{clean (upper $.module)}}Repository) Count(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result {
	ret := _m.Called(ctx, filter)

	var r0 <-chan golibshared.Result
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Filter) <-chan golibshared.Result); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan golibshared.Result)
		}
	}

	return r0
}

// Insert provides a mock function with given fields: ctx, newData
func (_m *{{clean (upper $.module)}}Repository) Insert(ctx context.Context, newData *domain.{{clean (upper $.module)}}) <-chan golibshared.Result {
	ret := _m.Called(ctx, newData)

	var r0 <-chan golibshared.Result
	if rf, ok := ret.Get(0).(func(context.Context, *domain.{{clean (upper $.module)}}) <-chan golibshared.Result); ok {
		r0 = rf(ctx, newData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan golibshared.Result)
		}
	}

	return r0
}

// Save provides a mock function with given fields: ctx, data
func (_m *{{clean (upper $.module)}}Repository) Save(ctx context.Context, data *domain.{{clean (upper $.module)}}) <-chan golibshared.Result {
	ret := _m.Called(ctx, data)

	var r0 <-chan golibshared.Result
	if rf, ok := ret.Get(0).(func(context.Context, *domain.{{clean (upper $.module)}}) <-chan golibshared.Result); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan golibshared.Result)
		}
	}

	return r0
}
`

//InterfaceRepositoryTemplate template interface repository
const InterfaceRepositoryTemplate = `package interfaces

import (
	"context"
	"{{$.GoModules}}/internal/modules/{{$.module}}/domain"

	"github.com/mrapry/go-lib/golibshared"
)

// {{clean (upper $.module)}}Repository abstract interface
type {{clean (upper $.module)}}Repository interface {
	FindAll(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result
	Count(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result
	Find(ctx context.Context, obj domain.{{clean (upper $.module)}}) <-chan golibshared.Result
	FindByID(tx context.Context, ID string) <-chan golibshared.Result
	Save(ctx context.Context, data *domain.{{clean (upper $.module)}}) <-chan golibshared.Result
	Insert(ctx context.Context, newData *domain.{{clean (upper $.module)}}) <-chan golibshared.Result
}
`

//ImplementRepositoryTemplate template implement repository with db mongodb
const ImplementRepositoryTemplate = `package mongodb

import (
	"context"
	"{{$.GoModules}}/internal/modules/{{$.module}}/domain"
	"{{$.GoModules}}/internal/modules/{{$.module}}/repository/interfaces"

	db "github.com/Kamva/mgm/v3"
	"github.com/Kamva/mgm/v3/operator"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/tracer"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type {{$.module}}RepoMongDB struct {
	readDB, writeDB *mongo.Database
}

// New{{clean (upper $.module)}}Repo create new {{$.module}} repository
func New{{clean (upper $.module)}}Repo(readDB, writeDB *mongo.Database) interfaces.{{clean (upper $.module)}}Repository {
	return &{{$.module}}RepoMongDB{readDB, writeDB}
}

func (r *{{$.module}}RepoMongDB) FindAll(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.{{clean (upper $.module)}}{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set offset
		filter.CalculateOffset()

		// set sort
		filter.SetSort()

		// set order by
		orderBy := filter.SetOrderBy(domain.FieldMap)

		// set default query
		where := []bson.M{}

		// set search
		fields := []string{"name"}
		where = filter.SetSearch(where, fields)

		// set show all
		if !filter.ShowAll {
			where = append(where, bson.M{"isActive": true})
		}

		// set option
		limit := cast.ToInt64(filter.Limit)
		offset := cast.ToInt64(filter.Offset)
		findOptions := &options.FindOptions{
			Limit: &limit,
			Skip:  &offset,
			Sort:  orderBy,
		}

		// set query
		query := bson.M{operator.And: where}

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.Find,
			Filter:     query,
			Sort:       findOptions.Sort,
			Limit:      *findOptions.Limit,
			Skip:       *findOptions.Skip,
		}
		trace.SetTags(ctx)

		var {{$.module}} = []*domain.{{clean (upper $.module)}}{}
		if err := coll.SimpleFind(&{{$.module}}, query, findOptions); err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: {{$.module}}}
	}()

	return output
}

func (r *{{$.module}}RepoMongDB) Count(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.{{clean (upper $.module)}}{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set default query
		where := []bson.M{}

		// set search
		fields := []string{"name"}
		where = filter.SetSearch(where, fields)

		// set show all
		if !filter.ShowAll {
			where = append(where, bson.M{"isActive": true})
		}

		// set query
		query := bson.M{operator.And: where}

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.CountDocument,
			Filter:     query,
		}
		trace.SetTags(ctx)

		count, err := coll.CountDocuments(ctx, query)
		if err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: count}
	}()

	return output
}

func (r *{{$.module}}RepoMongDB) Find(ctx context.Context, obj domain.{{clean (upper $.module)}}) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.{{clean (upper $.module)}}{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set data to bson M
		query := golibshared.ToBSON(obj)

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.FindOne,
			Filter:     query,
		}
		trace.SetTags(ctx)

		if err := coll.First(query, model); err != nil {
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: model}
	}()

	return output
}

func (r *{{$.module}}RepoMongDB) FindByID(ctx context.Context, ID string) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.{{clean (upper $.module)}}{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.readDB, collName)

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.FindOne,
			Filter:     ID,
		}
		trace.SetTags(ctx)

		if err := coll.FindByID(ID, model); err != nil {
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: model}
	}()

	return output
}

func (r *{{$.module}}RepoMongDB) Save(ctx context.Context, data *domain.{{clean (upper $.module)}}) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.{{clean (upper $.module)}}{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.writeDB, collName)

		// set version
		data.Version = data.Version + 1

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.UpdateOne,
			Filter:     data,
		}
		trace.SetTags(ctx)

		if err := coll.Update(data); err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: data}
	}()

	return output
}

func (r *{{$.module}}RepoMongDB) Insert(ctx context.Context, newData *domain.{{clean (upper $.module)}}) <-chan golibshared.Result {
	output := make(chan golibshared.Result)

	go func() {
		defer close(output)

		// set model
		model := &domain.{{clean (upper $.module)}}{}

		// set collection name
		collName := db.CollName(model)

		// set collection
		coll := db.NewCollection(r.writeDB, collName)

		// set version
		newData.Version = newData.Version + 1

		// set tracer mongo
		trace := &tracer.TraceMongo{
			Collection: collName,
			Method:     tracer.InsertOne,
			Filter:     newData,
		}
		trace.SetTags(ctx)

		if err := coll.Create(newData); err != nil {
			tracer.SetError(ctx, err)
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: newData}

	}()

	return output
}
`

//RepositoryTemplate template repository abstraction
const RepositoryTemplate = `package repository

import (
	"{{$.GoModules}}/internal/modules/{{$.module}}/repository/interfaces"
	"{{$.GoModules}}/internal/modules/{{$.module}}/repository/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repository parent
type Repository struct {
	readDB, writeDB *mongo.Database
	{{clean (upper $.module)}}            interfaces.{{clean (upper $.module)}}Repository
}

// NewRepository create new repository
func NewRepository(read, write *mongo.Database) *Repository {
	return &Repository{
		readDB: read, writeDB: write,
		{{clean (upper $.module)}}: mongodb.New{{clean (upper $.module)}}Repo(read, write),
	}
}
`