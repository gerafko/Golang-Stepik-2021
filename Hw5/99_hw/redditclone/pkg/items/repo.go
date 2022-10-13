package Posts

import (
	"errors"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type PostsRepo struct {
	//lastID uint32
	data []*Post
}

func NewRepo() *PostsRepo {
	return &PostsRepo{
		data: make([]*Post, 0),
	}
}

func (repo *PostsRepo) GetAll() ([]*Post, error) {
	return repo.data, nil
}

func (repo *PostsRepo) GetSpec(spec string) ([]*Post, error) {
	var searchSpec []*Post
	for _, post := range repo.data {
		if post.Category == spec {
			searchSpec = append(searchSpec, post)
		}
	}
	return searchSpec, nil
}

func (repo *PostsRepo) GetByPostId(id string) (*Post, error) {
	for _, post := range repo.data {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, nil
}

func (repo *PostsRepo) GetByLogin(user_login string) ([]*Post, error) {
	searchSpec := NewRepo()
	for _, post := range repo.data {
		if post.Author.Username == user_login {
			searchSpec.data = append(searchSpec.data, post)
		}
	}
	return searchSpec.data, nil
}

func (repo *PostsRepo) Rate(post Post, post_id string) (*Post, error) {
	for i, post := range repo.data {
		if post.ID == post_id {
			repo.data[i] = post
			return post, nil
		}
	}
	return nil, errors.New("POST 404")
}

func (repo *PostsRepo) RateDown(post_id string) (*Post, error) {
	for _, post := range repo.data {
		if post.ID == post_id {
			post.Votes = append(post.Votes, Votes{"Anton", post.Score - 1})
			post.Score--
			return post, nil
		}
	}
	return nil, errors.New("POST_ID 404")
}

func (repo *PostsRepo) AddPost(p Post) (*Post, error) {
	p.ID = StringWithCharset(15, charset)
	repo.data = append(repo.data, &p)
	return &p, nil
}

func (repo *PostsRepo) DelPost(post_id string) (Message, error) {
	for i, _ := range repo.data {
		if repo.data[i].ID == post_id {
			repo.data = removePost(repo.data, i)
			return Message{"success"}, nil
		}
	}
	return Message{"fail"}, nil
}

func (repo *PostsRepo) AddComm(ic IncomingComment, post_id, u_id, u_name string) (*Post, error) {
	var p Post
	var c Comment
	c.ID = StringWithCharset(15, charset)
	c.Author.ID = u_id
	c.Author.Username = u_name
	c.Created = time.Now()
	c.Body = ic.Body
	for i, post := range repo.data {
		if post.ID == post_id {
			repo.data[i].Comments = append(repo.data[i].Comments, c)
			p = *repo.data[i]
		}
	}
	return &p, nil
}

func (repo *PostsRepo) DelComm(post_id string, comm_id string) (*Post, error) { //есть вариант все данные взять из URL и с помощью func remove ниже
	var p Post
	for i, post := range repo.data {
		if post.ID == post_id {
			for j, _ := range repo.data[i].Comments {
				if repo.data[i].Comments[j].ID == comm_id {
					repo.data[i].Comments = removeComment(repo.data[i].Comments, j)
				}
			}
			p = *repo.data[i]
		}
	}
	return &p, nil
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func removeComment(slice []Comment, s int) []Comment {
	return append(slice[:s], slice[s+1:]...)
}

func removePost(slice []*Post, s int) []*Post {
	return append(slice[:s], slice[s+1:]...)
}
