<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# PBS Job Operator - Kubernetes Operator for PBS Jobs

This project is a Kubernetes operator written in Go that manages PBS (Portable Batch System) jobs through custom resources.

## Project Structure
- **Custom Resource Definition (CRD)**: Defines PBSJob custom resource
- **Controller**: Watches PBSJob resources and manages their lifecycle
- **PBS Integration**: Uses qsub to submit jobs and qstat to monitor status
- **Reconciliation Loop**: Periodically updates job status and handles deletions

## Development Guidelines
- Use controller-runtime framework for operator development
- Implement proper error handling and retry logic for PBS commands
- Follow Kubernetes operator best practices
- Use structured logging with logr
- Implement proper cleanup when PBSJob resources are deleted
- Handle PBS job states (Queued, Running, Completed, Error)

## PBS Integration
- Use `qsub` command to submit jobs to PBS
- Use `qstat` command to query job status
- Use `qdel` command to delete PBS jobs
- Parse PBS command outputs properly
- Handle PBS command failures gracefully

## Testing
- Unit tests for controllers and PBS integration
- Integration tests with mock PBS commands
- End-to-end tests with actual PBS system (if available)