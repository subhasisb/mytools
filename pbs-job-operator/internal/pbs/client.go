/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pbs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	batchv1 "github.com/subhas/pbs-job-operator/api/v1"
)

// PBSClient handles interactions with PBS commands
type PBSClient struct {
	logger logr.Logger
}

// NewPBSClient creates a new PBS client
func NewPBSClient(logger logr.Logger) *PBSClient {
	return &PBSClient{
		logger: logger,
	}
}

// JobInfo represents information about a PBS job
type JobInfo struct {
	ID           string
	Name         string
	State        string
	Queue        string
	SubmitTime   time.Time
	StartTime    time.Time
	CompleteTime time.Time
	ExitCode     int32
}

// SubmitJob submits a job to PBS using qsub
func (c *PBSClient) SubmitJob(job *batchv1.PBSJob) (string, error) {
	// Create temporary script file
	scriptFile, err := c.createScriptFile(job)
	if err != nil {
		return "", fmt.Errorf("failed to create script file: %w", err)
	}
	defer os.Remove(scriptFile)

	// Build qsub command
	args := []string{scriptFile}

	if job.Spec.JobName != "" {
		args = append([]string{"-N", job.Spec.JobName}, args...)
	}

	if job.Spec.Queue != "" {
		args = append([]string{"-q", job.Spec.Queue}, args...)
	}

	if job.Spec.Walltime != "" {
		args = append([]string{"-l", "walltime=" + job.Spec.Walltime}, args...)
	}

	if job.Spec.OutputPath != "" {
		args = append([]string{"-o", job.Spec.OutputPath}, args...)
	}

	if job.Spec.ErrorPath != "" {
		args = append([]string{"-e", job.Spec.ErrorPath}, args...)
	}

	// Add resource requirements
	for key, value := range job.Spec.Resources {
		args = append([]string{"-l", key + "=" + value}, args...)
	}

	fmt.Printf("/opt/pbs/bin/qsub -f %s\n", strings.Join(args, " "))

	// Execute qsub command
	cmd := exec.Command("/opt/pbs/bin/qsub", args...)

	// Set working directory if specified
	if job.Spec.WorkingDirectory != "" {
		cmd.Dir = job.Spec.WorkingDirectory
	}

	// Set environment variables
	if len(job.Spec.Environment) > 0 {
		env := os.Environ()
		for key, value := range job.Spec.Environment {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("qsub command failed: %w", err)
	}

	// Extract job ID from qsub output (format: "job_id.server_name")
	jobID := strings.TrimSpace(string(output))
	c.logger.Info("Job submitted successfully", "jobID", jobID)

	return jobID, nil
}

// GetJobInfo queries PBS for job information using qstat
func (c *PBSClient) GetJobInfo(jobID string) (*JobInfo, error) {
	// Use qstat -f for detailed job information
	cmd := exec.Command("/opt/pbs/bin/qstat", "-f", jobID)
	output, err := cmd.Output()
	if err != nil {
		// Job might not exist anymore
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 153 {
			return nil, fmt.Errorf("job not found: %s", jobID)
		}
		return nil, fmt.Errorf("qstat command failed: %w", err)
	}

	return c.parseJobInfo(string(output))
}

// DeleteJob deletes a PBS job using qdel
func (c *PBSClient) DeleteJob(jobID string) error {
	cmd := exec.Command("/opt/pbs/bin/qdel", jobID)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("qdel command failed: %w", err)
	}

	c.logger.Info("Job deleted successfully", "jobID", jobID)
	return nil
}

// createScriptFile creates a temporary script file from the job spec
func (c *PBSClient) createScriptFile(job *batchv1.PBSJob) (string, error) {
	tmpFile, err := os.CreateTemp("", "pbs-job-*.sh")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Write script content
	_, err = tmpFile.WriteString(job.Spec.Script)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	// Make script executable
	err = os.Chmod(tmpFile.Name(), 0755)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// parseJobInfo parses qstat -f output into JobInfo struct
func (c *PBSClient) parseJobInfo(output string) (*JobInfo, error) {
	jobInfo := &JobInfo{}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse job ID from first line
		if strings.Contains(line, "Job Id:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				jobInfo.ID = strings.TrimSpace(parts[1])
			}
			continue
		}

		// Parse key-value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch {
			case strings.HasSuffix(key, "Job_Name"):
				jobInfo.Name = value
			case strings.HasSuffix(key, "job_state"):
				jobInfo.State = value
			case strings.HasSuffix(key, "queue"):
				jobInfo.Queue = value
			case strings.HasSuffix(key, "qtime"):
				if t, err := parseTimeStamp(value); err == nil {
					jobInfo.SubmitTime = t
				}
			case strings.HasSuffix(key, "start_time"):
				if t, err := parseTimeStamp(value); err == nil {
					jobInfo.StartTime = t
				}
			case strings.HasSuffix(key, "comp_time"):
				if t, err := parseTimeStamp(value); err == nil {
					jobInfo.CompleteTime = t
				}
			case strings.HasSuffix(key, "exit_status"):
				if code, err := strconv.Atoi(value); err == nil {
					jobInfo.ExitCode = int32(code)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error parsing qstat output: %w", err)
	}

	return jobInfo, nil
}

// parseTimeStamp parses PBS timestamp format
func parseTimeStamp(timestamp string) (time.Time, error) {
	// PBS typically uses format like "Mon Oct 16 10:30:00 2025"
	layouts := []string{
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan 2 15:04:05 2006",
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timestamp); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestamp)
}

// GetJobState maps PBS job state to our custom phases
func GetJobState(pbsState string) string {
	switch pbsState {
	case "Q":
		return "Queued"
	case "R":
		return "Running"
	case "C":
		return "Completed"
	case "E", "H":
		return "Failed"
	default:
		return "Unknown"
	}
}
