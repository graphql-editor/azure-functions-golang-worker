// Package pluginloader implements loader using go plugin functionality.
package pluginloader

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"sync"

	"github.com/graphql-editor/azure-functions-golang-worker/api"
	"github.com/graphql-editor/azure-functions-golang-worker/worker"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func newGoBuild() (gobuild gobuild, err error) {
	gobuild.Context = build.Default
	if goroot := os.Getenv("GOROOT"); goroot != "" {
		gobuild.GOROOT = goroot
	}
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		gobuild.GOPATH = gopath
	}
	if goarch := os.Getenv("GOARCH"); goarch != "" {
		gobuild.GOARCH = goarch
	}
	if goos := os.Getenv("GOOS"); goos != "" {
		gobuild.GOOS = goos
	}
	if useCgo := os.Getenv("CGO_ENABLED"); useCgo != "" {
		gobuild.CgoEnabled = useCgo != "0"
	}
	if gocache := os.Getenv("GOCACHE"); gocache != "" {
		gobuild.GOCACHE = gocache
	} else {
		var cacheDir string
		cacheDir, err = os.UserCacheDir()
		if err != nil {
			return
		}
		gobuild.GOCACHE = filepath.Join(cacheDir, "function-go-build")
	}
	return
}

type gobuild struct {
	build.Context
	GOCACHE string
}

func (g gobuild) tmpdir() (string, error) {
	dir, err := ioutil.TempDir("", "azure-golang-function")
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (g gobuild) build(logger api.Logger, path, out string) error {
	goexe := filepath.Join(g.GOROOT, "bin", "go"+ext)
	args := []string{"build", "-buildmode=plugin", "-o", out}
	if raceBuild {
		args = append(args, "-race")
	}
	args = append(args, path)
	cmd := exec.Command(goexe, args...)
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"GOPATH=" + g.GOPATH,
		"GOROOT=" + g.GOROOT,
		"GOARCH=" + g.GOARCH,
		"GOOS=" + g.GOOS,
		"GOCACHE=" + g.GOCACHE,
	}
	_, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Error(string(exitErr.Stderr))
		}
	}
	return err
}

// Loader is a plugin loader that builds go function and returns reflection of function type
type Loader struct {
	binaries []string
	lock     sync.Mutex
}

// GetFunctionType returns reflection of function type from go plugin
func (l *Loader) GetFunctionType(fi worker.FunctionInfo, logger api.Logger) (reflect.Type, error) {
	gobuild, err := newGoBuild()
	if err != nil {
		return nil, err
	}
	fpath, err := gobuild.tmpdir()
	if err != nil {
		return nil, err
	}
	l.lock.Lock()
	l.binaries = append(l.binaries, fpath)
	l.lock.Unlock()
	fpath = filepath.Join(fpath, "function")
	path := "./"
	if fi.ScriptFile != "" {
		path = fi.ScriptFile
	}
	if err := gobuild.build(logger, path, fpath); err != nil {
		return nil, errors.Wrap(err, "function build failed")
	}
	plug, err := plugin.Open(fpath)
	if err != nil {
		return nil, errors.Wrap(err, "failed loading function plugin")
	}
	entrypoint, err := plug.Lookup(fi.EntryPoint)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed loooking up function entrypoint: %s", fi.EntryPoint))
	}
	return reflect.TypeOf(entrypoint).Elem(), nil
}

// Close cleans up after loader. Must be called before program exit to cleanup temporary binaries created by loader.
func (l *Loader) Close() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	var collectErrors error
	for _, bin := range l.binaries {
		if err := os.RemoveAll(bin); err != nil {
			collectErrors = multierror.Append(collectErrors, err)
		}
	}
	return collectErrors
}

// NewLoader creates new loader
func NewLoader() *Loader {
	return &Loader{}
}
