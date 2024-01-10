package commoninputs

type PagingParams struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}
