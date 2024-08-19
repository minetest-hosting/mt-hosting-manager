package db

import (
	"database/sql"
	"mt-hosting-manager/types"
	"strings"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type UserNodeRepository struct {
	dbu *dbutil.DBUtil[*types.UserNode]
}

func (r *UserNodeRepository) Insert(n *types.UserNode) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *UserNodeRepository) Update(n *types.UserNode) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *UserNodeRepository) GetByID(id string) (*types.UserNode, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *UserNodeRepository) GetAll() ([]*types.UserNode, error) {
	return r.dbu.SelectMulti("")
}

func (r *UserNodeRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}

func (r *UserNodeRepository) Search(search *types.UserNodeSearch) ([]*types.UserNode, error) {
	q := strings.Builder{}
	params := []any{}
	q.WriteString("where true")

	if search.ID != nil {
		q.WriteString(" and id = %s")
		params = append(params, *search.ID)
	}

	if search.Name != nil {
		q.WriteString(" and name = %s")
		params = append(params, *search.Name)
	}

	if search.UserID != nil {
		q.WriteString(" and user_id = %s")
		params = append(params, *search.UserID)
	}

	if search.State != nil {
		q.WriteString(" and state = %s")
		params = append(params, *search.State)
	}

	if search.ValidUntil != nil {
		q.WriteString(" and valid_until < %s")
		params = append(params, *search.ValidUntil)
	}

	return r.dbu.SelectMulti(q.String(), params)
}
