package store

import (
	"github.com/maxim-kuderko/mosyne/entities"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func TestInMem_ZSet(t *testing.T) {
	var tests = []struct {
		name    string
		setReq  entities.ZSetRequest
		setResp entities.ZSetResponse
		getReq  entities.ZGetRequest
		getResp entities.ZGetResponse
	}{
		{
			name: `simple`,
			setReq: entities.ZSetRequest{
				Key:   "a",
				Value: `test`,
				Score: 1,
			},
			setResp: entities.ZSetResponse{
				Value: `test`,
				Error: nil,
			},
			getReq: entities.ZGetRequest{
				Key:      "a",
				ScoreMin: 0,
				ScoreMax: 1,
			},
			getResp: entities.ZGetResponse{
				Values: []entities.ZGetStruct{{Value: `test`, Score: 1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInMemStore(&Options{})
			if got := i.ZSet(tt.setReq); !reflect.DeepEqual(got, tt.setResp) && got.Error == nil {
				t.Errorf("ZSet() = %v, want %v", got, tt.setResp)
			}
			if got := i.ZGet(tt.getReq); !reflect.DeepEqual(got, tt.getResp) && got.Error == nil {
				t.Errorf("ZGet() = %v, want %v", got, tt.setResp)
			}
		})
	}
}

func TestInMem_ZSet_Cuncurrent(t *testing.T) {
	type test struct {
		name    string
		setReq  entities.ZSetRequest
		setResp entities.ZSetResponse
		getReq  entities.ZGetRequest
		getResp entities.ZGetResponse
	}
	var tests = make([]test, 0, 1e5)
	for i := 0; i < 1e5; i++ {
		key := strconv.Itoa(i)
		val := `test-` + strconv.Itoa(i)
		score := rand.Float64()
		tests = append(tests, test{
			name: key,
			setReq: entities.ZSetRequest{
				Key:   key,
				Value: val,
				Score: score,
			},
			setResp: entities.ZSetResponse{
				Value: val,
				Error: nil,
			},
			getReq: entities.ZGetRequest{
				Key:      key,
				ScoreMin: 0,
				ScoreMax: score,
			},
			getResp: entities.ZGetResponse{
				Values: []entities.ZGetStruct{{Value: val, Score: score}},
			},
		})
	}
	i := NewInMemStore(&Options{})
	wg := sync.WaitGroup{}
	wg.Add(len(tests))
	for _, tt := range tests {

		go func(tt test) {
			defer func() {
				wg.Done()
			}()
			t.Run(tt.name, func(t *testing.T) {

				if got := i.ZSet(tt.setReq); !reflect.DeepEqual(got, tt.setResp) && got.Error == nil {
					t.Errorf("ZSet() = %v, want %v", got, tt.setResp)
				}
				if got := i.ZGet(tt.getReq); !reflect.DeepEqual(got, tt.getResp) && got.Error == nil {
					t.Errorf("ZGet() = %v, want %v", got, tt.setResp)
				}
			})
		}(tt)

	}
	wg.Wait()
}

type test struct {
	name    string
	setReq  entities.ZSetRequest
	setResp entities.ZSetResponse
	getReq  entities.ZGetRequest
	getResp entities.ZGetResponse
}

func BenchmarkInMem_ZSet(b *testing.B) {
	b.ReportAllocs()
	c := int(b.N)
	var tests = make([]test, 0, c)
	for i := 0; i < c; i++ {
		key := strconv.Itoa(i)
		val := `test-` + strconv.Itoa(i)
		score := rand.Float64()
		tests = append(tests, test{
			name: key,
			setReq: entities.ZSetRequest{
				Key:   key,
				Value: val,
				Score: score,
			},
			setResp: entities.ZSetResponse{
				Value: val,
				Error: nil,
			},
			getReq: entities.ZGetRequest{
				Key:      key,
				ScoreMin: 0,
				ScoreMax: score,
			},
			getResp: entities.ZGetResponse{
				Values: []entities.ZGetStruct{{Value: val, Score: score}},
			},
		})
	}
	i := NewInMemStore(&Options{})

	for _, tt := range tests {
		if got := i.ZSet(tt.setReq); !reflect.DeepEqual(got, tt.setResp) && got.Error == nil {
			b.Errorf("ZSet() = %v, want %v", got, tt.setResp)
		}
		if got := i.ZGet(tt.getReq); !reflect.DeepEqual(got, tt.getResp) && got.Error == nil {
			b.Errorf("ZGet() = %v, want %v", got, tt.setResp)
		}

	}

}
