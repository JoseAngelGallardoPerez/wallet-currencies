package workers

import (
	"fmt"
	"time"
)

// Offset for CET(UTC +1) time in second
const cetOffset = 60 * 60

// Converts "4:00PM" (CET) to local time
// in format 16:00 (with offset) needed for workers
func FixedCetTimeToLocal(kitchenTime string) (result string, err error) {
	fixedCetLocation := time.FixedZone("CET", cetOffset)
	return convertFromLocation(fixedCetLocation, kitchenTime)
}

func DynamicCetToLocal(kitchenTime string) (cetTime string, err error) {
	cetLocation, _ := time.LoadLocation("CET")
	return convertFromLocation(cetLocation, kitchenTime)
}

func convertFromLocation(loc *time.Location, kitchenTime string) (result string, err error) {
	timeResult, err := time.ParseInLocation(time.Kitchen, kitchenTime, loc)
	if err == nil {
		_, currentOffset := time.Now().Zone()
		_, locationOffset := time.Now().In(loc).Zone()
		duration := time.Duration(currentOffset-locationOffset) * time.Second
		timeResult = timeResult.Add(duration)
		result = fmt.Sprintf("%02d:%02d", timeResult.Hour(), timeResult.Minute())
	}
	return
}
