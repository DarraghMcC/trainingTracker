package main

import (
	"fmt"
	"strings"
	"time"
	"google.golang.org/api/calendar/v3"
	"strconv"
)

func Filter(vs []*calendar.Event, f func(*calendar.Event) bool) []*calendar.Event {
	vsf := make([]*calendar.Event, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func isBjjEvent(event *calendar.Event) bool {
	return strings.Contains(strings.ToUpper(event.Summary), "BJJ")
}

func getTimeOfEvent(event *calendar.Event) (t time.Time) {
	t, err := time.Parse(time.RFC3339, event.Start.DateTime)
	if event.Start.DateTime == "" {
		t, err = time.Parse("2006-01-02", event.Start.Date)
	}
	if err != nil {
		fmt.Println("failure is " + event.Start.Date)
		panic(fmt.Sprintf("Failure to get time of event: %v\n", err))
	}
	return t

}

func getWeeklyBreakdownOfEventType(eventItems []*calendar.Event, filterCondition func(*calendar.Event) bool) []WeekYear {
	weekMap := getAllWeekYearsBetween(getTimeOfEvent(eventItems[0]), getTimeOfEvent(eventItems[len(eventItems)-1]))

	for _, event := range Filter(eventItems, filterCondition) {
		week := getWeekYear(getTimeOfEvent(event)).String()

		fmt.Println("week year is " + week)
		weekMap[week].increment()
	}

	//map just for quick look ups, discard it now
	weeks := make([]WeekYear, 0, len(weekMap))

	for _ , week := range weekMap {
		weeks = append(weeks, *week)
	}
	return weeks
}

func getAllWeekYearsBetween(startDate time.Time, endDate time.Time) map[string]*WeekYear {
	currentDate := startDate
	weeks := make(map[string]*WeekYear)
	for currentDate.Before(endDate) {
		wy := getWeekYear(currentDate)
		weeks[wy.String()] = &wy
		currentDate = currentDate.AddDate(0, 0, 7)
	}
	return weeks
}

func getWeekYear(time time.Time) WeekYear {
	year, week := time.ISOWeek()
	var wY WeekYear
	wY.Week = week
	wY.Year = year
	wY.Count = 0.0
	return wY
}

type WeekYear struct {
	Week int
	Year int

	Count float64
}

func (wy WeekYear) String() string {
	return strconv.Itoa(wy.Week) + "  " + strconv.Itoa(wy.Year)
}

func (wy *WeekYear) increment() {
	wy.Count = wy.Count + 1.0
}
