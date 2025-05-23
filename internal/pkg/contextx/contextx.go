package contextx

import "context"

type AuthKey string

const StudentIDKey AuthKey = "student_id"

// SetStudentID 存入 studentId 到 context
func SetStudentID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, StudentIDKey, id)
}

// GetStudentID 从 context 中提取 studentId
func GetStudentID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(StudentIDKey).(string)
	return v, ok
}
