package config

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var validDotenvConfig = []byte(`APP_ADDR=:8080
APP_CONFIG_NAME=.env
APP_CONFIG_PATH=.
DATABASE_TYPE=memory`)

var invalidDotenvConfig = []byte(`edededde`)

func TestNew(t *testing.T) {
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, ".env", validDotenvConfig, 0755)

	t.Run("Load from valid dotenv file", func(t *testing.T) {
		c := New()
		assert.NotEmpty(t, c)
		s := bufio.NewScanner(strings.NewReader(string(validDotenvConfig)))
		for s.Scan() {
			assert.Equal(t, strings.Split(s.Text(), "=")[1], c.GetString(strings.Split(s.Text(), "=")[0]))
		}
	})
}

// https://talks.golang.org/2014/testing.slide
func TestNewNonExistentDotenv(t *testing.T) {
	if os.Getenv("BE_CRASHER") != "1" {

		t.Run("Load from non existent dotenv file", func(t *testing.T) {
			os.Setenv("APP_CONFIG_NAME", "FileThatDoesntExist")
			_ = New()
		})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNewNonExistentDotenv")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.ErrorAs(t, err, err.(*exec.ExitError))
}

func TestNewInvalidDotenv(t *testing.T) {
	if os.Getenv("BE_CRASHER") != "1" {
		appFS := afero.NewMemMapFs()
		afero.WriteFile(appFS, ".env", invalidDotenvConfig, 0755)

		t.Run("Load from invalid dotenv file", func(t *testing.T) {
			_ = New()
		})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNewInvalidDotenv")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.ErrorAs(t, err, err.(*exec.ExitError))
}
