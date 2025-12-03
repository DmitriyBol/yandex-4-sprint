package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	params := strings.Split(data, ",")

	if len(params) != 2 {
		return 0, 0, errors.New("invalid data params in :parsePackage:daysteps")
	}

	steps, err := strconv.Atoi(params[0])
	if err != nil {
		return 0, 0, errors.New("unable to convert steps to integer in :parsePackage")
	}
	if steps <= 0 {
		return 0, 0, errors.New("steps cannot be 0 in :parsePackage")
	}

	duration, err := time.ParseDuration(params[1])
	if err != nil {
		return 0, 0, errors.New("unable to convert duration to time in :parsePackage")
	}
	if duration <= 0 {
		return 0, 0, errors.New("duration must be positive in :parsePackage")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distance := float64(steps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return "error to calculate calories in :DayActionInfo"
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
