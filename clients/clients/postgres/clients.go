package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgtype" // required for SQL access
	"github.com/mainflux/mainflux/clients/clients"
	"github.com/mainflux/mainflux/clients/groups"
	"github.com/mainflux/mainflux/clients/postgres"
	"github.com/mainflux/mainflux/pkg/errors"
)

var _ clients.ClientRepository = (*clientRepo)(nil)

type clientRepo struct {
	db postgres.Database
}

// NewClientRepo instantiates a PostgreSQL
// implementation of Clients repository.
func NewClientRepo(db postgres.Database) clients.ClientRepository {
	return &clientRepo{
		db: db,
	}
}

func (repo clientRepo) Save(ctx context.Context, c clients.Client) (clients.Client, error) {
	q := `INSERT INTO clients (id, name, tags, owner, identity, secret, metadata, created_at, updated_at, status)
        VALUES (:id, :name, :tags, :owner, :identity, :secret, :metadata, :created_at, :updated_at, :status)
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`
	if c.Owner == "" {
		q = `INSERT INTO clients (id, name, tags, identity, secret, metadata, created_at, updated_at, status)
        VALUES (:id, :name, :tags, :identity, :secret, :metadata, :created_at, :updated_at, :status)
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`
	}
	dbc, err := toDBClient(c)
	if err != nil {
		return clients.Client{}, errors.Wrap(errors.ErrCreateEntity, err)
	}

	row, err := repo.db.NamedQueryContext(ctx, q, dbc)
	if err != nil {
		return clients.Client{}, postgres.HandleError(err, errors.ErrCreateEntity)
	}

	defer row.Close()
	row.Next()
	var rClient dbClient
	if err := row.StructScan(&rClient); err != nil {
		return clients.Client{}, err
	}

	return toClient(rClient)
}

func (repo clientRepo) RetrieveByID(ctx context.Context, id string) (clients.Client, error) {
	q := `SELECT id, name, tags, COALESCE(owner, '') AS owner, identity, secret, metadata, created_at, updated_at, status 
        FROM clients
        WHERE id = $1`

	dbc := dbClient{
		ID: id,
	}

	if err := repo.db.QueryRowxContext(ctx, q, id).StructScan(&dbc); err != nil {
		if err == sql.ErrNoRows {
			return clients.Client{}, errors.Wrap(errors.ErrNotFound, err)

		}
		return clients.Client{}, errors.Wrap(errors.ErrViewEntity, err)
	}

	return toClient(dbc)
}

func (repo clientRepo) RetrieveByIdentity(ctx context.Context, identity string) (clients.Client, error) {
	q := fmt.Sprintf(`SELECT id, name, tags, COALESCE(owner, '') AS owner, identity, secret, metadata, created_at, updated_at, status
        FROM clients
        WHERE identity = $1 AND status = %d`, clients.EnabledStatus)

	dbc := dbClient{
		Identity: identity,
	}

	if err := repo.db.QueryRowxContext(ctx, q, identity).StructScan(&dbc); err != nil {
		if err == sql.ErrNoRows {
			return clients.Client{}, errors.Wrap(errors.ErrNotFound, err)

		}
		return clients.Client{}, errors.Wrap(errors.ErrViewEntity, err)
	}

	return toClient(dbc)
}

func (repo clientRepo) RetrieveAll(ctx context.Context, pm clients.Page) (clients.ClientsPage, error) {
	query, err := pageQuery(pm)
	if err != nil {
		return clients.ClientsPage{}, errors.Wrap(errors.ErrViewEntity, err)
	}

	q := fmt.Sprintf(`SELECT c.id, c.name, c.tags, c.identity, c.metadata, COALESCE(c.owner, '') AS owner, c.status, c.created_at
						FROM clients c %s ORDER BY c.created_at LIMIT :limit OFFSET :offset;`, query)

	dbPage, err := toDBClientsPage(pm)
	if err != nil {
		return clients.ClientsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}
	rows, err := repo.db.NamedQueryContext(ctx, q, dbPage)
	if err != nil {
		return clients.ClientsPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}
	defer rows.Close()

	var items []clients.Client
	for rows.Next() {
		dbc := dbClient{}
		if err := rows.StructScan(&dbc); err != nil {
			return clients.ClientsPage{}, errors.Wrap(errors.ErrViewEntity, err)
		}

		c, err := toClient(dbc)
		if err != nil {
			return clients.ClientsPage{}, err
		}

		items = append(items, c)
	}
	cq := fmt.Sprintf(`SELECT COUNT(*) FROM clients c %s;`, query)

	total, err := postgres.Total(ctx, repo.db, cq, dbPage)
	if err != nil {
		return clients.ClientsPage{}, errors.Wrap(errors.ErrViewEntity, err)
	}

	page := clients.ClientsPage{
		Clients: items,
		Page: clients.Page{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}

	return page, nil
}

func (repo clientRepo) Members(ctx context.Context, groupID string, pm clients.Page) (clients.MembersPage, error) {
	emq, err := pageQuery(pm)
	if err != nil {
		return clients.MembersPage{}, err
	}

	q := fmt.Sprintf(`SELECT c.id, c.name, c.tags, c.metadata, c.identity, c.status, c.created_at FROM clients c
		INNER JOIN policies ON c.id=policies.subject %s AND policies.object = :group_id
		AND EXISTS (SELECT 1 FROM policies WHERE policies.subject = '%s' AND '%s'=ANY(actions))
	  	ORDER BY c.created_at LIMIT :limit OFFSET :offset;`, emq, pm.Subject, pm.Action)
	dbPage, err := toDBClientsPage(pm)
	if err != nil {
		return clients.MembersPage{}, errors.Wrap(postgres.ErrFailedToRetrieveAll, err)
	}
	dbPage.GroupID = groupID
	rows, err := repo.db.NamedQueryContext(ctx, q, dbPage)
	if err != nil {
		return clients.MembersPage{}, errors.Wrap(postgres.ErrFailedToRetrieveMembers, err)
	}
	defer rows.Close()

	var items []clients.Client
	for rows.Next() {
		dbc := dbClient{}
		if err := rows.StructScan(&dbc); err != nil {
			return clients.MembersPage{}, errors.Wrap(postgres.ErrFailedToRetrieveMembers, err)
		}

		c, err := toClient(dbc)
		if err != nil {
			return clients.MembersPage{}, err
		}

		items = append(items, c)
	}
	cq := fmt.Sprintf(`SELECT COUNT(*) FROM clients c INNER JOIN policies ON c.id=policies.subject %s AND policies.object = :group_id;`, emq)

	total, err := postgres.Total(ctx, repo.db, cq, dbPage)
	if err != nil {
		return clients.MembersPage{}, errors.Wrap(postgres.ErrFailedToRetrieveMembers, err)
	}

	page := clients.MembersPage{
		Members: items,
		Page: clients.Page{
			Total:  total,
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}
	return page, nil
}

func (repo clientRepo) Update(ctx context.Context, client clients.Client) (clients.Client, error) {
	var query []string
	var upq string
	if client.Name != "" {
		query = append(query, "name = :name,")
	}
	if client.Metadata != nil {
		query = append(query, "metadata = :metadata,")
	}
	if len(query) > 0 {
		upq = strings.Join(query, " ")
	}
	q := fmt.Sprintf(`UPDATE clients SET %s updated_at = :updated_at
        WHERE id = :id AND status = %d
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`,
		upq, clients.EnabledStatus)

	dbu, err := toDBClient(client)
	if err != nil {
		return clients.Client{}, errors.Wrap(errors.ErrUpdateEntity, err)
	}

	row, err := repo.db.NamedQueryContext(ctx, q, dbu)
	if err != nil {
		return clients.Client{}, postgres.HandleError(err, errors.ErrCreateEntity)
	}

	defer row.Close()
	if ok := row.Next(); !ok {
		return clients.Client{}, errors.Wrap(errors.ErrNotFound, row.Err())
	}
	var rClient dbClient
	if err := row.StructScan(&rClient); err != nil {
		return clients.Client{}, err
	}

	return toClient(rClient)
}

func (repo clientRepo) UpdateTags(ctx context.Context, client clients.Client) (clients.Client, error) {
	q := fmt.Sprintf(`UPDATE clients SET tags = :tags, updated_at = :updated_at
        WHERE id = :id AND status = %d
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`,
		clients.EnabledStatus)

	dbu, err := toDBClient(client)
	if err != nil {
		return clients.Client{}, errors.Wrap(errors.ErrUpdateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbu)
	if err != nil {
		return clients.Client{}, postgres.HandleError(err, errors.ErrUpdateEntity)
	}

	defer row.Close()
	if ok := row.Next(); !ok {
		return clients.Client{}, errors.Wrap(errors.ErrNotFound, row.Err())
	}
	var rClient dbClient
	if err := row.StructScan(&rClient); err != nil {
		return clients.Client{}, err
	}

	return toClient(rClient)
}

func (repo clientRepo) UpdateIdentity(ctx context.Context, client clients.Client) (clients.Client, error) {
	q := fmt.Sprintf(`UPDATE clients SET identity = :identity, updated_at = :updated_at
        WHERE id = :id AND status = %d
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`,
		clients.EnabledStatus)

	dbc, err := toDBClient(client)
	if err != nil {
		return clients.Client{}, errors.Wrap(errors.ErrUpdateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbc)
	if err != nil {
		return clients.Client{}, postgres.HandleError(err, errors.ErrUpdateEntity)
	}

	defer row.Close()
	if ok := row.Next(); !ok {
		return clients.Client{}, errors.Wrap(errors.ErrNotFound, row.Err())
	}
	var rClient dbClient
	if err := row.StructScan(&rClient); err != nil {
		return clients.Client{}, err
	}

	return toClient(rClient)
}

func (repo clientRepo) UpdateSecret(ctx context.Context, client clients.Client) (clients.Client, error) {
	q := fmt.Sprintf(`UPDATE clients SET secret = :secret, updated_at = :updated_at
        WHERE identity = :identity AND status = %d
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`,
		clients.EnabledStatus)

	dbc, err := toDBClient(client)
	if err != nil {
		return clients.Client{}, errors.Wrap(errors.ErrUpdateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbc)
	if err != nil {
		return clients.Client{}, postgres.HandleError(err, errors.ErrUpdateEntity)
	}

	defer row.Close()
	if ok := row.Next(); !ok {
		return clients.Client{}, errors.Wrap(errors.ErrNotFound, row.Err())
	}
	var rClient dbClient
	if err := row.StructScan(&rClient); err != nil {
		return clients.Client{}, err
	}

	return toClient(rClient)
}

func (repo clientRepo) UpdateOwner(ctx context.Context, client clients.Client) (clients.Client, error) {
	q := fmt.Sprintf(`UPDATE clients SET owner = :owner, updated_at = :updated_at
        WHERE id = :id AND status = %d
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`,
		clients.EnabledStatus)

	dbc, err := toDBClient(client)
	if err != nil {
		return clients.Client{}, errors.Wrap(errors.ErrUpdateEntity, err)
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbc)
	if err != nil {
		return clients.Client{}, postgres.HandleError(err, errors.ErrUpdateEntity)
	}

	defer row.Close()
	if ok := row.Next(); !ok {
		return clients.Client{}, errors.Wrap(errors.ErrNotFound, row.Err())
	}
	var rClient dbClient
	if err := row.StructScan(&rClient); err != nil {
		return clients.Client{}, err
	}

	return toClient(rClient)
}

func (repo clientRepo) ChangeStatus(ctx context.Context, id string, status clients.Status) (clients.Client, error) {
	q := fmt.Sprintf(`UPDATE clients SET status = %d WHERE id = :id
        RETURNING id, name, tags, identity, metadata, COALESCE(owner, '') AS owner, status, created_at, updated_at`, status)

	dbc := dbClient{
		ID: id,
	}
	row, err := repo.db.NamedQueryContext(ctx, q, dbc)
	if err != nil {
		return clients.Client{}, postgres.HandleError(err, errors.ErrUpdateEntity)
	}

	defer row.Close()
	if ok := row.Next(); !ok {
		return clients.Client{}, errors.Wrap(errors.ErrNotFound, row.Err())
	}
	var rClient dbClient
	if err := row.StructScan(&rClient); err != nil {
		return clients.Client{}, errors.Wrap(errors.ErrUpdateEntity, err)
	}

	return toClient(rClient)
}

type dbClient struct {
	ID        string           `db:"id"`
	Name      string           `db:"name,omitempty"`
	Tags      pgtype.TextArray `db:"tags,omitempty"`
	Identity  string           `db:"identity"`
	Owner     string           `db:"owner,omitempty"` // nullable
	Secret    string           `db:"secret"`
	Metadata  []byte           `db:"metadata,omitempty"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
	Groups    []groups.Group   `db:"groups"`
	Status    clients.Status   `db:"status"`
}

func toDBClient(c clients.Client) (dbClient, error) {
	data := []byte("{}")
	if len(c.Metadata) > 0 {
		b, err := json.Marshal(c.Metadata)
		if err != nil {
			return dbClient{}, errors.Wrap(errors.ErrMalformedEntity, err)
		}
		data = b
	}
	var tags pgtype.TextArray
	if err := tags.Set(c.Tags); err != nil {
		return dbClient{}, err
	}

	return dbClient{
		ID:        c.ID,
		Name:      c.Name,
		Tags:      tags,
		Owner:     c.Owner,
		Identity:  c.Credentials.Identity,
		Secret:    c.Credentials.Secret,
		Metadata:  data,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Status:    c.Status,
	}, nil
}

func toClient(c dbClient) (clients.Client, error) {
	var metadata clients.Metadata
	if c.Metadata != nil {
		if err := json.Unmarshal([]byte(c.Metadata), &metadata); err != nil {
			return clients.Client{}, errors.Wrap(errors.ErrMalformedEntity, err)
		}
	}
	var tags []string
	for _, e := range c.Tags.Elements {
		tags = append(tags, e.String)
	}

	return clients.Client{
		ID:    c.ID,
		Name:  c.Name,
		Tags:  tags,
		Owner: c.Owner,
		Credentials: clients.Credentials{
			Identity: c.Identity,
			Secret:   c.Secret,
		},
		Metadata:  metadata,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Status:    c.Status,
	}, nil
}

func pageQuery(pm clients.Page) (string, error) {
	mq, _, err := postgres.CreateMetadataQuery("", pm.Metadata)
	if err != nil {
		return "", errors.Wrap(errors.ErrViewEntity, err)
	}
	var query []string
	var emq string
	if mq != "" {
		query = append(query, mq)
	}
	if pm.Name != "" {
		query = append(query, fmt.Sprintf("c.name = '%s'", pm.Name))
	}
	if pm.Tag != "" {
		query = append(query, fmt.Sprintf("'%s' = ANY(c.tags)", pm.Tag))
	}
	if pm.Status != clients.AllStatus {
		query = append(query, fmt.Sprintf("c.status = %d", pm.Status))
	}
	// For listing clients that the specified client owns but not sharedby
	if pm.OwnerID != "" && pm.SharedBy == "" {
		query = append(query, fmt.Sprintf("c.owner = '%s'", pm.OwnerID))
	}

	// For listing clients that the specified client owns and that are shared with the specified client
	if pm.OwnerID != "" && pm.SharedBy != "" {
		query = append(query, fmt.Sprintf("(c.owner = '%s' OR policies.object IN (SELECT object FROM policies WHERE subject = '%s' AND '%s'=ANY(actions)))", pm.OwnerID, pm.SharedBy, pm.Action))
	}
	// For listing clients that the specified client is shared with
	if pm.SharedBy != "" && pm.OwnerID == "" {
		query = append(query, fmt.Sprintf("c.owner != '%s' AND (policies.object IN (SELECT object FROM policies WHERE subject = '%s' AND '%s'=ANY(actions)))", pm.SharedBy, pm.SharedBy, pm.Action))
	}
	if len(query) > 0 {
		emq = fmt.Sprintf("WHERE %s", strings.Join(query, " AND "))
		if strings.Contains(emq, "policies") {
			emq = fmt.Sprintf("JOIN policies ON policies.subject = c.id %s", emq)
		}
	}
	return emq, nil

}

func toDBClientsPage(pm clients.Page) (dbClientsPage, error) {
	_, data, err := postgres.CreateMetadataQuery("", pm.Metadata)
	if err != nil {
		return dbClientsPage{}, errors.Wrap(errors.ErrViewEntity, err)
	}
	return dbClientsPage{
		Name:     pm.Name,
		Metadata: data,
		Owner:    pm.OwnerID,
		Total:    pm.Total,
		Offset:   pm.Offset,
		Limit:    pm.Limit,
		Status:   pm.Status,
		Tag:      pm.Tag,
	}, nil
}

type dbClientsPage struct {
	GroupID  string         `db:"group_id"`
	Name     string         `db:"name"`
	Owner    string         `db:"owner"`
	Identity string         `db:"identity"`
	Metadata []byte         `db:"metadata"`
	Tag      string         `db:"tag"`
	Status   clients.Status `db:"status"`
	Total    uint64         `db:"total"`
	Limit    uint64         `db:"limit"`
	Offset   uint64         `db:"offset"`
}