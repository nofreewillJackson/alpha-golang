// api/handlers.go
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nofreewilljackson/alpha-golang/common"
)

func getDigests(c *gin.Context) {
	rows, err := dbpool.Query(context.Background(),
		`SELECT id, digest, created_at FROM digests ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var digests []common.Digest
	for rows.Next() {
		var digest common.Digest
		err := rows.Scan(&digest.ID, &digest.Digest, &digest.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		digests = append(digests, digest)
	}

	c.JSON(http.StatusOK, gin.H{"digests": digests})
}
