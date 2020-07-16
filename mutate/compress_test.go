package mutate

import (
	"bytes"
	"io"
	"io/ioutil"
	"runtime"

	gzip "github.com/klauspost/pgzip"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	fact = "meshuggah rocks!!!"
)

func TestNoopCompressor(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBufferString(fact)
	c := NewNoopCompressor(fact)

	r, err := c.Compress(buf)
	assert.NoError(err)
	assert.Equal(fact, c.MediaType())

	content, err := ioutil.ReadAll(r)
	assert.NoError(err)

	assert.Equal(string(content), fact)
}

func TestGzipCompressor(t *testing.T) {
	assert := assert.New(t)

	buf := bytes.NewBufferString(fact)
	c := NewGzipCompressor(buf)

	r, err := c.Compress(buf)
	assert.NoError(err)

	r, err := gzip.NewReader(r)
	assert.NoError(err)

	content, err := ioutil.ReadAll(r)
	assert.NoError(err)

	assert.Equal(string(content), fact)
}
