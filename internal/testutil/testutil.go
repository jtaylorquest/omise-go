package testutil

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/omise/omise-go"
	a "github.com/stretchr/testify/assert"
)

func Keys() (string, string) {
	return os.Getenv("OMISE_PUB_KEY"),
		os.Getenv("OMISE_SECRET_KEY")
}

func Main(m *testing.M) {
	// never test against live key.
	pkey, skey := Keys()
	switch {
	case pkey == "" || skey == "":
		log.Fatalln("no test key specified, please set both $OMISE_PUB_KEY and " +
			"$OMISE_SECRET_KEY")
	case !strings.HasPrefix(pkey, "pkey_test_") || !strings.HasPrefix(skey, "skey_test_"):
		log.Fatalln("specified key is invalid or is not a test key!!! You might lose money!!!")
	}

	flag.Parse()
	os.Exit(m.Run())
}

func NewClient() (*omise.Client, error) {
	return omise.NewClient(Keys())
}

func AssertJSONEquals(t *testing.T, m1 map[string]interface{}, m2 map[string]interface{}) bool {
	for k, v := range m1 {
		v2, ok := m2[k]
		if !a.True(t, ok, "missing `"+k+"` key") {
			return false
		}
		if !a.IsType(t, v, v2, "mismatched type for `"+k+"` key") {
			return false
		} // postcond: v.(type) == v2.(type)

		if vmap, ok := v.(map[string]interface{}); ok {
			vmap2 := v2.(map[string]interface{}) // IsType
			if !AssertJSONEquals(t, vmap, vmap2) {
				return false
			}

		} else { // not map
			// TODO: Check arrays as well?
			if !a.Equal(t, v, v2, "mismatched value for `"+k+"` key") {
				return false
			}
		}
	}

	return true
}
