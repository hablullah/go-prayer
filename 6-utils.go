package prayer

import (
	"sort"
	"time"
)

func groupEmptyTimes(times []time.Time) [][]int {
	// Group the empty times
	lastIdx := -2
	var currentGroup []int
	var indexesGroups [][]int
	for i, t := range times {
		if t.IsZero() {
			if i > lastIdx+1 && len(currentGroup) > 0 {
				indexesGroups = append(indexesGroups, currentGroup)
				currentGroup = []int{i}
			} else {
				currentGroup = append(currentGroup, i)
			}
			lastIdx = i
		}
	}

	// If:
	// 1) the final group not empty,
	// 2) the final group end with the last index of `times`
	// 3) the first group start with 0
	// then merge the final group with the first one
	if len(currentGroup) > 0 {
		firstOfFirstGroup := -1
		if len(indexesGroups) > 0 && len(indexesGroups[0]) > 0 {
			firstOfFirstGroup = indexesGroups[0][0]
		}

		lastOfFinalGroup := currentGroup[len(currentGroup)-1]
		if firstOfFirstGroup == 0 && lastOfFinalGroup == len(times)-1 {
			indexesGroups[0] = append(currentGroup, indexesGroups[0]...)
			currentGroup = []int{} // empty the current group
		}
	}

	// If current group is not empty, save
	if len(currentGroup) > 0 {
		indexesGroups = append(indexesGroups, currentGroup)
	}

	// Sort index groups by its length
	sort.Slice(indexesGroups, func(i, j int) bool {
		return len(indexesGroups[i]) > len(indexesGroups[j])
	})

	return indexesGroups
}

// limitTimeDiff limit the difference between the `current` and `reference`. If `forward` is true,
// it means `reference` is behind the `current`.
func limitTimeDiff(curent, reference time.Time, maxDiff time.Duration, forward bool) time.Time {
	currentDiff := curent.Sub(reference).Abs()
	if currentDiff < maxDiff {
		return curent
	}

	if forward {
		if curent.Before(reference) {
			return reference.Add(-maxDiff)
		} else {
			return reference.Add(maxDiff)
		}
	} else {
		if curent.After(reference) {
			return reference.Add(maxDiff)
		} else {
			return reference.Add(-maxDiff)
		}
	}
}

func getRealIndex(idx int, list []time.Time) int {
	if idx < 0 {
		return len(list) - 1
	} else if idx > len(list)-1 {
		return 0
	} else {
		return idx
	}
}
