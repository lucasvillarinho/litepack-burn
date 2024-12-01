package main

import (
	"fmt"

	"github.com/lucasvillarinho/litepack-burn/attacker"
)

func main() {
	cacheAtk := attacker.NewCacheAttacker()

	if err := cacheAtk.Attack(); err != nil {
		fmt.Print(err)
		return
	}
}
