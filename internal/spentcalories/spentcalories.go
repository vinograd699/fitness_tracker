package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку с данными о тренировке
func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid input format, expected 'steps,activity,duration'")
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, errors.New("invalid steps format")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("steps must be positive")
	}

	activity := strings.TrimSpace(parts[1])
	if activity == "" {
		return 0, "", 0, errors.New("activity cannot be empty")
	}

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, errors.New("invalid duration format")
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}

	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в километрах
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	return distanceKm
}

// meanSpeed вычисляет среднюю скорость в километрах в час
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	hours := duration.Hours()

	speed := dist / hours

	return speed
}

// TrainingInfo обрабатывает данные и возвращает информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	if weight <= 0 || height <= 0 {
		return "", errors.New("вес и рост должны быть положительными")
	}

	var distanceKm, speed, calories float64
	var errCalc error

	switch strings.ToLower(activity) {
	case "бег":
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, errCalc = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба":
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, errCalc = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if errCalc != nil {
		return "", errCalc
	}

	durationHours := fmt.Sprintf("%.2f", duration.Hours())

	info := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %s ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		durationHours,
		distanceKm,
		speed,
		calories,
	)

	return info, nil
}

// RunningSpentCalories вычисляет количество затраченных калорий при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories вычисляет количество затраченных калорий при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()

	baseCalories := (weight * speed * durationMinutes) / minInH

	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}
