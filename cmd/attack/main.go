package main

import (
	"log"

	"github.com/lucasvillarinho/litepack-burn/attacker"
)

func main() {
	cacheAtk := attacker.NewCacheAttacker()

	if err := cacheAtk.AttackCacheSet(); err != nil {
		log.Fatalf("Error running benchmark: %v", err)
	}
}
