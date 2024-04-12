package console

import (
	"fmt"
	"sync"
	"time"
)

type TerminalInterface interface {
	AddDomain(domain string)
	AddStatus(status string)
	AddStarted()
	AddFinished()
	AddAttachments()
	AddDownloading()
	AddDownloaded()
	AddErrors(err string)
}

type Terminal struct {
	out   ouput
	mutex sync.Mutex
}

type ouput struct {
	domain      string
	status      string
	started     int
	scanning    int
	finished    int
	attachments int
	downloaded  int
	downloading int
	errors      int
	start       time.Time
}

func New() *Terminal {
	// Clear the console // call this only once
	fmt.Print("\033[2J")
	return &Terminal{
		out: ouput{
			status: "Starting",
			start:  time.Now(),
		},
	}
}

// Println prints the output, clearing the entire console before printing
func (t *Terminal) print() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Move the cursor to the top left corner
	fmt.Print("\033[H")

	// Status shouldn't be longer than 80 characters
	if len(t.out.status) > 78 {
		t.out.status = t.out.status[:78] + ".."
	}

	// pad the status with spaces
	for len(t.out.status) < 80 {
		t.out.status += " "
	}

	// Print each line separately, moving the cursor to the correct position before each print
	fmt.Printf("\033[1;1HDomain: %s\n", t.out.domain)
	fmt.Printf("\033[2;1HStatus: %s\n", t.out.status)
	fmt.Printf("\033[3;1HTime Started: %s\n", t.out.start)
	fmt.Printf("\033[4;1HPages Started: %d\n", t.out.started)
	fmt.Printf("\033[5;1HScanning: %d\n", t.out.scanning)
	fmt.Printf("\033[6;1HPages Finished: %d\n", t.out.finished)
	fmt.Printf("\033[7;1HAttachments to Download: %d\n", t.out.attachments)
	fmt.Printf("\033[8;1HDownloading %d\n", t.out.downloading)
	fmt.Printf("\033[9;1HAttachments Downloaded: %d\n", t.out.downloaded)
	fmt.Printf("\033[10;1HErrors: %d\n", t.out.errors)
	fmt.Printf("\033[11;1HTime: %s		\n", time.Since(t.out.start))
}

func (t *Terminal) AddDomain(domain string) {
	t.out.domain = domain
	t.print()
}

func (t *Terminal) AddStatus(status string) {
	t.out.status = status
	t.print()
}

func (t *Terminal) AddStarted() {
	t.out.started++
	t.out.scanning++
	t.print()
}

func (t *Terminal) AddFinished() {
	t.out.scanning--
	t.out.finished++
	t.print()
}

func (t *Terminal) AddAttachments() {
	t.out.attachments++
	t.print()
}

func (t *Terminal) AddDownloaded() {
	t.out.downloaded++
	t.print()
}

func (t *Terminal) AddDownloading() {
	t.out.downloading++
	t.print()
}

func (t *Terminal) AddErrors(err string) {
	t.out.status = err
	t.out.errors++
	t.print()
}
