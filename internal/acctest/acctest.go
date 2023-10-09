package acctest

import (
	"os"
	"testing"
)

// Verify that any required environment variables are set before the test begins
func ConfigurationPreCheck(t *testing.T) {
	envVars := []string{
		//"EXAMPLE",
	}

	errorFound := false
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Errorf("The '%s' environment variable must be set to run acceptance tests", envVar)
			errorFound = true
		}
	}

	if errorFound {
		t.FailNow()
	}
}
