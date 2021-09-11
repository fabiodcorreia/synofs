package synofs

import "fmt"

type ClientError struct {
	API  string
	Code int `json:"code"`
}

func (e ClientError) Error() string {
	var description string

	switch e.API {
	case authAPI:
		description = authError(e.Code)
	}

	if description != "" {
		return description
	}

	description = commonError(e.Code)

	if description != "" {
		return description
	}

	return fmt.Sprintf("unexpected error with code %d for api %q", e.Code, e.API)
}

func commonError(code int) string {
	switch code {
	case 100:
		return "unknown error"
	case 101:
		return "no parameter of api, method or version"
	case 102:
		return "the requested API does not exist"
	case 103:
		return "the requested method does not exist"
	case 104:
		return "the requested version does not support the functionality"
	case 105:
		return "the logged in session does not have permission"
	case 106:
		return "session timeout"
	case 107:
		return "session interrupted by duplicate login"
	case 109:
		return "sid not found"
	case 400:
		return "invalid parameter of file operation"
	case 401:
		return "unknown error of file operation"
	case 402:
		return "system is too busy"
	case 403:
		return "invalid user does this file operation"
	case 404:
		return "invalid group does this file operation"
	case 405:
		return "invalid user and group does this file operation"
	case 406:
		return "can't get user/group information from the account server"
	case 407:
		return "operation not permitted"
	case 408:
		return "no such file or directory"
	case 409:
		return "non-supported file system"
	case 410:
		return "failed to connect internet-based file system (e.g., CIFS)"
	case 411:
		return "read-only file system"
	case 412:
		return "filename too long in the non-encrypted file system"
	case 413:
		return "filename too long in the encrypted file system"
	case 414:
		return "file already exists"
	case 415:
		return "disk quota exceeded"
	case 416:
		return "no space left on device"
	case 417:
		return "input/output error"
	case 418:
		return "illegal name or path"
	case 419:
		return "illegal file name"
	case 420:
		return "illegal file name on FAT file system"
	case 421:
		return "device or resource busy"
	case 599:
		return "no such task of the file operation"
	}
	return ""
}

func authError(code int) string {
	switch code {
	case 400:
		return "username/password incorrect"
	case 401:
		return "account disabled"
	case 402:
		return "permission denied"
	case 403:
		return "2-step verification code required"
	case 404:
		return "failed to authenticate 2-step verification code"
	}
	return ""
}
