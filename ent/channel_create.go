// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/taythebot/discord_chatbot/ent/channel"
	"github.com/taythebot/discord_chatbot/ent/guild"
	"github.com/taythebot/discord_chatbot/ent/message"
)

// ChannelCreate is the builder for creating a Channel entity.
type ChannelCreate struct {
	config
	mutation *ChannelMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetName sets the "name" field.
func (cc *ChannelCreate) SetName(s string) *ChannelCreate {
	cc.mutation.SetName(s)
	return cc
}

// SetModel sets the "model" field.
func (cc *ChannelCreate) SetModel(s string) *ChannelCreate {
	cc.mutation.SetModel(s)
	return cc
}

// SetPrompt sets the "prompt" field.
func (cc *ChannelCreate) SetPrompt(s string) *ChannelCreate {
	cc.mutation.SetPrompt(s)
	return cc
}

// SetGuildID sets the "guild_id" field.
func (cc *ChannelCreate) SetGuildID(s string) *ChannelCreate {
	cc.mutation.SetGuildID(s)
	return cc
}

// SetCreatedAt sets the "created_at" field.
func (cc *ChannelCreate) SetCreatedAt(t time.Time) *ChannelCreate {
	cc.mutation.SetCreatedAt(t)
	return cc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cc *ChannelCreate) SetNillableCreatedAt(t *time.Time) *ChannelCreate {
	if t != nil {
		cc.SetCreatedAt(*t)
	}
	return cc
}

// SetUpdatedAt sets the "updated_at" field.
func (cc *ChannelCreate) SetUpdatedAt(t time.Time) *ChannelCreate {
	cc.mutation.SetUpdatedAt(t)
	return cc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (cc *ChannelCreate) SetNillableUpdatedAt(t *time.Time) *ChannelCreate {
	if t != nil {
		cc.SetUpdatedAt(*t)
	}
	return cc
}

// SetID sets the "id" field.
func (cc *ChannelCreate) SetID(s string) *ChannelCreate {
	cc.mutation.SetID(s)
	return cc
}

// SetOwnerID sets the "owner" edge to the Guild entity by ID.
func (cc *ChannelCreate) SetOwnerID(id string) *ChannelCreate {
	cc.mutation.SetOwnerID(id)
	return cc
}

// SetOwner sets the "owner" edge to the Guild entity.
func (cc *ChannelCreate) SetOwner(g *Guild) *ChannelCreate {
	return cc.SetOwnerID(g.ID)
}

// AddMessageIDs adds the "messages" edge to the Message entity by IDs.
func (cc *ChannelCreate) AddMessageIDs(ids ...int) *ChannelCreate {
	cc.mutation.AddMessageIDs(ids...)
	return cc
}

// AddMessages adds the "messages" edges to the Message entity.
func (cc *ChannelCreate) AddMessages(m ...*Message) *ChannelCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return cc.AddMessageIDs(ids...)
}

// Mutation returns the ChannelMutation object of the builder.
func (cc *ChannelCreate) Mutation() *ChannelMutation {
	return cc.mutation
}

// Save creates the Channel in the database.
func (cc *ChannelCreate) Save(ctx context.Context) (*Channel, error) {
	cc.defaults()
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *ChannelCreate) SaveX(ctx context.Context) *Channel {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *ChannelCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *ChannelCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *ChannelCreate) defaults() {
	if _, ok := cc.mutation.CreatedAt(); !ok {
		v := channel.DefaultCreatedAt()
		cc.mutation.SetCreatedAt(v)
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		v := channel.DefaultUpdatedAt()
		cc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *ChannelCreate) check() error {
	if _, ok := cc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Channel.name"`)}
	}
	if v, ok := cc.mutation.Name(); ok {
		if err := channel.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Channel.name": %w`, err)}
		}
	}
	if _, ok := cc.mutation.Model(); !ok {
		return &ValidationError{Name: "model", err: errors.New(`ent: missing required field "Channel.model"`)}
	}
	if v, ok := cc.mutation.Model(); ok {
		if err := channel.ModelValidator(v); err != nil {
			return &ValidationError{Name: "model", err: fmt.Errorf(`ent: validator failed for field "Channel.model": %w`, err)}
		}
	}
	if _, ok := cc.mutation.Prompt(); !ok {
		return &ValidationError{Name: "prompt", err: errors.New(`ent: missing required field "Channel.prompt"`)}
	}
	if v, ok := cc.mutation.Prompt(); ok {
		if err := channel.PromptValidator(v); err != nil {
			return &ValidationError{Name: "prompt", err: fmt.Errorf(`ent: validator failed for field "Channel.prompt": %w`, err)}
		}
	}
	if _, ok := cc.mutation.GuildID(); !ok {
		return &ValidationError{Name: "guild_id", err: errors.New(`ent: missing required field "Channel.guild_id"`)}
	}
	if v, ok := cc.mutation.GuildID(); ok {
		if err := channel.GuildIDValidator(v); err != nil {
			return &ValidationError{Name: "guild_id", err: fmt.Errorf(`ent: validator failed for field "Channel.guild_id": %w`, err)}
		}
	}
	if _, ok := cc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Channel.created_at"`)}
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Channel.updated_at"`)}
	}
	if v, ok := cc.mutation.ID(); ok {
		if err := channel.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Channel.id": %w`, err)}
		}
	}
	if len(cc.mutation.OwnerIDs()) == 0 {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Channel.owner"`)}
	}
	return nil
}

func (cc *ChannelCreate) sqlSave(ctx context.Context) (*Channel, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected Channel.ID type: %T", _spec.ID.Value)
		}
	}
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *ChannelCreate) createSpec() (*Channel, *sqlgraph.CreateSpec) {
	var (
		_node = &Channel{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(channel.Table, sqlgraph.NewFieldSpec(channel.FieldID, field.TypeString))
	)
	_spec.OnConflict = cc.conflict
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := cc.mutation.Name(); ok {
		_spec.SetField(channel.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := cc.mutation.Model(); ok {
		_spec.SetField(channel.FieldModel, field.TypeString, value)
		_node.Model = value
	}
	if value, ok := cc.mutation.Prompt(); ok {
		_spec.SetField(channel.FieldPrompt, field.TypeString, value)
		_node.Prompt = value
	}
	if value, ok := cc.mutation.CreatedAt(); ok {
		_spec.SetField(channel.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := cc.mutation.UpdatedAt(); ok {
		_spec.SetField(channel.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := cc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   channel.OwnerTable,
			Columns: []string{channel.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guild.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.GuildID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.MessagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   channel.MessagesTable,
			Columns: []string{channel.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Channel.Create().
//		SetName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ChannelUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (cc *ChannelCreate) OnConflict(opts ...sql.ConflictOption) *ChannelUpsertOne {
	cc.conflict = opts
	return &ChannelUpsertOne{
		create: cc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Channel.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cc *ChannelCreate) OnConflictColumns(columns ...string) *ChannelUpsertOne {
	cc.conflict = append(cc.conflict, sql.ConflictColumns(columns...))
	return &ChannelUpsertOne{
		create: cc,
	}
}

type (
	// ChannelUpsertOne is the builder for "upsert"-ing
	//  one Channel node.
	ChannelUpsertOne struct {
		create *ChannelCreate
	}

	// ChannelUpsert is the "OnConflict" setter.
	ChannelUpsert struct {
		*sql.UpdateSet
	}
)

// SetName sets the "name" field.
func (u *ChannelUpsert) SetName(v string) *ChannelUpsert {
	u.Set(channel.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ChannelUpsert) UpdateName() *ChannelUpsert {
	u.SetExcluded(channel.FieldName)
	return u
}

// SetModel sets the "model" field.
func (u *ChannelUpsert) SetModel(v string) *ChannelUpsert {
	u.Set(channel.FieldModel, v)
	return u
}

// UpdateModel sets the "model" field to the value that was provided on create.
func (u *ChannelUpsert) UpdateModel() *ChannelUpsert {
	u.SetExcluded(channel.FieldModel)
	return u
}

// SetPrompt sets the "prompt" field.
func (u *ChannelUpsert) SetPrompt(v string) *ChannelUpsert {
	u.Set(channel.FieldPrompt, v)
	return u
}

// UpdatePrompt sets the "prompt" field to the value that was provided on create.
func (u *ChannelUpsert) UpdatePrompt() *ChannelUpsert {
	u.SetExcluded(channel.FieldPrompt)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ChannelUpsert) SetUpdatedAt(v time.Time) *ChannelUpsert {
	u.Set(channel.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ChannelUpsert) UpdateUpdatedAt() *ChannelUpsert {
	u.SetExcluded(channel.FieldUpdatedAt)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Channel.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(channel.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ChannelUpsertOne) UpdateNewValues() *ChannelUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(channel.FieldID)
		}
		if _, exists := u.create.mutation.GuildID(); exists {
			s.SetIgnore(channel.FieldGuildID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(channel.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Channel.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ChannelUpsertOne) Ignore() *ChannelUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ChannelUpsertOne) DoNothing() *ChannelUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ChannelCreate.OnConflict
// documentation for more info.
func (u *ChannelUpsertOne) Update(set func(*ChannelUpsert)) *ChannelUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ChannelUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "name" field.
func (u *ChannelUpsertOne) SetName(v string) *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ChannelUpsertOne) UpdateName() *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdateName()
	})
}

// SetModel sets the "model" field.
func (u *ChannelUpsertOne) SetModel(v string) *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.SetModel(v)
	})
}

// UpdateModel sets the "model" field to the value that was provided on create.
func (u *ChannelUpsertOne) UpdateModel() *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdateModel()
	})
}

// SetPrompt sets the "prompt" field.
func (u *ChannelUpsertOne) SetPrompt(v string) *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.SetPrompt(v)
	})
}

// UpdatePrompt sets the "prompt" field to the value that was provided on create.
func (u *ChannelUpsertOne) UpdatePrompt() *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdatePrompt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ChannelUpsertOne) SetUpdatedAt(v time.Time) *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ChannelUpsertOne) UpdateUpdatedAt() *ChannelUpsertOne {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *ChannelUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ChannelCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ChannelUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ChannelUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: ChannelUpsertOne.ID is not supported by MySQL driver. Use ChannelUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ChannelUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ChannelCreateBulk is the builder for creating many Channel entities in bulk.
type ChannelCreateBulk struct {
	config
	err      error
	builders []*ChannelCreate
	conflict []sql.ConflictOption
}

// Save creates the Channel entities in the database.
func (ccb *ChannelCreateBulk) Save(ctx context.Context) ([]*Channel, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Channel, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ChannelMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *ChannelCreateBulk) SaveX(ctx context.Context) []*Channel {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *ChannelCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *ChannelCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Channel.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ChannelUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (ccb *ChannelCreateBulk) OnConflict(opts ...sql.ConflictOption) *ChannelUpsertBulk {
	ccb.conflict = opts
	return &ChannelUpsertBulk{
		create: ccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Channel.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ccb *ChannelCreateBulk) OnConflictColumns(columns ...string) *ChannelUpsertBulk {
	ccb.conflict = append(ccb.conflict, sql.ConflictColumns(columns...))
	return &ChannelUpsertBulk{
		create: ccb,
	}
}

// ChannelUpsertBulk is the builder for "upsert"-ing
// a bulk of Channel nodes.
type ChannelUpsertBulk struct {
	create *ChannelCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Channel.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(channel.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ChannelUpsertBulk) UpdateNewValues() *ChannelUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(channel.FieldID)
			}
			if _, exists := b.mutation.GuildID(); exists {
				s.SetIgnore(channel.FieldGuildID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(channel.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Channel.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ChannelUpsertBulk) Ignore() *ChannelUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ChannelUpsertBulk) DoNothing() *ChannelUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ChannelCreateBulk.OnConflict
// documentation for more info.
func (u *ChannelUpsertBulk) Update(set func(*ChannelUpsert)) *ChannelUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ChannelUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "name" field.
func (u *ChannelUpsertBulk) SetName(v string) *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ChannelUpsertBulk) UpdateName() *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdateName()
	})
}

// SetModel sets the "model" field.
func (u *ChannelUpsertBulk) SetModel(v string) *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.SetModel(v)
	})
}

// UpdateModel sets the "model" field to the value that was provided on create.
func (u *ChannelUpsertBulk) UpdateModel() *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdateModel()
	})
}

// SetPrompt sets the "prompt" field.
func (u *ChannelUpsertBulk) SetPrompt(v string) *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.SetPrompt(v)
	})
}

// UpdatePrompt sets the "prompt" field to the value that was provided on create.
func (u *ChannelUpsertBulk) UpdatePrompt() *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdatePrompt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ChannelUpsertBulk) SetUpdatedAt(v time.Time) *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ChannelUpsertBulk) UpdateUpdatedAt() *ChannelUpsertBulk {
	return u.Update(func(s *ChannelUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *ChannelUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the ChannelCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ChannelCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ChannelUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
