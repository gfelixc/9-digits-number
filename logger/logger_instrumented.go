package logger

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

const reportFrequency = 10 * time.Second

type LoggerInstrumented struct {
	reportTicker                    *time.Ticker
	uniqueNumbersFromLastReport     uint32
	duplicatedNumbersFromLastReport uint32
	uniqueNumbersCounter            uint32
	duplicatedNumbersCounter        uint32
	*Logger
}

func NewLoggerInstrumented(logger *Logger) *LoggerInstrumented {
	l := &LoggerInstrumented{
		reportTicker: time.NewTicker(reportFrequency),
		Logger:       logger,
	}
	return l
}

func (i *LoggerInstrumented) OnlyNumbers(s string) error {
	err := i.Logger.OnlyNumbers(s)

	if err == nil {
		atomic.AddUint32(&i.uniqueNumbersCounter, 1)
	}

	if errors.Is(err, ErrDuplicatedNumber) {
		atomic.AddUint32(&i.duplicatedNumbersCounter, 1)
	}

	return err
}

func (i *LoggerInstrumented) Shutdown() {
	i.Logger.Shutdown()
	i.reportTicker.Stop()
	i.printReport()
}

func (i *LoggerInstrumented) printReport() {
	u := atomic.LoadUint32(&i.uniqueNumbersCounter)
	deltaUniqueNumbers := u - i.uniqueNumbersFromLastReport
	i.uniqueNumbersFromLastReport = u

	d := atomic.LoadUint32(&i.duplicatedNumbersCounter)
	deltaDuplicatedNumbers := d - i.duplicatedNumbersFromLastReport
	i.duplicatedNumbersFromLastReport = d

	_, _ = fmt.Fprintf(
		os.Stdout,
		"Received %d unique numbers, %d duplicates. Unique total: %d\n",
		deltaUniqueNumbers,
		deltaDuplicatedNumbers,
		atomic.LoadUint32(&i.uniqueNumbersCounter),
	)
}
