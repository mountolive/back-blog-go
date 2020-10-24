package usecase

type createPostTestCase struct {
	Name        string
	Description string
	Dto         *CreatePostDto
	ExpErr      bool
}

type updatePostTestCase struct {
	Name        string
	Description string
	Dto         *UpdatePostDto
	ExpErr      bool
}

type filterByTagTestCase struct {
	Name        string
	Description string
	Dto         *ByTagDto
	ExpErr      bool
}

type filterByDateRangeTestCase struct {
	Name        string
	Description string
	Dto         *ByDateRangeDto
	ExpErr      bool
}
