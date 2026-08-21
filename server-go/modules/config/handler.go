package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JBailes/aimee/server-go/bus"
	configcontract "github.com/JBailes/aimee/server-go/config"
)

// NewHandler serves the complete config store contract. The Store owns all
// parsing, validation, versioning, and persistence; the handler only translates
// the language-neutral event payload.
func NewHandler(store *Store) bus.ModuleHandler {
	return func(invocation bus.ModuleInvocation, body []byte) ([]byte, bus.ModuleStatus) {
		if invocation.StageID != configcontract.StageConfig || store == nil {
			return nil, bus.ModuleStatusInvalidRequest
		}
		if invocation.Cancelled() {
			return nil, bus.ModuleStatusCancelled
		}
		var request configcontract.Request
		decoder := json.NewDecoder(strings.NewReader(string(body)))
		decoder.UseNumber()
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil || request.Operation == "" {
			return nil, bus.ModuleStatusInvalidRequest
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			return nil, bus.ModuleStatusInvalidRequest
		}
		response := handleRequest(store, request)
		encoded, err := json.Marshal(response)
		if err != nil || uint32(len(encoded)) > bus.ModuleMessageMaxBody {
			return nil, bus.ModuleStatusInternal
		}
		return encoded, bus.ModuleStatusOK
	}
}

// DefaultPath resolves configuration inside the config process. Callers do not
// pass filesystem paths over the bus and never open the module's file.
func DefaultPath() (string, error) {
	if configured := os.Getenv("AIMEE_CONFIG_PATH"); configured != "" {
		return configured, nil
	}
	home := os.Getenv("AIMEE_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		home = filepath.Join(userHome, ".config", "aimee")
	}
	return filepath.Join(home, "aimee.yaml"), nil
}

func NewDefaultHandler() (bus.ModuleHandler, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	store, err := NewStore(path)
	if err != nil {
		return nil, err
	}
	return NewHandler(store), nil
}

func handleRequest(store *Store, request configcontract.Request) configcontract.Response {
	var response configcontract.Response
	var err error
	switch request.Operation {
	case configcontract.OpValues:
		response.Values, err = store.Values()
	case configcontract.OpValue:
		response.Value, err = store.Value(request.Key)
	case configcontract.OpStringValue:
		var found bool
		response.Value, found, err = store.StringValue(request.Key)
		if err == nil && !found {
			err = errors.New("config key not found")
		}
	case configcontract.OpSetVersioned:
		err = store.SetVersioned(request.Key, request.Value, request.PreviousVersion)
	case configcontract.OpVersion:
		response.Version, err = store.Version(request.Key)
	case configcontract.OpTriggerRules:
		response.Rules, err = store.TriggerRules()
	default:
		return configcontract.Response{Code: "invalid_operation", Error: "unknown config operation"}
	}
	if err != nil {
		code := "invalid"
		switch {
		case strings.Contains(err.Error(), "unknown config key") ||
			strings.Contains(err.Error(), "not found"):
			code = "not_found"
		case strings.Contains(err.Error(), "version conflict"):
			code = "conflict"
		}
		return configcontract.Response{Code: code, Error: err.Error()}
	}
	response.OK = true
	return response
}
