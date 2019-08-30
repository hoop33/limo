package xog

import "github.com/theplant/cldr"

var calendar = cldr.Calendar{
	Formats: cldr.CalendarFormats{
		Date:     cldr.CalendarDateFormat{Full: "EEEE, d MMMM y", Long: "d MMMM y", Medium: "d MMM y", Short: "dd/MM/y"},
		Time:     cldr.CalendarDateFormat{Full: "h:mm:ss a zzzz", Long: "h:mm:ss a z", Medium: "h:mm:ss a", Short: "h:mm a"},
		DateTime: cldr.CalendarDateFormat{},
	},
	FormatNames: cldr.CalendarFormatNames{
		Months: cldr.CalendarMonthFormatNames{
			Abbreviated: cldr.CalendarMonthFormatNameValue{Jan: "Jan", Feb: "Feb", Mar: "Mar", Apr: "Apu", May: "Maa", Jun: "Juu", Jul: "Jul", Aug: "Agu", Sep: "Seb", Oct: "Oki", Nov: "Nov", Dec: "Des"},
			Narrow:      cldr.CalendarMonthFormatNameValue{Jan: "J", Feb: "F", Mar: "M", Apr: "A", May: "M", Jun: "J", Jul: "J", Aug: "A", Sep: "S", Oct: "O", Nov: "N", Dec: "D"},
			Short:       cldr.CalendarMonthFormatNameValue{},
			Wide:        cldr.CalendarMonthFormatNameValue{Jan: "Janwaliyo", Feb: "Febwaliyo", Mar: "Marisi", Apr: "Apuli", May: "Maayi", Jun: "Juuni", Jul: "Julaayi", Aug: "Agusito", Sep: "Sebuttemba", Oct: "Okitobba", Nov: "Novemba", Dec: "Desemba"},
		},
		Days: cldr.CalendarDayFormatNames{
			Abbreviated: cldr.CalendarDayFormatNameValue{Sun: "Sabi", Mon: "Bala", Tue: "Kubi", Wed: "Kusa", Thu: "Kuna", Fri: "Kuta", Sat: "Muka"},
			Narrow:      cldr.CalendarDayFormatNameValue{Sun: "S", Mon: "B", Tue: "B", Wed: "S", Thu: "K", Fri: "K", Sat: "M"},
			Short:       cldr.CalendarDayFormatNameValue{},
			Wide:        cldr.CalendarDayFormatNameValue{Sun: "Sabiiti", Mon: "Balaza", Tue: "Owokubili", Wed: "Owokusatu", Thu: "Olokuna", Fri: "Olokutaanu", Sat: "Olomukaaga"},
		},
		Periods: cldr.CalendarPeriodFormatNames{
			Abbreviated: cldr.CalendarPeriodFormatNameValue{},
			Narrow:      cldr.CalendarPeriodFormatNameValue{},
			Short:       cldr.CalendarPeriodFormatNameValue{},
			Wide:        cldr.CalendarPeriodFormatNameValue{AM: "Munkyo", PM: "Eigulo"},
		},
	},
}
