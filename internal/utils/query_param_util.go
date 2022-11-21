package utils

func QueryParamSort(value string) string {
	if value == "asc" {
		return "ASC"
	}
	return "DESC"
}
