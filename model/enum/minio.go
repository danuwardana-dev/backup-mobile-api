package enum

type MinioPathName string

const (
	MINIO_KYC_OCR              MinioPathName = "KYC_OCR"
	MINIO_KYC_FACE_COMPARE     MinioPathName = "KYC_FACE_COMPARE"
	MINIO_USER_PROFILE_PICTURE MinioPathName = "USER_PROFILE_PICTURE"
)

var MinioPathNameMap = map[MinioPathName]string{
	MINIO_KYC_OCR:              "user/kyc/ocr",
	MINIO_KYC_FACE_COMPARE:     "user/kyc/face-compare",
	MINIO_USER_PROFILE_PICTURE: "user/profile-picture",
}
