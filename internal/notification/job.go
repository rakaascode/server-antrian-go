package notification

import "encoding/json"

// JobType tipe job notifikasi
type JobType string

const (
	JobTypeWA JobType = "whatsapp"
)

// NotificationJob adalah payload yang di-push ke Redis queue
type NotificationJob struct {
	Type    JobType `json:"type"`
	Target  string  `json:"target"`  // nomor HP tujuan
	Message string  `json:"message"` // isi pesan
}

func (j NotificationJob) ToJSON() ([]byte, error) {
	return json.Marshal(j)
}

func ParseJob(data []byte) (*NotificationJob, error) {
	var job NotificationJob
	err := json.Unmarshal(data, &job)
	return &job, err
}

// Process mengeksekusi job sesuai tipenya
func (j *NotificationJob) Process() error {
	switch j.Type {
	case JobTypeWA:
		return SendWA(j.Target, j.Message)
	default:
		return nil
	}
}
