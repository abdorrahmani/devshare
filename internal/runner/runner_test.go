package runner

import (
	"reflect"
	"testing"
)

func TestPmRunArgs(t *testing.T) {
	// npm needs "run <script>" and a "--" separator before script flags
	if got := pmRunArgs("npm", "dev", []string{"--port", "5173"}); !reflect.DeepEqual(got, []string{"run", "dev", "--", "--port", "5173"}) {
		t.Errorf("npm with flags: got %v", got)
	}
	// npm with no flags: no trailing separator
	if got := pmRunArgs("npm", "dev", nil); !reflect.DeepEqual(got, []string{"run", "dev"}) {
		t.Errorf("npm no flags: got %v", got)
	}
	// yarn/pnpm forward extra args to the script directly
	if got := pmRunArgs("yarn", "dev", []string{"--port", "5173"}); !reflect.DeepEqual(got, []string{"dev", "--port", "5173"}) {
		t.Errorf("yarn: got %v", got)
	}
}
