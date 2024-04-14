/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package timex

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestTimex_New(t *testing.T) {
	t.Parallel()

	now := time.Now()
	actual := Timex{now}
	assert.Equal(t, now, actual.Time)

	actual = Timex{time.Time{}}
	assert.Equal(t, time.Time{}, actual.Time)
}

func TestTimex_SerializeNow(t *testing.T) {
	t.Parallel()

	jt1 := Timex{time.Now()}
	b, err := json.Marshal(jt1)
	assert.NoError(t, err)
	var jt2 Timex
	err = json.Unmarshal(b, &jt2)
	assert.NoError(t, err)
	assert.True(t, jt1.Equal(jt2))
}

func TestTimex_SerializeZero(t *testing.T) {
	t.Parallel()

	jt1 := Timex{time.Time{}}
	b, err := json.Marshal(jt1)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), b)
	var jt2 Timex
	err = json.Unmarshal(b, &jt2)
	assert.NoError(t, err)
	assert.True(t, jt1.Equal(jt2))
}

func TestTimex_Format(t *testing.T) {
	t.Parallel()

	jt1 := Timex{time.Time{}}
	s1 := jt1.Format(time.RFC3339)
	assert.Equal(t, "", s1)

	jt2, err := Parse(time.RFC3339, "")
	assert.NoError(t, err)
	assert.True(t, jt1.Equal(jt2))

	jt3 := MustParse("", "2015-01-14T11:12:13Z")
	s3 := jt3.Format(time.RFC3339)
	assert.Equal(t, "2015-01-14T11:12:13Z", s3)
}

func TestTimex_JSONEmpty(t *testing.T) {
	t.Parallel()

	var jt1 Timex
	err := json.Unmarshal([]byte(`""`), &jt1)
	assert.NoError(t, err)
	assert.True(t, jt1.IsZero())

	err = json.Unmarshal([]byte(`null`), &jt1)
	assert.NoError(t, err)
	assert.True(t, jt1.IsZero())

	jt1 = Timex{time.Now()}
	err = json.Unmarshal([]byte(`""`), &jt1)
	assert.NoError(t, err)
	assert.True(t, jt1.IsZero())

	jt1 = Timex{}
	b, err := json.Marshal(jt1)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestTimex_YAMLEmpty(t *testing.T) {
	t.Parallel()

	var jt1 Timex
	err := yaml.Unmarshal([]byte(`""`), &jt1)
	assert.NoError(t, err)
	assert.True(t, jt1.IsZero())

	err = yaml.Unmarshal([]byte(`null`), &jt1)
	assert.NoError(t, err)
	assert.True(t, jt1.IsZero())

	jt1 = Timex{time.Now()}
	err = yaml.Unmarshal([]byte(`""`), &jt1)
	assert.NoError(t, err)
	assert.True(t, jt1.IsZero())

	jt1 = Timex{}
	b, err := yaml.Marshal(jt1)
	assert.NoError(t, err)
	assert.Equal(t, "null\n", string(b))
}

func TestTimex_UnmarshalJSONInvalid(t *testing.T) {
	t.Parallel()

	var jt2 Timex
	err := json.Unmarshal(nil, &jt2)
	assert.Error(t, err)

	err = json.Unmarshal([]byte("not-a-time"), &jt2)
	assert.Error(t, err)
}

func TestTimex_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var jt1 Timex
	err := json.Unmarshal([]byte(`"2021-08-11T10:00:00Z"`), &jt1)
	assert.NoError(t, err)
	assert.False(t, jt1.IsZero())

	jt2 := Timex{time.Now()}
	b, err := json.Marshal(jt2)
	assert.NoError(t, err)
	assert.NotNil(t, b)

	var jt3 Timex
	err = json.Unmarshal(b, &jt3)
	assert.NoError(t, err)
	assert.True(t, jt2.Equal(jt3))
}

func TestTimex_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	var jt1 Timex
	err := yaml.Unmarshal([]byte(`"2021-08-11T10:00:00Z"`), &jt1)
	assert.NoError(t, err)
	assert.False(t, jt1.IsZero())

	jt2 := Timex{time.Now()}
	b, err := yaml.Marshal(jt2)
	assert.NoError(t, err)
	assert.NotNil(t, b)

	var jt3 Timex
	err = yaml.Unmarshal(b, &jt3)
	assert.NoError(t, err)
	assert.True(t, jt2.Equal(jt3))
}

func TestTimex_ParseNullTime(t *testing.T) {
	t.Parallel()

	assert.Equal(t, time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC), MustParse("", "2015-01-02").Time)
	assert.Equal(t, time.Date(2015, 1, 2, 11, 12, 13, 0, time.UTC), MustParse("", "2015-01-02T11:12:13").Time)
	assert.Equal(t, time.Date(2015, 1, 2, 11, 12, 13, 0, time.UTC), MustParse("", "2015-01-02 11:12:13").Time)
	assert.Equal(t, time.Date(2015, 1, 2, 11, 12, 13, 0, time.UTC), MustParse("", "2015-01-02T11:12:13Z").Time)
	assert.Equal(t, time.Date(2015, 1, 2, 11, 12, 13, 500000, time.UTC), MustParse("", "2015-01-02T11:12:13.000500000Z").Time)
}
