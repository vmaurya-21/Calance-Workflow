package repository

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// GetTags returns all tags for a specific repository
// GET /api/repositories/:owner/:repo/tags
func (h *Handler) GetTags(c *gin.Context) {
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

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	tags, err := h.repositoryService.GetTags(c.Request.Context(), accessToken, owner, repo)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to fetch repository tags", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Successfully retrieved %d tags", len(tags)), gin.H{
		"owner":      owner,
		"repository": repo,
		"tags":       tags,
		"count":      len(tags),
	})
}

// CreateTag creates and pushes a tag for a specific commit
// POST /api/repositories/tags
func (h *Handler) CreateTag(c *gin.Context) {
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

	var tagRequest struct {
		Owner     string `json:"owner"`
		Repo      string `json:"repo" binding:"required"`
		TagName   string `json:"tag_name" binding:"required"`
		CommitSHA string `json:"commit_sha" binding:"required"`
	}

	if err := c.ShouldBindJSON(&tagRequest); err != nil {
		pkghttp.BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	if len(tagRequest.CommitSHA) != 40 {
		pkghttp.BadRequestResponse(c, "Invalid commit SHA format. Must be 40 character hash")
		return
	}

	accessToken, err := h.getAccessToken(userUUID)
	if err != nil {
		pkghttp.UnauthorizedResponse(c, "Access token not found. Please login again.")
		return
	}

	owner := tagRequest.Owner
	if owner == "" {
		owner = "default-owner" // Should get from user
	}

	reference, err := h.repositoryService.CreateTag(c.Request.Context(), accessToken, owner, tagRequest.Repo, tagRequest.TagName, tagRequest.CommitSHA)
	if err != nil {
		pkghttp.InternalServerErrorResponse(c, "Failed to create tag", err)
		return
	}

	pkghttp.SuccessResponse(c, http.StatusCreated, "Tag created and pushed successfully", gin.H{
		"owner":      owner,
		"repository": tagRequest.Repo,
		"tag_name":   tagRequest.TagName,
		"commit_sha": tagRequest.CommitSHA,
		"ref":        reference.Ref,
		"object_sha": reference.ObjectSHA,
		"url":        reference.URL,
	})
}
