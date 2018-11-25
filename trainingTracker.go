package main

import (
	"log"
	"strconv"
	"sort"
	"github.com/wcharczuk/go-chart"
	"bytes"
	"os"
)


func main() {

	eventItems := getEventsFromCalender("Trained")
	weeks := getWeeklyBreakdownOfEventType(eventItems, isBjjEvent)

	drawGraph(transformToContinuousSeries(weeks))
}

func drawGraph(ts1 chart.ContinuousSeries){

	graph := chart.Chart{

		XAxis: chart.XAxis{
			Name:           "The XAxis",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(),
			ValueFormatter: chart.TimeMinuteValueFormatter, //TimeHourValueFormatter,
		},

		YAxis: chart.YAxis{
			Name:      "The YAxis",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},

		Series: []chart.Series{
			ts1,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	log.Printf("DRAWING...")
	err := graph.Render(chart.PNG, buffer)
	log.Printf("DONE!!")
	if err != nil {
		log.Fatal(err)
	}

	fo, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	if _, err := fo.Write(buffer.Bytes()); err != nil {
		panic(err)
	}

}

func transformToContinuousSeries(weeks []WeekYear) chart.ContinuousSeries{
	currentXAxis := 0.0
	xAxis := make([]float64, 0, len(weeks))
	yAxis := make([]float64, 0, len(weeks))
	log.Printf("length before" + strconv.Itoa(len(weeks)))

	sort.Slice(weeks, func(i, j int) bool {
		if weeks[i].Year < weeks[j].Year {
			return true
		}
		if weeks[i].Year > weeks[j].Year {
			return false
		}
		if weeks[i].Week < weeks[j].Week {
			return true
		}
		return weeks[i].Week > weeks[j].Year
	})


	for _, week := range weeks {

		xAxis = append(xAxis, currentXAxis)
		yAxis = append(yAxis, week.Count)
		currentXAxis ++
		log.Printf("graphing : " + strconv.FormatFloat(currentXAxis, 'f', 6, 64) + ":" + week.String() + ": " + strconv.FormatFloat(week.Count, 'f', 6, 64))
	}

	log.Printf("x axis length after" + strconv.Itoa(len(xAxis)))
	log.Printf("got AXIS:" + strconv.FormatFloat(currentXAxis, 'f', 6, 64))
	return chart.ContinuousSeries{ //TimeSeries{
		Name:    "Time Series",
		Style:   chart.StyleShow(),
		XValues: xAxis,
		YValues: yAxis,
	}
}