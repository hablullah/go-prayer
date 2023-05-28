package prayer

import "time"

type abnormalRange struct {
	Start   time.Time
	End     time.Time
	Indexes []int
}

func (ar abnormalRange) IsEmpty() bool {
	return len(ar.Indexes) == 0
}

func extractAbnormalSchedules(schedules []Schedule) (abnormalSummer, abnormalWinter abnormalRange) {
	// If there are no schedules, return empty
	if len(schedules) == 0 {
		return
	}

	// Loop each schedule
	lastAbnormalIdx := -2
	var ranges []abnormalRange
	var currentRange abnormalRange
	for i, s := range schedules {
		// If schedule is normal, skip
		if s.IsNormal {
			continue
		}

		// If current schedule is continuation of the previous, just add the index
		if i == lastAbnormalIdx+1 {
			currentRange.End = s.Zuhr
			currentRange.Indexes = append(currentRange.Indexes, i)
		} else {
			// If current range is not empty, save it to range list
			if len(currentRange.Indexes) > 0 {
				ranges = append(ranges, currentRange)
			}

			// Re-initiate current range
			currentRange = abnormalRange{
				Start:   s.Zuhr,
				End:     s.Zuhr,
				Indexes: []int{i},
			}
		}

		lastAbnormalIdx = i
	}

	// Handle leftover range
	if len(currentRange.Indexes) > 0 {
		// Check if we can merge the leftover range to the first range:
		// - there are more than one range
		// - the first index of first range = zero
		// - the last index of current range = last schedule index
		// If we can't merge it, just append current range to the end
		if len(ranges) > 0 &&
			ranges[0].Indexes[0] == 0 &&
			currentRange.Indexes[len(currentRange.Indexes)-1] == len(schedules)-1 {
			ranges[0].Start = currentRange.Start.AddDate(-1, 0, 0) // move start to last year
			ranges[0].Indexes = append(currentRange.Indexes, ranges[0].Indexes...)
		} else {
			ranges = append(ranges, currentRange)
		}
	}

	// At this point we at most only have two abnormal periods: one for summer and one for winter
	for _, sr := range ranges {
		// Extract months in range
		months := make(map[int]int)
		for tmp := sr.Start; tmp.Before(sr.End.AddDate(0, 1, 0)); tmp = tmp.AddDate(0, 1, 0) {
			months[int(tmp.Month())] = 1
		}

		// Calculate score:
		// - Summer is in June, July and August
		// - Winter is in December, January and February
		summerScore := months[6] + months[7] + months[8]
		winterScore := months[12] + months[1] + months[2]

		switch {
		case summerScore == 3, summerScore > winterScore:
			abnormalSummer = sr
		case winterScore == 3, winterScore > summerScore:
			abnormalWinter = sr
		}
	}

	return
}
