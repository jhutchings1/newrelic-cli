// +build unit

package edge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestList(t *testing.T) {
	assert.Equal(t, "list", cmdList.Name())

	testcobra.CheckCobraMetadata(t, cmdList)
	testcobra.CheckCobraRequiredFlags(t, cmdList, []string{})
}

func TestCreate(t *testing.T) {
	assert.Equal(t, "create", cmdCreate.Name())

	testcobra.CheckCobraMetadata(t, cmdCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdCreate, []string{})
}

func TestDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdDelete.Name())

	testcobra.CheckCobraMetadata(t, cmdDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdDelete, []string{})
}
