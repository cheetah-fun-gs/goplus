package log4go

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/log4go"
)

// RotateInterval 滚动间隔
type RotateInterval int

// 常量
const (
	RotateIntervalDaily RotateInterval = iota
	RotateIntervalHour12
	RotateIntervalHour8
	RotateIntervalHour6
	RotateIntervalHour4
	RotateIntervalHour3
	RotateIntervalHour2
	RotateIntervalHourly
	RotateIntervalMin30
	RotateIntervalMin20
	RotateIntervalMin15
	RotateIntervalMin10
)

// TimeFileLogWriter This log writer sends output to a file
type TimeFileLogWriter struct {
	rec chan *log4go.LogRecord

	// 日志名称
	fname string

	// The opened file
	filename string
	file     *os.File

	// The logging format
	format string

	// File header/trailer
	header, trailer string

	curlines int
	cursize  int

	rotate     RotateInterval
	rotateCurr string
}

// LogWrite This is the FileLogWriter's output method
func (w *TimeFileLogWriter) LogWrite(rec *log4go.LogRecord) {
	w.rec <- rec
}

// Close Close
func (w *TimeFileLogWriter) Close() {
	close(w.rec)
	w.file.Sync()
}

// NewTimeFileLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewTimeFileLogWriter(fname string, rotate RotateInterval) *TimeFileLogWriter {
	w := &TimeFileLogWriter{
		rec:      make(chan *log4go.LogRecord, log4go.LogBufferLength),
		fname:    fname,
		filename: fname,
		format:   "[%D %T] [%L] (%S) %M",
		rotate:   rotate,
	}

	// open the file for the first time
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
		return nil
	}

	go func() {
		defer func() {
			if w.file != nil {
				fmt.Fprint(w.file, log4go.FormatLogRecord(w.trailer, &log4go.LogRecord{Created: time.Now()}))
				w.file.Close()
			}
		}()

		for {
			select {
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				if err := w.intRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}

				// Perform the write
				n, err := fmt.Fprint(w.file, log4go.FormatLogRecord(w.format, rec))
				if err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}

				// Update the counts
				w.curlines++
				w.cursize += n
			}
		}
	}()

	return w
}

func (w *TimeFileLogWriter) currFilename() (filename, rotateCurr string) {
	if w.rotate == 0 {
		filename = w.fname
		return
	}

	now := time.Now()
	switch w.rotate {
	case RotateIntervalDaily:
		rotateCurr = now.Format("2006-01-02")
	case RotateIntervalHour12:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02"), now.Hour()/12*12)
	case RotateIntervalHour8:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02"), now.Hour()/8*8)
	case RotateIntervalHour6:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02"), now.Hour()/6*6)
	case RotateIntervalHour4:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02"), now.Hour()/4*4)
	case RotateIntervalHour3:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02"), now.Hour()/3*3)
	case RotateIntervalHour2:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02"), now.Hour()/2*2)
	case RotateIntervalHourly:
		rotateCurr = now.Format("2006-01-02_15")
	case RotateIntervalMin30:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02_15"), now.Minute()/30*30)
	case RotateIntervalMin20:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02_15"), now.Minute()/20*20)
	case RotateIntervalMin15:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02_15"), now.Minute()/15*15)
	case RotateIntervalMin10:
		rotateCurr = fmt.Sprintf("%s_%02d", now.Format("2006-01-02_15"), now.Minute()/10*10)
	default:
		rotateCurr = now.Format("2006-01-02")
	}
	filename = w.fname + "." + rotateCurr
	return
}

// If this is called in a threaded context, it MUST be synchronized
func (w *TimeFileLogWriter) intRotate() error {
	filename, rotateCurr := w.currFilename()
	if w.file != nil && w.rotateCurr == rotateCurr {
		return nil
	}

	// Close any log file that may be open
	if w.file != nil {
		fmt.Fprint(w.file, log4go.FormatLogRecord(w.trailer, &log4go.LogRecord{Created: time.Now()}))
		w.file.Close()
	}

	// Open the log file
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	w.filename = filename
	w.rotateCurr = rotateCurr
	w.file = fd

	now := time.Now()
	fmt.Fprint(w.file, log4go.FormatLogRecord(w.header, &log4go.LogRecord{Created: now}))

	// initialize rotation values
	w.curlines = 0
	w.cursize = 0

	return nil
}

// SetFormat Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *TimeFileLogWriter) SetFormat(format string) *TimeFileLogWriter {
	w.format = format
	return w
}

// SetHeadFoot Set the logfile header and footer (chainable).  Must be called before the first log
// message is written.  These are formatted similar to the FormatLogRecord (e.g.
// you can use %D and %T in your header/footer for date and time).
func (w *TimeFileLogWriter) SetHeadFoot(head, foot string) *TimeFileLogWriter {
	w.header, w.trailer = head, foot
	if w.curlines == 0 {
		fmt.Fprint(w.file, log4go.FormatLogRecord(w.header, &log4go.LogRecord{Created: time.Now()}))
	}
	return w
}
