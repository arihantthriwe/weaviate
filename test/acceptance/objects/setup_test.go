//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate/client/objects"

	"github.com/weaviate/weaviate/client/schema"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/test/helper"
)

// Tests that sort parameters are validated with the correct class
func TestSort(t *testing.T) {
	createObjectClass(t, &models.Class{
		Class: "ClassToSort",
		Properties: []*models.Property{
			{
				Name:     "name",
				DataType: []string{"string"},
			},
		},
	})
	defer deleteObjectClass(t, "ClassToSort")

	createObjectClass(t, &models.Class{
		Class: "OtherClass",
		Properties: []*models.Property{
			{
				Name:     "ref",
				DataType: []string{"ClassToSort"},
			},
		},
	})
	defer deleteObjectClass(t, "OtherClass")

	listParams := objects.NewObjectsListParams()
	nameClass := "ClassToSort"
	nameProp := "name"
	limit := int64(5)
	listParams.Class = &nameClass
	listParams.Sort = &nameProp
	listParams.Limit = &limit

	_, err := helper.Client(t).Objects.ObjectsList(listParams, nil)
	require.Nil(t, err, "should not error")
}

func Test_Objects(t *testing.T) {
	createObjectClass(t, &models.Class{
		Class: "TestObject",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"vectorizeClassName": true,
			},
		},
		Properties: []*models.Property{
			{
				Name:     "testString",
				DataType: []string{"string"},
			},
			{
				Name:     "testWholeNumber",
				DataType: []string{"int"},
			},
			{
				Name:     "testNumber",
				DataType: []string{"number"},
			},
			{
				Name:     "testDateTime",
				DataType: []string{"date"},
			},
			{
				Name:     "testTrueFalse",
				DataType: []string{"boolean"},
			},
			{
				Name:     "testPhoneNumber",
				DataType: []string{"phoneNumber"},
			},
		},
	})
	createObjectClass(t, &models.Class{
		Class:      "TestObjectCustomVector",
		Vectorizer: "none",
		Properties: []*models.Property{
			{
				Name:     "description",
				DataType: []string{"text"},
			},
		},
	})
	createObjectClass(t, &models.Class{
		Class:      "TestDeleteClassOne",
		Vectorizer: "none",
		Properties: []*models.Property{
			{
				Name:     "text",
				DataType: []string{"text"},
			},
		},
	})
	createObjectClass(t, &models.Class{
		Class:      "TestDeleteClassTwo",
		Vectorizer: "none",
		Properties: []*models.Property{
			{
				Name:     "text",
				DataType: []string{"text"},
			},
		},
	})

	// tests
	t.Run("listing objects", listingObjects)
	t.Run("searching for neighbors", searchNeighbors)
	t.Run("running a feature projection", featureProjection)
	t.Run("creating objects", creatingObjects)

	t.Run("custom vector journey", customVectors)
	t.Run("auto schema", autoSchemaObjects)
	t.Run("checking object's existence", checkObjects)
	t.Run("delete request deletes all objects with a given ID", deleteAllObjectsFromAllClasses)

	// tear down
	deleteObjectClass(t, "TestObject")
	deleteObjectClass(t, "TestObjectCustomVector")
	deleteObjectClass(t, "NonExistingClass")
	deleteObjectClass(t, "TestDeleteClassOne")
	deleteObjectClass(t, "TestDeleteClassTwo")
}

func createObjectClass(t *testing.T, class *models.Class) {
	params := schema.NewSchemaObjectsCreateParams().WithObjectClass(class)
	resp, err := helper.Client(t).Schema.SchemaObjectsCreate(params, nil)
	helper.AssertRequestOk(t, resp, err, nil)
}

func deleteObjectClass(t *testing.T, class string) {
	delParams := schema.NewSchemaObjectsDeleteParams().WithClassName(class)
	delRes, err := helper.Client(t).Schema.SchemaObjectsDelete(delParams, nil)
	helper.AssertRequestOk(t, delRes, err, nil)
}
