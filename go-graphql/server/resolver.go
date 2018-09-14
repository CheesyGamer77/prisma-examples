//go:generate gorunpkg github.com/99designs/gqlgen

package main

import (
	"context"

	"github.com/prisma/prisma-examples/go-graphql/prisma-client"
)

type Resolver struct {
	Prisma *prisma.Client
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Post() PostResolver {
	return &postResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateDraft(ctx context.Context, title string, content string, authorEmail string) (prisma.Post, error) {
	return r.Prisma.CreatePost(&prisma.PostCreateInput{
		Title:   &title,
		Content: &content,
		Author: &prisma.UserCreateOneWithoutPostsInput{
			Connect: &prisma.UserWhereUniqueInput{Email: &authorEmail},
		},
	}).Exec()
}
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*prisma.Post, error) {
	// TODO remove pointer once fixed: https://github.com/prisma/prisma/issues/3066
	post, err := r.Prisma.DeletePost(&prisma.PostWhereUniqueInput{ID: &id}).Exec()
	return &post, err
}
func (r *mutationResolver) Publish(ctx context.Context, id string) (*prisma.Post, error) {
	// TODO remove pointer once fixed: https://github.com/prisma/prisma/issues/3066
	isPublished := true
	post, err := r.Prisma.UpdatePost(&prisma.UpdatePostParams{
		Where: &prisma.PostWhereUniqueInput{ID: &id},
		Data:  &prisma.PostUpdateInput{IsPublished: &isPublished},
	}).Exec()
	return &post, err
}

type postResolver struct{ *Resolver }

func (r *postResolver) Author(ctx context.Context, obj *prisma.Post) (prisma.User, error) {
	return r.Prisma.Post(&prisma.PostWhereUniqueInput{ID: &obj.ID}).Author().Exec()
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Feed(ctx context.Context) ([]prisma.Post, error) {
	// TODO remove pointer once fixed: https://github.com/prisma/prisma/issues/3066
	isPublished := true
	return r.Prisma.Posts(&prisma.PostsParams{
		Where: &prisma.PostWhereInput{IsPublished: &isPublished},
	}).Exec()
}
func (r *queryResolver) Drafts(ctx context.Context) ([]prisma.Post, error) {
	// TODO remove pointer once fixed: https://github.com/prisma/prisma/issues/3066
	isPublished := false
	return r.Prisma.Posts(&prisma.PostsParams{
		Where: &prisma.PostWhereInput{IsPublished: &isPublished},
	}).Exec()
}
func (r *queryResolver) Post(ctx context.Context, id string) (*prisma.Post, error) {
	// TODO remove pointer once fixed: https://github.com/prisma/prisma/issues/3066
	post, err := r.Prisma.Post(&prisma.PostWhereUniqueInput{ID: &id}).Exec()
	return &post, err
}

type userResolver struct{ *Resolver }

func (r *userResolver) Posts(ctx context.Context, obj *prisma.User) ([]prisma.Post, error) {
	return r.Prisma.User(&prisma.UserWhereUniqueInput{ID: &obj.ID}).Posts(nil).Exec()
}
