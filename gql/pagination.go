package gql

const DefaultPage = 0
const DefaultPageSize = 10

func paginationToSQL(page, pageSize *int) (offset, limit int) {
	if page == nil {
		return 0, DefaultPageSize
	}
	if pageSize == nil {
		return *page * DefaultPageSize, DefaultPageSize
	}
	return *page * *pageSize, *pageSize
}
