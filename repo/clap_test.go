package repo

import (
	"context"
	"errors"
	"testing"

	mockDB "github.com/kh411d/clapi/db/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClap(t *testing.T) {
	mockDBConn := &mockDB.KV{}
	mockDBConn.On("WithContext", mock.Anything).Return(mockDBConn)
	mockDBConn.On("Get", mock.Anything).Return([]byte("9"), nil)
	mockDBConn.On("IncrBy", mock.Anything, mock.Anything).Return(nil)

	c := &claps{}
	r := c.GetClap(context.TODO(), mockDBConn, "anything-url")
	assert.Equal(t, r, "9", "wrong value")

	c.AddClap(context.TODO(), mockDBConn, "anything-url", 9)
}

func TestClapNeg(t *testing.T) {
	mockDBConn := &mockDB.KV{}
	mockDBConn.On("WithContext", mock.Anything).Return(mockDBConn)
	mockDBConn.On("Get", mock.Anything).Return(nil, nil)
	mockDBConn.On("IncrBy", mock.Anything, mock.Anything).Return(errors.New("error printed"))

	c := &claps{}
	r := c.GetClap(context.TODO(), mockDBConn, "anything-url")
	assert.Equal(t, r, "0", "wrong value")

	c.AddClap(context.TODO(), mockDBConn, "anything-url", 9)
}
