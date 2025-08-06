package main

import (
	"log"
	"test/internal/adapter/handler"
	"test/internal/core/service"

	"github.com/gin-gonic/gin"

	
)

func main() {
	svc, err := service.NewHelmService()

	if err != nil {
		log.Fatalf("failed to initialize Helm service: %v", err)
	}
	h := handler.NewHttpHandler(svc)

	r := gin.Default()
	r.POST("/api/deploy/:namespace/:release", h.Deploy)

	r.POST("/api/update/:namespace/:release", h.Update)

	r.DELETE("/api/delete/:namespace/:release", h.Delete)

	//1. Xem danh sách release trong namespace
	r.GET("/api/releases/:namespace/:release", h.GetRelease)

	// //2. Xem chi tiết 1 release
	// r.GET("/api/release/:namespace/:release", GetReleaseDetail)

	// //3. Xem giá trị values.yaml đã được áp dụng
	// r.GET("/api/release/:namespace/:release/values", GetReleaseValues)

	// //4. Xem chart metadata
	// r.GET("/api/release/:namespace/:release/chart", GetChartInfo)

	// //5. Xem image đang dùng trong release
	// r.GET("/api/release/:namespace/:release/images", GetReleaseImages)

	// //6. Xem manifest (YAML Kubernetes đã được render)
	// r.GET("/api/release/:namespace/:release/manifest", GetReleaseManifest)

	// //7. Xem logs của các Pod thuộc release (nếu có)
	// r.GET("/api/release/:namespace/:release/logs", GetReleaseLogs)

	// //8. Rollback release
	// r.POST("/api/release/:namespace/:release/rollback", RollbackRelease)

	// //9. Xem history release
	// r.GET("/api/release/:namespace/:release/history", GetReleaseHistory)

	// //10. Diff chart (giữa version mới và đang deploy)
	// r.POST("/api/release/:namespace/:release/diff", DiffRelease)

	r.Run(":8080")
}
