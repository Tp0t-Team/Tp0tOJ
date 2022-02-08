package user

type MutationResolver struct {
}

func (r *MutationResolver) Hello() string {
	return "hello world!"
}
