package db

import (
	"time"

	"github.com/gocql/gocql"
)

func InsertReport(reportID, imageID, userID, reason string, createdAt time.Time) error {
	imgUUID, err := gocql.ParseUUID(imageID)
	if err != nil {
		return err
	}
	repUUID, err := gocql.ParseUUID(reportID)
	if err != nil {
		return err
	}
	reporterUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return err
	}
	return Session.Query(`
		INSERT INTO reports_by_image (image_id, report_id, reporter_id, reason, reported_at)
		VALUES (?, ?, ?, ?, ?)`,
		imgUUID, repUUID, reporterUUID, reason, createdAt,
	).Exec()
}

func IncrementImageReportCounter(imageID string) error {
	return Session.Query(`
		UPDATE image_counters SET reports = reports + 1 WHERE image_id = ?`,
		imageID,
	).Exec()
}

func GetReportsByImage(imageID string) ([]map[string]interface{}, error) {
	iter := Session.Query(`
		SELECT report_id, reporter_id, reason, reported_at FROM reports_by_image WHERE image_id = ?`,
		imageID,
	).Iter()
	var reports []map[string]interface{}
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		reports = append(reports, m)
		m = map[string]interface{}{}
	}
	return reports, iter.Close()
}

func GetImageReportCount(imageID string) (int, error) {
	var count int
	err := Session.Query(`
		SELECT reports FROM image_counters WHERE image_id = ?`,
		imageID,
	).Scan(&count)
	return count, err
}
