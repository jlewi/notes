// main is a simple program to test whether GCP APIs support CORS
package main

import (
	"fmt"
	"github.com/jlewi/monogo/helpers"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

func run() error {
	url := "https://bigquery.googleapis.com/bigquery/v2/projects/dev-sailplane/datasets/traces/tables/AgentTraces?alt=json&prettyPrint=false"
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Origin", "http://localhost:8080")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "authorization")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Failed to make request to %v", url)
	}

	fmt.Printf("Response:\n")
	fmt.Printf("StatusCode: %v\n", resp.StatusCode)
	fmt.Printf("Status: %v\n", resp.Status)
	fmt.Printf("Headers:\n%+v\n", helpers.PrettyString(resp.Header))
	if resp.Body != nil {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrapf(err, "Failed to read response body")
		}

		fmt.Printf("\nBody: \n%v", string(b))
	} else {
		fmt.Print("No body")
	}

	// Is this the right way to verify CORS?
	if resp.StatusCode == http.StatusOK {
		fmt.Print("Request succeeded; CORS is probably supported")
	} else {
		fmt.Print("Request failed")
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Print("Error: %+v", err)
		os.Exit(1)
	}

}
