package restapi

import (
	"fmt"
	"testing"

	"github.com/gauravsarma1992/go-rest-api/gorpi"
	"github.com/stretchr/testify/assert"
)

type (
	GrandParentModel struct{}
	ParentModel      struct{}
	ChildModel       struct{}
)

func (dm *GrandParentModel) String() (name string) {
	name = "grand_parent"
	return
}

func (dm *GrandParentModel) Ancestor() (ancestor ResourceModel) {
	return
}

func (dm *ParentModel) String() (name string) {
	name = "parent"
	return
}

func (dm *ParentModel) Ancestor() (ancestor ResourceModel) {
	ancestor = &GrandParentModel{}
	return
}

func (dm *ChildModel) String() (name string) {
	name = "dummy"
	return
}

func (dm *ChildModel) Ancestor() (ancestor ResourceModel) {
	ancestor = &ParentModel{}
	return
}

func TestMgrInit(t *testing.T) {

	srv, err := gorpi.Default()
	assert.Nil(t, err)

	rMgr, err := NewRestApiManager(srv, nil)
	assert.Nil(t, err)
	assert.NotNil(t, rMgr)

	return
}

func TestMgrAddRoutes(t *testing.T) {

	srv, _ := gorpi.Default()

	rMgr, _ := NewRestApiManager(srv, nil)

	rRoute := &ResourceRoute{
		ResourceModel: &ChildModel{},
	}
	rMgr.AddResource(rRoute)
	fmt.Println(rMgr.GenerateRoutes())

	return
}
