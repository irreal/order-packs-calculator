package orders

import (
	"fmt"
	"math"
	"sort"

	"github.com/irreal/order-packs/models"
)

type PackingCalculation struct {
	Packs      map[models.Pack]int
	TotalItems int
	TotalPacks int
}

// Calculates the packs to be used given these rules:
//
// 1. Only whole packs can be sent. Packs cannot be broken open.
// 2. Within the constraints of Rule 1 above, send out the least amount of items to fulfil the order.
// 3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfil each order.
// note: rule 2 takes precedence over rule 3
func CalculatePack(availablePacks []models.Pack, requestedCount int) (*PackingCalculation, error) {
	if requestedCount <= 0 {
		return nil, fmt.Errorf("requested count must be greater than 0")
	}
	if len(availablePacks) == 0 {
		return nil, fmt.Errorf("no packs available to fulfill the order")
	}

	sort.Slice(availablePacks, func(i, j int) bool {
		return availablePacks[i] < availablePacks[j]
	})

	maxPackSize := int(availablePacks[len(availablePacks)-1])

	//Let's solve it dynamically for each count up from 0 to requested count
	//the overshoot can be at most the size of the max pack
	maxSize := requestedCount + maxPackSize
	minItems := make([]int, maxSize+1)
	minPacks := make([]int, maxSize+1)
	lastPackUsed := make(map[int]models.Pack)

	for i := 1; i <= maxSize; i++ {
		minItems[i] = math.MaxInt32
		minPacks[i] = math.MaxInt32
	}

	// base case, for 0 items, we need 0 packs
	minItems[0] = 0
	minPacks[0] = 0

	// find the solution for all order counts up to target (+ max pack)
	for i := 1; i <= maxSize; i++ {

		// consider all packs
		for _, packSize := range availablePacks {

			// the current pack is a candidate to solve the target. it reaches or overshoots the target, and we already solved the previous sub-target
			if i >= int(packSize) && minItems[i-int(packSize)] != math.MaxInt32 {

				// this candidate pack gives these results
				newTotalItems := minItems[i-int(packSize)] + int(packSize)
				newTotalPacks := minPacks[i-int(packSize)] + 1

				// rule 2, minimize the total number of items
				if newTotalItems < minItems[i] {
					minItems[i] = newTotalItems
					minPacks[i] = newTotalPacks
					lastPackUsed[i] = packSize
				} else if newTotalItems == minItems[i] {
					// we found a solution as good as an existing for rule 2, check if it is better in terms of rule 3, the count of packs
					if newTotalPacks < minPacks[i] {
						minPacks[i] = newTotalPacks
						lastPackUsed[i] = packSize
					}
				}
			}
		}
	}

	// first solution that equals or overshoots the requestedCount
	optimalCount := -1
	for i := requestedCount; i <= maxSize; i++ {
		if minItems[i] != math.MaxInt32 {
			optimalCount = i
			break
		}
	}

	// as long as we have some packs, in theory, we always have a solution. this shouldn't happen with the validation checks we have at the top
	if optimalCount == -1 {
		return nil, fmt.Errorf("no valid combination found")
	}

	// we now know the optimal solution. reconstruct the exact packs and counts
	finalPacks := make(map[models.Pack]int)

	// start with target count and reconstruct the packs, relying on the memory of packs used for computed sub-amounts
	remainingCount := optimalCount
	for remainingCount > 0 {
		pack := lastPackUsed[remainingCount]
		finalPacks[pack] = finalPacks[pack] + 1

		remainingCount -= int(pack)
	}

	// yay, create the response struct
	return &PackingCalculation{
		Packs:      finalPacks,
		TotalItems: minItems[optimalCount],
		TotalPacks: minPacks[optimalCount],
	}, nil
}
