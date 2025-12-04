package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага. unused?
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	params := strings.Split(data, ",")

	if len(params) != 3 {
		return 0, "", 0, errors.New("invalid data params in :parseTraining:spentcalories")
	}

	steps, err := strconv.Atoi(params[0]) // trim?
	if err != nil {
		return 0, "", 0, errors.New("unable to convert steps to integer in :parseTraining")
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps cannot be 0 in :parseTraining")
	}

	activity := params[1]
	if len(activity) == 0 {
		return 0, "", 0, errors.New("activity cannot be empty in :parseTraining")
	}

	duration, err := time.ParseDuration(params[2])
	if err != nil {
		return 0, "", 0, errors.New("unable to convert duration to time in :parseTraining")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be positive in :parseTraining")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	return height * stepLengthCoefficient * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	curDist := distance(steps, height)

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), curDist, meanSpeed(steps, height, duration), calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid parameters in :RunningSpentCalories")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	return weight * meanSpeed * duration.Minutes() / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid parameters in :WalkingSpentCalories")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	return weight * meanSpeed * duration.Minutes() / minInH * walkingCaloriesCoefficient, nil
}
