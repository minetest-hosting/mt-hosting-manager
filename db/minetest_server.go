package db

import (
	"database/sql"
	"mt-hosting-manager/types"
	"strings"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type MinetestServerRepository struct {
	dbu *dbutil.DBUtil[*types.MinetestServer]
}

func (r *MinetestServerRepository) Insert(n *types.MinetestServer) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.dbu.Insert(n)
}

func (r *MinetestServerRepository) Update(n *types.MinetestServer) error {
	return r.dbu.Update(n, "where id = %s", n.ID)
}

func (r *MinetestServerRepository) GetByID(id string) (*types.MinetestServer, error) {
	nt, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nt, err
}

func (r *MinetestServerRepository) GetAll() ([]*types.MinetestServer, error) {
	return r.dbu.SelectMulti("")
}

func (r *MinetestServerRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}

func (r *MinetestServerRepository) Search(search *types.MinetestServerSearch) ([]*types.MinetestServer, error) {
	q := strings.Builder{}
	params := []any{}

	q.WriteString("where true")

	if search.ID != nil {
		q.WriteString(" and id = %s")
		params = append(params, *search.ID)
	}

	if search.UserID != nil {
		q.WriteString(" and user_node_id in (select id from user_node where user_id = %s)")
		params = append(params, *search.UserID)
	}

	if search.NodeID != nil {
		q.WriteString(" and user_node_id = %s")
		params = append(params, *search.NodeID)
	}

	if search.State != nil {
		q.WriteString(" and state = %s")
		params = append(params, *search.State)
	}

	return r.dbu.SelectMulti(q.String(), params...)
}
