package examples

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/seanpar203/gobus"
)

/*

A multi stage event

Scrape Job -> Scrape Website -> Generate Report

*/

type ScrapeJobStatus int64

const (
	ScrapeJobCreated ScrapeJobStatus = iota
	ScrapeJobStarted
	ScrapeJobReportCreated
	ScrapeJobFinished
)

var (
	scrapeRW sync.RWMutex
	scrapeDB = map[string]ScrapeJob{}

	NEW_SCRAPE_JOB      = gobus.Event("new-scrape-job")
	START_SCRAPE_JOB    = gobus.Event("start-scrape-job")
	SCRAPE_JOB_COMPLETE = gobus.Event("generate-report")
)

type ScrapeJob struct {
	ID            string
	Website       string
	Status        ScrapeJobStatus
	ScrapeResults string
	ReportResults string
}

// writeJobToDB writes a new job to the DB.
//
// It takes a ScrapeJob as a parameter and returns an error.
func writeJobToDB(job ScrapeJob) (ScrapeJob, error) {

	scrapeRW.Lock()
	defer scrapeRW.Unlock()

	scrapeDB[job.ID] = ScrapeJob{
		ID:      uuid.NewString(),
		Website: job.Website,
		Status:  ScrapeJobCreated,
	}

	return scrapeDB[job.ID], nil
}

// getJobById retrieves a ScrapeJob from the scrapeDB map based on the provided ID.
//
// Parameters:
// - id: the ID of the ScrapeJob to retrieve.
//
// Returns:
// - ScrapeJob: the retrieved ScrapeJob.
// - error: an error if the ScrapeJob with the provided ID does not exist.
func getJobById(id string) (ScrapeJob, error) {
	scrapeRW.RLock()
	defer scrapeRW.RUnlock()

	job, ok := scrapeDB[id]

	if !ok {
		return ScrapeJob{}, errors.New("unable to create job")
	}

	return job, nil
}

// updateScrapeJob updates a ScrapeJob in the scrapeDB map.
//
// Parameters:
// - job: The ScrapeJob to be updated.
//
// Returns:
// - error: An error if the update operation fails.
func updateScrapeJob(job ScrapeJob) error {
	scrapeRW.Lock()
	defer scrapeRW.Unlock()

	scrapeDB[job.ID] = job

	return nil

}

// scrapeWebsite is a function that takes a ScrapeJob as a parameter and returns a string and an error.
//
// It actually goes out to a website and does some work.
// It returns the string "I did it!" and a nil error if successful.
func scrapeWebsite(job ScrapeJob) (string, error) {
	// Actually go out to website here and do some work.

	return "I did it!", nil
}

// generateReport generates a report based on the given ScrapeJob.
//
// It takes a ScrapeJob as a parameter and returns a string and an error.
func generateReport(job ScrapeJob) (string, error) {
	// Actually generate a report here.
	return "I generated the report", nil
}

// NewScrapeJobCreateDBRecord creates a new database record for a scrape job.
//
// The function takes a single argument `args` of any type. It first checks if
// the argument is of type `ScrapeJob`. If it's not, it returns an
// `InvalidArgError` with the details of the error. If it is of type `ScrapeJob`,
// it writes the job to the database using the `writeJobToDB` function.
//
// The function returns an error if there was an issue writing the job to the
// database.
func NewScrapeJobCreateDBRecord(args any) error {
	_, ok := args.(ScrapeJob)

	if !ok {
		return gobus.NewInvalidArgError(NEW_SCRAPE_JOB, "NewScrapeJobEvent", ScrapeJob{}, args)
	}

	job, err := writeJobToDB(args.(ScrapeJob))

	if err != nil {
		return err
	}

	return gobus.Emit(NEW_SCRAPE_JOB, job, nil)
}

// StartScrapeJobScrapeWebsite is a function that performs some action.
//
// It takes in a parameter called `args` of type `any`.
// It returns an error.
func StartScrapeJobScrapeWebsite(args any) error {
	_, ok := args.(ScrapeJob)

	if !ok {
		return gobus.NewInvalidArgError(START_SCRAPE_JOB, "StartScrapeJobEvent", ScrapeJob{}, args)
	}

	job := args.(ScrapeJob)

	job.Status = ScrapeJobStarted

	if err := updateScrapeJob(job); err != nil {
		return err
	}

	results, err := scrapeWebsite(job)

	if err != nil {
		return err
	}

	job.ScrapeResults = results

	if err = updateScrapeJob(job); err != nil {
		return err
	}

	return gobus.Emit(SCRAPE_JOB_COMPLETE, job, nil)
}

// GenerateReport generates a report for a given ScrapeJob.
//
// args: any - The arguments for the function.
// Returns: error - An error if the report generation fails.
func GenerateReport(args any) error {
	_, ok := args.(ScrapeJob)

	if !ok {
		return gobus.NewInvalidArgError(START_SCRAPE_JOB, "GenerateReport", ScrapeJob{}, args)
	}

	job := args.(ScrapeJob)

	job.Status = ScrapeJobReportCreated
	job.ReportResults = "Some results"

	return updateScrapeJob(job)
}

func ScrapeJobExample() {

	// Producer is on the outside
	//
	// Could be from an API request, a cron job, etc.
	err := gobus.Emit(NEW_SCRAPE_JOB, ScrapeJob{Website: "example.com"}, nil)

	if err != nil {
		fmt.Printf("error emitting event: %s", err)
	}

	fmt.Println("Done with scrape job!")
}
