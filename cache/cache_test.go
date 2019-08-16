package cache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jivalabs/picsum-photos/cache"
	"github.com/jivalabs/picsum-photos/cache/mock"
)

var mockLoaderFunc cache.LoaderFunc = func(ctx context.Context, key string) (data []byte, err error) {
	if key == "notfounderr" {
		return nil, fmt.Errorf("notfounderr")
	}

	return []byte("notfound"), nil
}

func TestAuto(t *testing.T) {
	auto := &cache.Auto{
		Provider: &mock.Provider{},
		Loader:   mockLoaderFunc,
	}

	tests := []struct {
		Key           string
		ExpectedError error
	}{
		{"foo", nil},
		{"notfound", nil},
		{"notfounderr", fmt.Errorf("notfounderr")},
		{"seterror", fmt.Errorf("seterror")},
	}

	for _, test := range tests {
		data, err := auto.Get(context.Background(), test.Key)
		if err != nil {
			if test.ExpectedError == nil {
				t.Errorf("%s: %s", test.Key, err)
				continue
			}

			if test.ExpectedError.Error() != err.Error() {
				t.Errorf("%s: wrong error: %s", test.Key, err)
				continue
			}

			continue
		}

		if string(data) != test.Key {
			t.Errorf("%s: wrong data", test.Key)
		}
	}

}
