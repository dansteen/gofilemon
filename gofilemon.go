package main

import "os"
import "fmt"
import "sync"
import "runtime"
import "github.com/jessevdk/go-flags"
import "github.com/hpcloud/tail"
import "regexp"

//import "reflect"

/* ------------ Process Command Line ----------- */
// these help us process our command line
type search struct {
	String regexp.Regexp `the string to look for`
	Status string        `whether to succeed or fail`
}

var opts struct {
	// successes
	Successes []map[string]string `short:"s" long:"succeed" description:"Monitor the file for the supplied regex, and exit with success if found.  Use the following format: <file to monitor>:<regex to look for>"`

	// failures
	Failures []map[string]string `short:"f" long:"fail" description:"Monitor the file for the supplied regex, and exit with failure if found.  Use the following format: <file to monitor>:<regex to look for>"`
}

// this does the heavy lifting
func parseFile(filename string, searches []search) {
	// print out some information
	fmt.Printf("Monitoring: %v\n", filename)
	fmt.Printf("For the following Strings:\n")
	for _, search := range searches {
		fmt.Printf("\t%v: /%v/\n", search.Status, search.String.String())
	}

	// follow our file
	tail, _ := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true, MustExist: false, Logger: tail.DiscardingLogger})
	for line := range tail.Lines {
		for _, search := range searches {
			if search.String.Match([]byte(line.Text)) == true {
				fmt.Printf("%v matches %v\n", line.Text, search.String.String())
				if search.Status == "failure" {
					os.Exit(1)
				} else {
					os.Exit(0)
				}
			}
		}
	}
}

/* --------------- End Process Command Line ---- */
var wg sync.WaitGroup

func main() {
	// parse our command line flags
	parser := flags.NewParser(&opts, flags.Default)
	description := "This command will monitor files for specific regular expressions, and, depending on the options provided, will exit with either a successful return code, or a failure return code when those regular expressions are found.  The file does not need to exist prior to running this command.\n"

	_, err := parser.Parse()
	ourErr, _ := err.(*flags.Error)

	if ourErr != nil {
		if ourErr.Type == flags.ErrHelp {
			fmt.Printf(description)
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// see if any options were supplied
	if len(opts.Successes) == 0 && len(opts.Failures) == 0 {
		parser.WriteHelp(os.Stdout)
		fmt.Printf("\n%v", description)
	}

	// put together our files and search strings
	files := make(map[string][]search)
	// first process our successes
	// we compile the regular expressions before we put them in
	for _, element := range opts.Successes {
		for filename, value := range element {
			files[filename] = append(files[filename], search{*regexp.MustCompile(value), "success"})
		}
	}
	// do the same for our failures
	for _, element := range opts.Failures {
		for filename, value := range element {
			files[filename] = append(files[filename], search{*regexp.MustCompile(value), "failure"})
		}
	}

	// we need these for goroutines to work
	// make sure we can spin up enough goroutines
	runtime.GOMAXPROCS(len(files) + 1)

	// parse our files and look for the patterns listed
	for file, searches := range files {
		wg.Add(1)
		go parseFile(file, searches)
	}
	wg.Wait()

}
