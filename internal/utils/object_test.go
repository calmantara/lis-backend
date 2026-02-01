package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectMapperSuccess(t *testing.T) {
	type Source struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	type Destination struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	source := Source{
		Field1: "test",
		Field2: 123,
	}

	var destination Destination
	err := ObjectMapper(source, &destination)

	assert.NoError(t, err, "ObjectMapper should not return an error")
	assert.Equal(t, source.Field1, destination.Field1, "Field1 should match")
	assert.Equal(t, source.Field2, destination.Field2, "Field2 should match")
}

func TestObjectMapperInvalidInput(t *testing.T) {
	type Source struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	type Destination struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"` // Different type to cause an error
	}

	source := Source{
		Field1: "test",
		Field2: 123,
	}

	var destination Destination
	err := ObjectMapper(source, &destination)

	assert.Error(t, err, "ObjectMapper should return an error for incompatible types")
}

func TestObjectMapperInvalidChannelInput(t *testing.T) {
	input := make(chan int)

	type Destination struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"` // Different type to cause an error
	}

	var destination Destination
	err := ObjectMapper(input, &destination)

	assert.Error(t, err, "ObjectMapper should return an error for incompatible types")
}

func TestObjectMapperEmptyInput(t *testing.T) {
	type Source struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	type Destination struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	source := Source{}
	var destination Destination
	err := ObjectMapper(source, &destination)

	assert.NoError(t, err, "ObjectMapper should not return an error for empty input")
	assert.Equal(t, source.Field1, destination.Field1, "Field1 should match")
	assert.Equal(t, source.Field2, destination.Field2, "Field2 should match")
}
