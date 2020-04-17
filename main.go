package main

import (
	"fmt"
	"github.com/maxim-kuderko/mosyne/entities"
	store "github.com/maxim-kuderko/mosyne/store"
)

func main() {
	store := store.NewInMemStore(&store.Options{})
	resp := store.ZSet(entities.ZSetRequest{
		Key:   `a`,
		Value: `valueeeee`,
		Score: 151621.1,
	})
	fmt.Println(resp)

}
