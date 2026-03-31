package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/med-000/notifyclass/pkg/service"
)

func main() {
	// --- lock ---
	lockFile := "/tmp/notifyclass.lock"

	if _, err := os.Stat(lockFile); err == nil {
		log.Println("already running, skip")
		return
	}

	err := os.WriteFile(lockFile, []byte("lock"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(lockFile)

	// --- env ---
	_ = godotenv.Load()

	// =========================
	// ① 全授業取得
	// =========================
	req := service.GetCourseRequest{
		UserID:   os.Getenv("USER_ID"),
		Password: os.Getenv("PASSWORD"),
		Year:     2025,
		Term:     1,
	}

	courses, err := service.FetchCourses(req)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := json.MarshalIndent(courses, "", "  ")

	filename := "courses_" + time.Now().Format("20060102_150405") + ".json"
	_ = os.WriteFile(filename, data, 0644)

	log.Println("saved courses:", filename)

	// =========================
	// ② 特定授業取得（ここ追加）
	// =========================
	classReq := service.GetClassRequest{
		UserID:   os.Getenv("USER_ID"),
		Password: os.Getenv("PASSWORD"),
		Year:     2025,
		Term:     1,
		Day:      2, // ← 月曜
		Period:   1, // ← 2限
	}

	class, err := service.FetchClassByRequest(classReq)
	if err != nil {
		log.Println("class fetch error:", err)
		return
	}

	if class == nil {
		log.Println("class not found")
		return
	}

	classData, _ := json.MarshalIndent(class, "", "  ")

	classFile := "class_" + time.Now().Format("20060102_150405") + ".json"
	_ = os.WriteFile(classFile, classData, 0644)

	log.Println("saved class:", classFile)
}
