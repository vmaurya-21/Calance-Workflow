package repository

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetWorkflowRuns returns workflow runs for a specific repository
// GET /api/repositories/:owner/:repo/actions/runs
func (h *Handler) GetWorkflowRuns(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		pkghttp.BadRequestResponse(c, "Owner and repository name are required")
		return
	}

	perPage := 30
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		var parsed int
		if _, err := fmt.Sscanf(perPageStr, "%d", &parsed); err == nil {
			if parsed > 0 && parsed <= 100 {
				perPage = parsed
			}
		}
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	runs, err := h.repositoryService.GetWorkflowRuns(c.Request.Context(), accessToken, owner, repo, perPage)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch workflow runs", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Workflow runs fetched successfully", gin.H{
		"owner":      owner,
		"repository": repo,
		"runs":       runs,
		"run_count":  len(runs),
	})
}

// GetWorkflowRunDetail returns detailed information about a specific workflow run
// GET /api/repositories/:owner/:repo/actions/runs/:run_id
func (h *Handler) GetWorkflowRunDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	owner := c.Param("owner")
	repo := c.Param("repo")
	runIDStr := c.Param("run_id")

	if owner == "" || repo == "" || runIDStr == "" {
		pkghttp.BadRequestResponse(c, "Owner, repository name, and run ID are required")
		return
	}

	var runID int64
	if _, err := fmt.Sscanf(runIDStr, "%d", &runID); err != nil {
		pkghttp.BadRequestResponse(c, "Invalid run ID format")
		return
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	runDetail, jobs, err := h.repositoryService.GetWorkflowRunDetail(c.Request.Context(), accessToken, owner, repo, runID)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch workflow run detail", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Workflow run detail fetched successfully", gin.H{
		"owner":      owner,
		"repository": repo,
		"run":        runDetail,
		"jobs":       jobs,
		"job_count":  len(jobs),
	})
}

// GetJobLogs returns ALL logs for a job in a single response
// GET /api/repositories/:owner/:repo/actions/jobs/:job_id/logs
func (h *Handler) GetJobLogs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkghttp.UnauthorizedResponse(c, "User not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Invalid user ID format", err)
		return
	}

	owner := c.Param("owner")
	repo := c.Param("repo")
	jobIDStr := c.Param("job_id")

	if owner == "" || repo == "" || jobIDStr == "" {
		pkghttp.BadRequestResponse(c, "Owner, repository, and job ID are required")
		return
	}

	var jobID int64
	if _, err := fmt.Sscanf(jobIDStr, "%d", &jobID); err != nil {
		pkghttp.BadRequestResponse(c, "Invalid job ID format")
		return
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	logs, err := h.repositoryService.GetJobLogs(c.Request.Context(), accessToken, owner, repo, jobID)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch job logs", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, "Job logs fetched successfully", gin.H{
		"owner":      owner,
		"repository": repo,
		"job_id":     jobID,
		"logs":       logs,
		"size_bytes": len(logs),
	})
}
