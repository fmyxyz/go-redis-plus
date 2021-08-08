package redis

import (
	"context"
	"testing"
	"time"
)

type kstruct struct {
	K         string        `json:"k"`
	K1        uint8         `json:"k_1"`
	K1_1      float64       `json:"k_1_1"`
	Ktrue     bool          `json:"ktrue"`
	Ktime     time.Time     `json:"ktime"`
	Kduration time.Duration `json:"kduration"`
	Kstruct   subkstruct    `json:"kstruct"`
}

type subkstruct struct {
	K         string        `json:"k"`
	K1        int           `json:"k_1"`
	K1_1      float32       `json:"k_1_1"`
	Ktrue     bool          `json:"ktrue"`
	Ktime     time.Time     `json:"ktime"`
	Kduration time.Duration `json:"kduration"`
}

var client *Client

func init() {
	opt := &Options{}
	opt.Addr = "localhost:6379"
	opt.Expiration = time.Second * 7
	client = NewClient(opt)
}

func TestClient_SetMapValue(t *testing.T) {
	err := client.SetMapValue(context.Background(), "kmap[interface{}]string", map[interface{}]string{
		"v":        "k",
		1:          "k1",
		1.1:        "k1.1",
		true:       "ktrue",
		time.Now(): "ktime",
		subkstruct{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     false,
			Ktime:     time.Now(),
			Kduration: time.Second,
		}: "kstruct",
	})
	if err != nil {
		t.Error(err)
	}
	err = client.SetMapValue(context.Background(), "kmap[string]interface{}", map[string]interface{}{
		"k":     "v",
		"k1":    1,
		"k1.1":  1.1,
		"ktrue": true,
		"ktime": time.Now(),
		"k[]struct": []subkstruct{{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     true,
			Ktime:     time.Now(),
			Kduration: time.Second,
		}},
		"k[]string": []string{"test1", "test2", "test3"},
		"kstruct": subkstruct{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     false,
			Ktime:     time.Now(),
			Kduration: time.Second,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClient_GetMapValue(t *testing.T) {
	m := map[string]interface{}{}
	err := client.GetMapValue(context.Background(), "kmap[interface{}]string", m)
	if err != nil {
		t.Error(err)
	}
	t.Log("")
	for k, v := range m {
		t.Logf("%+v:%+v", k, v)
	}

	mm := map[string]interface{}{}
	err = client.GetMapValue(context.Background(), "kmap[string]interface{}", mm)
	if err != nil {
		t.Error(err)
	}
	t.Log("")
	for k, v := range mm {
		t.Logf("%+v:%+v", k, v)
	}

	kmapstringstring := make(map[string]string)
	err = client.GetMapValue(context.Background(), "kmap[string]interface{}", &kmapstringstring)
	if err != nil {
		t.Error(err)
	}
	t.Log("")
	for k, v := range kmapstringstring {
		t.Logf("%+v:%+v", k, v)
	}
}

func TestClient_SetSingleValue(t *testing.T) {
	err := client.SetSingleValue(context.Background(), "k", "v")
	if err != nil {
		t.Error(err)
	}
	err = client.SetSingleValue(context.Background(), "k1", 1)
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "k1.1", 1.1)
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "ktrue", true)
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "ktime", time.Now())
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "kduration", time.Second)
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "k[]string", []string{"test1", "test2", "test3"})
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "k[]int", []int{1, 2, 3, 4, 5})
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "k[]struct", []subkstruct{{
		K:         "v",
		K1:        1,
		K1_1:      1.1,
		Ktrue:     true,
		Ktime:     time.Now(),
		Kduration: time.Second,
	}})
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "kstruct", kstruct{
		K:         "v",
		K1:        1,
		K1_1:      1.1,
		Ktrue:     true,
		Ktime:     time.Now(),
		Kduration: time.Second,
		Kstruct: subkstruct{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     true,
			Ktime:     time.Now(),
			Kduration: time.Second,
		},
	})
	if err != nil {
		t.Error(err)
	}

	err = client.SetSingleValue(context.Background(), "kmap[string]interface{}", map[string]interface{}{
		"k":     "v",
		"k1":    1,
		"k1.1":  1.1,
		"ktrue": true,
		"ktime": time.Now(),
		"k[]struct": []subkstruct{{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     true,
			Ktime:     time.Now(),
			Kduration: time.Second,
		}},
		"k[]string": []string{"test1", "test2", "test3"},
		"kstruct": subkstruct{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     false,
			Ktime:     time.Now(),
			Kduration: time.Second,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClient_GetSingleValue(t *testing.T) {
	k := ""
	err := client.GetSingleValue(context.Background(), "k", &k)
	if err != nil {
		t.Error(err)
	}
	t.Log(k)

	k1 := 0
	err = client.GetSingleValue(context.Background(), "k1", &k1)
	if err != nil {
		t.Error(err)
	}
	t.Log(k1)

	var k1_1 float64
	err = client.GetSingleValue(context.Background(), "k1.1", &k1_1)
	if err != nil {
		t.Error(err)
	}
	t.Log(k1_1)

	var ktrue bool
	err = client.GetSingleValue(context.Background(), "ktrue", &ktrue)
	if err != nil {
		t.Error(err)
	}
	t.Log(ktrue)

	var ktime time.Time
	err = client.GetSingleValue(context.Background(), "ktime", &ktime)
	if err != nil {
		t.Error(err)
	}
	t.Log(ktime.String())

	var kduration time.Duration
	err = client.GetSingleValue(context.Background(), "kduration", &kduration)
	if err != nil {
		t.Error(err)
	}
	t.Log(kduration.String())

	var kstrings []string
	err = client.GetSingleValue(context.Background(), "k[]string", &kstrings)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstrings)

	var kints []int
	err = client.GetSingleValue(context.Background(), "k[]int", &kints)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kints)

	var kstructs []subkstruct
	err = client.GetSingleValue(context.Background(), "k[]struct", &kstructs)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstructs)

	var kstruct subkstruct
	err = client.GetSingleValue(context.Background(), "kstruct", &kstruct)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstruct)

	kmapstringinterface := make(map[string]interface{})
	err = client.GetSingleValue(context.Background(), "kmap[string]interface{}", &kmapstringinterface)
	if err != nil {
		t.Error(err)
	}
	t.Log("")
	for k, v := range kmapstringinterface {
		t.Logf("%+v:%+v", k, v)
	}
}

func TestClient_SetSliceValue(t *testing.T) {
	err := client.SetSliceValue(context.Background(), "k[]string", []string{"test1", "test2", "test3"})
	if err != nil {
		t.Error(err)
	}

	err = client.SetSliceValue(context.Background(), "k[]int", []int{1, 2, 3, 4, 5})
	if err != nil {
		t.Error(err)
	}

	err = client.SetSliceValue(context.Background(), "k[]struct", []subkstruct{{
		K:         "v",
		K1:        1,
		K1_1:      1.1,
		Ktrue:     true,
		Ktime:     time.Now(),
		Kduration: time.Second,
	}, {
		K:         "v1",
		K1:        2,
		K1_1:      2.1,
		Ktrue:     false,
		Ktime:     time.Now(),
		Kduration: time.Second * 3,
	}})
	if err != nil {
		t.Error(err)
	}
}

func TestClient_GetSliceValue(t *testing.T) {
	var kstrings []string
	err := client.GetSliceValue(context.Background(), "k[]string", &kstrings)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstrings)

	var kints [3]int
	err = client.GetSliceValue(context.Background(), "k[]int", &kints)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kints)

	var kstructs []subkstruct
	err = client.GetSliceValue(context.Background(), "k[]struct", &kstructs)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstructs)
}

func TestClient_SetStructValue(t *testing.T) {
	err := client.SetStructValue(context.Background(), "kstruct", kstruct{
		K:         "v",
		K1:        1,
		K1_1:      1.1,
		Ktrue:     true,
		Ktime:     time.Now(),
		Kduration: time.Second,
		Kstruct: subkstruct{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     true,
			Ktime:     time.Now(),
			Kduration: time.Second,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClient_GetStructValue(t *testing.T) {
	var kstruct kstruct
	err := client.GetStructValue(context.Background(), "kstruct", &kstruct)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstruct)
}

func TestClient_SetValue(t *testing.T) {
	err := client.SetValue(context.Background(), "k", "v")
	if err != nil {
		t.Error(err)
	}
	err = client.SetValue(context.Background(), "k1", 1)
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "k1.1", 1.1)
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "ktrue", true)
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "ktime", time.Now())
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "kduration", time.Second)
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "k[]string", []string{"test1", "test2", "test3"})
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "k[]int", []int{1, 2, 3, 4, 5})
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "k[]struct", []subkstruct{{
		K:         "v",
		K1:        1,
		K1_1:      1.1,
		Ktrue:     true,
		Ktime:     time.Now(),
		Kduration: time.Second,
	}})
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "kstruct", kstruct{
		K:         "v",
		K1:        1,
		K1_1:      1.1,
		Ktrue:     true,
		Ktime:     time.Now(),
		Kduration: time.Second,
		Kstruct: subkstruct{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     true,
			Ktime:     time.Now(),
			Kduration: time.Second,
		},
	})
	if err != nil {
		t.Error(err)
	}

	err = client.SetValue(context.Background(), "kmap[string]interface{}", map[string]interface{}{
		"k":     "v",
		"k1":    1,
		"k1.1":  1.1,
		"ktrue": true,
		"ktime": time.Now(),
		"k[]struct": []subkstruct{{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     true,
			Ktime:     time.Now(),
			Kduration: time.Second,
		}},
		"k[]string": []string{"test1", "test2", "test3"},
		"kstruct": subkstruct{
			K:         "v",
			K1:        1,
			K1_1:      1.1,
			Ktrue:     false,
			Ktime:     time.Now(),
			Kduration: time.Second,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClient_GetValue(t *testing.T) {
	k := ""
	err := client.GetValue(context.Background(), "k", &k)
	if err != nil {
		t.Error(err)
	}
	t.Log(k)

	k1 := 0
	err = client.GetValue(context.Background(), "k1", &k1)
	if err != nil {
		t.Error(err)
	}
	t.Log(k1)

	var k1_1 float64
	err = client.GetValue(context.Background(), "k1.1", &k1_1)
	if err != nil {
		t.Error(err)
	}
	t.Log(k1_1)

	var ktrue bool
	err = client.GetValue(context.Background(), "ktrue", &ktrue)
	if err != nil {
		t.Error(err)
	}
	t.Log(ktrue)

	var ktime time.Time
	err = client.GetValue(context.Background(), "ktime", &ktime)
	if err != nil {
		t.Error(err)
	}
	t.Log(ktime.String())

	var kduration time.Duration
	err = client.GetValue(context.Background(), "kduration", &kduration)
	if err != nil {
		t.Error(err)
	}
	t.Log(kduration.String())

	var kstrings []string
	err = client.GetValue(context.Background(), "k[]string", &kstrings)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstrings)

	var kints []int
	err = client.GetValue(context.Background(), "k[]int", &kints)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kints)

	var kstructs []subkstruct
	err = client.GetValue(context.Background(), "k[]struct", &kstructs)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstructs)

	var kstruct subkstruct
	err = client.GetValue(context.Background(), "kstruct", &kstruct)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", kstruct)

	kmapstringinterface := make(map[string]interface{})
	err = client.GetValue(context.Background(), "kmap[string]interface{}", &kmapstringinterface)
	if err != nil {
		t.Error(err)
	}
	t.Log("")
	for k, v := range kmapstringinterface {
		t.Logf("%+v:%+v", k, v)
	}
	kmapstringstring := make(map[string]string)
	err = client.GetValue(context.Background(), "kmap[string]interface{}", &kmapstringstring)
	if err != nil {
		t.Error(err)
	}
	t.Log("")
	for k, v := range kmapstringstring {
		t.Logf("%+v:%+v", k, v)
	}
}
