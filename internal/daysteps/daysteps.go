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

// parsePackage парсит строку с данными о шагах и времени
func parsePackage(data string) (int, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 2 {
		return 0, 0, errors.New("Неверно заданы данные: неверный формат")
	}
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, 0, errors.New("Неверно заданы шаги: неверный формат")
	}
	if steps <= 0 {
		return 0, 0, errors.New("Неверно заданы шаги: они должны быть больше 0")
	}
	duration, err := time.ParseDuration(slice[1])
	if err != nil {
		return 0, 0, errors.New("Неверно задано время: неправильный формат")
	}
	if duration <= 0 {
		return 0, 0, errors.New("Неверно задано время: оно должно быть чем 0")
	}
	return steps, duration, nil
}

// DayActionInfo обрабатывает данные и возвращает информацию о прогулке
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка:%s", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка расчета калорий:%s", err)
		return ""
	}

	info := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)

	return info
}
